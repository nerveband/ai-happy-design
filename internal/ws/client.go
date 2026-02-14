package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client is a WebSocket client that connects to the relay server.
type Client struct {
	url     string
	ws      *websocket.Conn
	channel string
	mu      sync.Mutex

	pending   map[string]chan *Message
	pendingMu sync.Mutex

	done            chan struct{}
	joined          chan struct{}
	progressHandler func(*Message)
}

// NewClient creates a new WebSocket client targeting the given URL.
func NewClient(url string) *Client {
	return &Client{
		url:     url,
		pending: make(map[string]chan *Message),
		done:    make(chan struct{}),
		joined:  make(chan struct{}, 1),
	}
}

// SetProgressHandler installs an optional callback for progress_update messages.
func (c *Client) SetProgressHandler(handler func(*Message)) {
	c.progressHandler = handler
}

// Connect establishes a WebSocket connection and joins the specified channel.
func (c *Client) Connect(channelKey string) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url+"/ws", nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.done = make(chan struct{})
	c.joined = make(chan struct{}, 1)
	c.ws = conn
	c.channel = channelKey

	// Start read pump
	go c.readPump()

	// Send join message
	joinMsg := Message{
		Type:    "join",
		Channel: channelKey,
	}
	data, _ := json.Marshal(joinMsg)
	if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send join: %w", err)
	}

	select {
	case <-c.joined:
		log.Printf("[client] joined channel %s", channelKey)
		return nil
	case <-time.After(10 * time.Second):
		_ = c.Close()
		return fmt.Errorf("timed out waiting to join channel %s", channelKey)
	case <-c.done:
		return fmt.Errorf("connection closed before join completed")
	}
}

// JoinChannel connects to the relay and joins the specified channel.
// Blocks until the connection is closed or an error occurs.
func (c *Client) JoinChannel(channelKey string) error {
	if err := c.Connect(channelKey); err != nil {
		return err
	}
	// Block until done
	<-c.done
	return nil
}

// SendCommand sends a command to the plugin via the relay and waits for a response.
func (c *Client) SendCommand(command string, params map[string]interface{}) (json.RawMessage, error) {
	if c.ws == nil {
		return nil, fmt.Errorf("not connected")
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

	msg := Message{
		ID:      id,
		Type:    "command",
		Channel: c.channel,
		Command: command,
		Domain:  domain,
		Action:  action,
		Params:  paramsJSON,
	}

	respCh := make(chan *Message, 1)
	c.pendingMu.Lock()
	c.pending[id] = respCh
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	c.mu.Lock()
	err = c.ws.WriteMessage(websocket.TextMessage, data)
	c.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed to send: %w", err)
	}

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

// Send writes an arbitrary JSON message to the relay.
func (c *Client) Send(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ws.WriteMessage(websocket.TextMessage, data)
}

// Close shuts down the client connection.
func (c *Client) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}

	c.mu.Lock()
	ws := c.ws
	c.ws = nil
	c.mu.Unlock()

	if ws != nil {
		return ws.Close()
	}
	return nil
}

// readPump reads messages from the WebSocket and routes responses.
func (c *Client) readPump() {
	defer func() {
		select {
		case <-c.done:
		default:
			close(c.done)
		}
	}()

	// Capture conn reference so Close() nil-ing c.ws doesn't crash us.
	c.mu.Lock()
	ws := c.ws
	c.mu.Unlock()
	if ws == nil {
		return
	}

	ws.SetReadLimit(64 * 1024 * 1024) // 64 MB – export payloads can be large

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("[client] read error: %v", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("[client] invalid message: %v", err)
			continue
		}

		switch msg.Type {
		case "response", "result", "error":
			c.pendingMu.Lock()
			ch, ok := c.pending[msg.ID]
			c.pendingMu.Unlock()
			if ok {
				ch <- &msg
			}
		case "message":
			if inner, ok := extractWrappedResponse(&msg); ok {
				c.pendingMu.Lock()
				ch, exists := c.pending[inner.ID]
				c.pendingMu.Unlock()
				if exists {
					ch <- inner
				}
			}
		case "joined":
			log.Printf("[client] confirmed join on channel %s", msg.Channel)
			select {
			case c.joined <- struct{}{}:
			default:
			}
		case "progress_update":
			if c.progressHandler != nil {
				c.progressHandler(&msg)
			} else {
				log.Printf("[client] progress: %s", string(msg.Result))
			}
		}
	}
}
