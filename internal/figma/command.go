package figma

import (
	"encoding/json"
	"fmt"

	"github.com/nerveband/ai-happy-design-v2/internal/ws"
)

// Commander sends commands to the Figma plugin through the WebSocket relay.
type Commander struct {
	wsServer *ws.Server
}

// NewCommander creates a Commander that sends commands via the given WebSocket server.
func NewCommander(wsServer *ws.Server) *Commander {
	return &Commander{
		wsServer: wsServer,
	}
}

// SendCommand sends a command to the Figma plugin on the currently active channel
// and returns the decoded result.
func (c *Commander) SendCommand(command string, params map[string]interface{}) (interface{}, error) {
	channel := c.wsServer.GetChannelKey()
	if channel == "" {
		return nil, fmt.Errorf("no Figma plugin connected - use 'connect' command in the Figma plugin first")
	}

	resultJSON, err := c.wsServer.SendCommand(channel, command, params)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resultJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to decode result: %w", err)
	}
	return result, nil
}

// SendCommandRaw sends a command and returns the raw JSON response.
func (c *Commander) SendCommandRaw(command string, params map[string]interface{}) (json.RawMessage, error) {
	channel := c.wsServer.GetChannelKey()
	if channel == "" {
		return nil, fmt.Errorf("no Figma plugin connected - use 'connect' command in the Figma plugin first")
	}
	return c.wsServer.SendCommand(channel, command, params)
}

// IsConnected returns true if at least one Figma plugin is connected.
func (c *Commander) IsConnected() bool {
	return c.wsServer.GetChannelKey() != ""
}
