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

	done chan struct{}
}

// NewClient creates a new WebSocket client targeting the given URL.
func NewClient(url string) *Client {
	return &Client{
		url:     url,
		pending: make(map[string]chan *Message),
		done:    make(chan struct{}),
	}
}

// JoinChannel connects to the relay and joins the specified channel.
// Blocks until the connection is closed or an error occurs.
func (c *Client) JoinChannel(channelKey string) error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url+"/ws", nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
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

	log.Printf("[client] joined channel %s", channelKey)

	// Block until done
	<-c.done
	return nil
}

// SendCommand sends a command to the plugin via the relay and waits for a response.
func (c *Client) SendCommand(command string, params map[string]interface{}) (json.RawMessage, error) {
	id := uuid.New().String()

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	msg := Message{
		ID:      id,
		Type:    "command",
		Channel: c.channel,
		Command: command,
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
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("command %q timed out after 60s", command)
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
	close(c.done)
	if c.ws != nil {
		return c.ws.Close()
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

	for {
		_, data, err := c.ws.ReadMessage()
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
		case "response", "error":
			c.pendingMu.Lock()
			ch, ok := c.pending[msg.ID]
			c.pendingMu.Unlock()
			if ok {
				ch <- &msg
			}
		case "joined":
			log.Printf("[client] confirmed join on channel %s", msg.Channel)
		case "progress_update":
			log.Printf("[client] progress: %s", string(msg.Result))
		}
	}
}
