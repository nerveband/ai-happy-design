package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message is the envelope for all WebSocket messages.
type Message struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Command string          `json:"command,omitempty"`
	Domain  string          `json:"domain,omitempty"`
	Action  string          `json:"action,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   string          `json:"error,omitempty"`
	Message json.RawMessage `json:"message,omitempty"`
	Role    string          `json:"role,omitempty"`
}

// Conn wraps a WebSocket connection with metadata.
type Conn struct {
	ws      *websocket.Conn
	id      string
	channel string
	role    string
	sendCh  chan []byte
	closed  bool // true after eviction; prevents send-on-closed-channel panics
	mu      sync.Mutex
}

// safeSend sends data to the connection's send channel, returning false if
// the connection has been evicted or the buffer is full.
func (c *Conn) safeSend(data []byte) bool {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return false
	}
	c.mu.Unlock()
	select {
	case c.sendCh <- data:
		return true
	default:
		return false
	}
}

// Server is the WebSocket relay server.
type Server struct {
	port             int
	channels         map[string]map[string]*Conn // channel -> connID -> Conn
	preferredChannel string
	mu               sync.RWMutex
	upgrader         websocket.Upgrader

	// Pending commands: requestID -> response channel
	pending   map[string]chan *Message
	pendingMu sync.Mutex

	// Client mode: when port is already in use, delegate to an embedded client
	client *Client

	// Idle auto-shutdown
	idleTimeout time.Duration
	lastActivity time.Time
	idleTimer   *time.Timer
	startedAt   time.Time
	httpServer  *http.Server
}

// NewServer creates a new WebSocket relay server on the given port.
func NewServer(port int) *Server {
	return &Server{
		port:     port,
		channels: make(map[string]map[string]*Conn),
		pending:  make(map[string]chan *Message),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		startedAt:    time.Now(),
		lastActivity: time.Now(),
	}
}

// SetIdleTimeout configures auto-shutdown after a period of no activity.
// Zero disables idle shutdown.
func (s *Server) SetIdleTimeout(d time.Duration) {
	s.idleTimeout = d
}

// touchActivity resets the idle timer on meaningful activity.
func (s *Server) touchActivity() {
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()
	if s.idleTimeout > 0 {
		s.resetIdleTimer()
	}
}

// resetIdleTimer resets or starts the idle shutdown timer.
func (s *Server) resetIdleTimer() {
	if s.idleTimer != nil {
		s.idleTimer.Stop()
	}
	s.idleTimer = time.AfterFunc(s.idleTimeout, func() {
		log.Printf("[ws] Relay shutting down after %s of inactivity. CLI will auto-restart on next command.", s.idleTimeout)
		if s.httpServer != nil {
			s.httpServer.Close()
		}
	})
}

// Start launches the HTTP server with WebSocket and status endpoints.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	mux.HandleFunc("/status", s.handleStatus)

	addr := fmt.Sprintf(":%d", s.port)
	s.httpServer = &http.Server{Addr: addr, Handler: mux}

	// Start idle timer if configured
	if s.idleTimeout > 0 {
		s.resetIdleTimer()
		log.Printf("[ws] idle auto-shutdown enabled: %s", s.idleTimeout)
	}

	log.Printf("[ws] relay server listening on %s", addr)
	return s.httpServer.ListenAndServe()
}

// ConnectAsClient connects to an existing relay as a WebSocket client.
// This is used when the MCP server can't start its own relay (port in use).
func (s *Server) ConnectAsClient(wsURL string) error {
	// NewClient expects a base URL without /ws — Client.Connect appends "/ws" itself.
	baseURL := fmt.Sprintf("ws://localhost:%d", s.port)
	c := NewClient(baseURL)

	// Probe the existing relay for its preferred channel
	probeConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to relay: %w", err)
	}
	probeConn.Close()

	// Get channel from relay status endpoint
	statusURL := fmt.Sprintf("http://localhost:%d/status", s.port)
	resp, err := http.Get(statusURL)
	if err != nil {
		return fmt.Errorf("failed to get relay status: %w", err)
	}
	defer resp.Body.Close()
	var status struct {
		PreferredChannel string         `json:"preferredChannel"`
		Channels         map[string]int `json:"channels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return fmt.Errorf("failed to decode relay status: %w", err)
	}

	channel := status.PreferredChannel
	if channel == "" {
		for ch := range status.Channels {
			channel = ch
			break
		}
	}
	if channel == "" {
		return fmt.Errorf("no active channel on relay")
	}

	if err := c.Connect(channel); err != nil {
		return fmt.Errorf("failed to join channel %s: %w", channel, err)
	}

	s.client = c
	s.preferredChannel = channel
	log.Printf("[mcp] connected as client to relay channel %s", channel)
	return nil
}

// SendCommand sends a command to all connections in the given channel and waits
// for a response. Returns the result or an error.
func (s *Server) SendCommand(channel, command string, params map[string]interface{}) (json.RawMessage, error) {
	// If running in client mode, delegate to the embedded client
	if s.client != nil {
		return s.client.SendCommand(command, params)
	}

	id := uuid.New().String()

	domain, action, err := resolveCommandRoute(command, params)
	if err != nil {
		return nil, err
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	msg := &Message{
		ID:      id,
		Type:    "command",
		Channel: channel,
		Command: command,
		Domain:  domain,
		Action:  action,
		Params:  paramsJSON,
	}

	respCh := make(chan *Message, 1)
	s.pendingMu.Lock()
	s.pending[id] = respCh
	s.pendingMu.Unlock()

	defer func() {
		s.pendingMu.Lock()
		delete(s.pending, id)
		s.pendingMu.Unlock()
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	s.broadcast(channel, data, "")

	select {
	case resp := <-respCh:
		if resp.Error != "" {
			return nil, fmt.Errorf("figma error: %s", resp.Error)
		}
		return resp.Result, nil
	case <-time.After(300 * time.Second):
		return nil, fmt.Errorf("command %q timed out after 300s", command)
	}
}

// GetChannelKey returns the preferred active channel key, or a deterministic
// active fallback if no preferred channel is currently connected.
func (s *Server) GetChannelKey() string {
	// In client mode, return the channel we joined
	if s.client != nil {
		return s.preferredChannel
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.preferredChannel != "" {
		if conns, ok := s.channels[s.preferredChannel]; ok && len(conns) > 0 {
			return s.preferredChannel
		}
	}

	next := firstActiveChannelLocked(s.channels)
	s.preferredChannel = next
	return next
}

// HasChannel returns true if the named channel has at least one connection.
func (s *Server) HasChannel(channel string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conns, ok := s.channels[channel]
	return ok && len(conns) > 0
}

// ChannelCount returns the number of active channels.
func (s *Server) ChannelCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, conns := range s.channels {
		if len(conns) > 0 {
			count++
		}
	}
	return count
}

// handleWS upgrades HTTP to WebSocket and manages the connection lifecycle.
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	conn := &Conn{
		ws:     wsConn,
		id:     uuid.New().String(),
		sendCh: make(chan []byte, 256),
	}

	// Read pump
	go s.readPump(conn)
	// Write pump
	go s.writePump(conn)
}

// readPump reads messages from the WebSocket and routes them.
func (s *Server) readPump(conn *Conn) {
	defer func() {
		s.removeConn(conn)
		conn.ws.Close()
	}()

	conn.ws.SetReadLimit(64 * 1024 * 1024) // 64 MB – export payloads can be large
	conn.ws.SetReadDeadline(time.Now().Add(120 * time.Second))
	conn.ws.SetPongHandler(func(string) error {
		conn.ws.SetReadDeadline(time.Now().Add(120 * time.Second))
		return nil
	})

	for {
		_, data, err := conn.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[ws] read error: %v", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("[ws] invalid message: %v", err)
			continue
		}

		s.routeMessage(conn, &msg, data)
	}
}

// writePump writes queued messages to the WebSocket.
func (s *Server) writePump(conn *Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		conn.ws.Close()
	}()

	for {
		select {
		case msg, ok := <-conn.sendCh:
			if !ok {
				conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.ws.SetWriteDeadline(time.Now().Add(120 * time.Second))
			if err := conn.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// routeMessage handles incoming messages based on their type.
func (s *Server) routeMessage(conn *Conn, msg *Message, raw []byte) {
	if msg.Type != "ping" {
		s.touchActivity()
	}

	switch msg.Type {
	case "ping":
		pong, _ := json.Marshal(Message{Type: "pong"})
		conn.safeSend(pong)

	case "join":
		s.handleJoin(conn, msg)

	case "response", "result", "error":
		if !s.routePendingResponse(msg) && conn.channel != "" {
			// Keep CLI clients working by forwarding unmatched responses.
			s.broadcast(conn.channel, raw, conn.id)
		}

	case "message":
		inner, ok := extractWrappedResponse(msg)
		if ok {
			if !s.routePendingResponse(inner) && conn.channel != "" {
				// Broadcast wrapped envelope for WebSocket clients waiting on this ID.
				s.broadcast(conn.channel, raw, conn.id)
			}
			return
		}
		if conn.channel != "" {
			s.broadcast(conn.channel, raw, conn.id)
		}

	case "progress_update":
		// Broadcast to other channel members (for UI updates)
		s.broadcast(conn.channel, raw, conn.id)

	default:
		// Broadcast to channel
		if conn.channel != "" {
			s.broadcast(conn.channel, raw, conn.id)
		}
	}
}

func (s *Server) routePendingResponse(msg *Message) bool {
	if msg.ID == "" {
		return false
	}

	s.pendingMu.Lock()
	ch, ok := s.pending[msg.ID]
	s.pendingMu.Unlock()
	if ok {
		ch <- msg
		return true
	}
	return false
}

func extractWrappedResponse(msg *Message) (*Message, bool) {
	if len(msg.Message) == 0 {
		return nil, false
	}

	// Presence checks (instead of truthiness) prevent falsey result values from timing out.
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(msg.Message, &payload); err != nil {
		return nil, false
	}

	idRaw, hasID := payload["id"]
	_, hasResult := payload["result"]
	_, hasError := payload["error"]
	if !hasID || (!hasResult && !hasError) {
		return nil, false
	}

	var id string
	if err := json.Unmarshal(idRaw, &id); err != nil || id == "" {
		return nil, false
	}

	out := &Message{ID: id}
	if hasResult {
		out.Type = "response"
		out.Result = payload["result"]
	}
	if hasError {
		out.Type = "error"
		var errText string
		_ = json.Unmarshal(payload["error"], &errText)
		out.Error = errText
	}
	return out, true
}

// handleJoin adds a connection to a channel.
func (s *Server) handleJoin(conn *Conn, msg *Message) {
	channel := msg.Channel
	if channel == "" {
		channel = GenerateChannelKey()
	}

	conn.channel = channel
	conn.role = msg.Role

	s.mu.Lock()
	if s.channels[channel] == nil {
		s.channels[channel] = make(map[string]*Conn)
	}

	// Evict stale connections with the same role (e.g. old plugin instances).
	// This prevents duplicate command execution when a plugin reconnects.
	if conn.role != "" {
		var evict []*Conn
		for id, existing := range s.channels[channel] {
			if existing.role == conn.role && id != conn.id {
				evict = append(evict, existing)
				delete(s.channels[channel], id)
			}
		}
		for _, old := range evict {
			log.Printf("[ws] evicting stale %s connection %s from channel %s", old.role, old.id, channel)
			old.mu.Lock()
			old.closed = true
			old.mu.Unlock()
			close(old.sendCh)
		}
	}

	s.channels[channel][conn.id] = conn
	if s.preferredChannel == "" {
		s.preferredChannel = channel
	}
	s.mu.Unlock()

	log.Printf("[ws] %s joined channel %s (role=%s)", conn.id, channel, conn.role)

	// Send confirmation
	resp := Message{
		Type:    "joined",
		Channel: channel,
	}
	data, _ := json.Marshal(resp)
	conn.safeSend(data)
}

// broadcast sends a message to all connections in a channel, optionally
// excluding a specific connection.
func (s *Server) broadcast(channel string, data []byte, excludeID string) {
	s.mu.RLock()
	conns := s.channels[channel]
	s.mu.RUnlock()

	for id, conn := range conns {
		if id == excludeID {
			continue
		}
		if !conn.safeSend(data) {
			log.Printf("[ws] send buffer full or evicted for %s, dropping message", id)
		}
	}
}

// removeConn removes a connection from its channel.
func (s *Server) removeConn(conn *Conn) {
	if conn.channel == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if conns, ok := s.channels[conn.channel]; ok {
		delete(conns, conn.id)
		if len(conns) == 0 {
			delete(s.channels, conn.channel)
		}
	}
	if conn.channel == s.preferredChannel {
		if conns, ok := s.channels[conn.channel]; !ok || len(conns) == 0 {
			s.preferredChannel = firstActiveChannelLocked(s.channels)
		}
	}
	log.Printf("[ws] %s left channel %s", conn.id, conn.channel)
}

// handleStatus returns basic server status as JSON.
func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	preferred := s.GetChannelKey()

	s.mu.RLock()
	channels := make(map[string]int)
	for ch, conns := range s.channels {
		channels[ch] = len(conns)
	}
	s.mu.RUnlock()

	status := map[string]interface{}{
		"status":           "ok",
		"channels":         channels,
		"preferredChannel": preferred,
		"uptime":           time.Since(s.startedAt).Truncate(time.Second).String(),
	}

	if s.idleTimeout > 0 {
		s.mu.RLock()
		lastAct := s.lastActivity
		s.mu.RUnlock()
		remaining := s.idleTimeout - time.Since(lastAct)
		if remaining < 0 {
			remaining = 0
		}
		status["idleTimeout"] = s.idleTimeout.String()
		status["idleRemaining"] = remaining.Truncate(time.Second).String()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func firstActiveChannelLocked(channels map[string]map[string]*Conn) string {
	names := make([]string, 0, len(channels))
	for ch, conns := range channels {
		if len(conns) > 0 {
			names = append(names, ch)
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return names[0]
}
