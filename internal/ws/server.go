package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Message is the envelope for all WebSocket messages.
type Message struct {
	ID      string          `json:"id,omitempty"`
	Type    string          `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Command string          `json:"command,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// Conn wraps a WebSocket connection with metadata.
type Conn struct {
	ws       *websocket.Conn
	id       string
	channel  string
	sendCh   chan []byte
	mu       sync.Mutex
}

// Server is the WebSocket relay server.
type Server struct {
	port     int
	channels map[string]map[string]*Conn // channel -> connID -> Conn
	mu       sync.RWMutex
	upgrader websocket.Upgrader

	// Pending commands: requestID -> response channel
	pending   map[string]chan *Message
	pendingMu sync.Mutex
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
	}
}

// Start launches the HTTP server with WebSocket and status endpoints.
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	mux.HandleFunc("/status", s.handleStatus)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("[ws] relay server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

// SendCommand sends a command to all connections in the given channel and waits
// for a response. Returns the result or an error.
func (s *Server) SendCommand(channel, command string, params map[string]interface{}) (json.RawMessage, error) {
	id := uuid.New().String()

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	msg := &Message{
		ID:      id,
		Type:    "command",
		Channel: channel,
		Command: command,
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
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("command %q timed out after 60s", command)
	}
}

// GetChannelKey returns the first active channel key, or generates a new one
// if no channels exist. This is used by the MCP server to find the current
// plugin connection.
func (s *Server) GetChannelKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.channels {
		if len(s.channels[ch]) > 0 {
			return ch
		}
	}
	return ""
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

	conn.ws.SetReadLimit(1024 * 1024) // 1 MB
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
			conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
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
	switch msg.Type {
	case "join":
		s.handleJoin(conn, msg)

	case "response", "error":
		// Route to pending command
		s.pendingMu.Lock()
		ch, ok := s.pending[msg.ID]
		s.pendingMu.Unlock()
		if ok {
			ch <- msg
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

// handleJoin adds a connection to a channel.
func (s *Server) handleJoin(conn *Conn, msg *Message) {
	channel := msg.Channel
	if channel == "" {
		channel = GenerateChannelKey()
	}

	conn.channel = channel

	s.mu.Lock()
	if s.channels[channel] == nil {
		s.channels[channel] = make(map[string]*Conn)
	}
	s.channels[channel][conn.id] = conn
	s.mu.Unlock()

	log.Printf("[ws] %s joined channel %s", conn.id, channel)

	// Send confirmation
	resp := Message{
		Type:    "joined",
		Channel: channel,
	}
	data, _ := json.Marshal(resp)
	conn.sendCh <- data
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
		select {
		case conn.sendCh <- data:
		default:
			log.Printf("[ws] send buffer full for %s, dropping message", id)
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
	log.Printf("[ws] %s left channel %s", conn.id, conn.channel)
}

// handleStatus returns basic server status as JSON.
func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	channels := make(map[string]int)
	for ch, conns := range s.channels {
		channels[ch] = len(conns)
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"channels": channels,
	})
}
