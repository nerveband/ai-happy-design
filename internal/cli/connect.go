package cli

import (
	"fmt"

	"github.com/nerveband/ai-happy-design/internal/ws"
)

// Connect joins a WebSocket channel and blocks until the connection is closed.
func Connect(serverURL, channelKey string) error {
	client := ws.NewClient(serverURL)
	fmt.Printf("Connecting to %s on channel %s...\n", serverURL, channelKey)
	return client.JoinChannel(channelKey)
}
