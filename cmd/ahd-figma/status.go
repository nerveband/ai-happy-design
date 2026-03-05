package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/ws"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show relay, plugin, and current page status (JSON)",
	Long: `Returns a single JSON object with relay health, plugin connectivity,
current page, and channel info. Designed for LLM agents to orient themselves
before issuing commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)

		cfg := loadConfig()
		out := map[string]interface{}{
			"relay": map[string]interface{}{
				"host": cfg.ServerHost,
				"port": cfg.Port,
			},
		}

		// Check relay health
		statusURL := fmt.Sprintf("http://%s:%d/status", cfg.ServerHost, cfg.Port)
		relayHealthy := false
		var relayInfo map[string]interface{}

		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(statusURL)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				relayHealthy = true
				if err := json.NewDecoder(resp.Body).Decode(&relayInfo); err == nil {
					out["relay"] = relayInfo
				}
			}
		}

		if !relayHealthy {
			outRelay := out["relay"].(map[string]interface{})
			outRelay["healthy"] = false
			outRelay["error"] = "Relay not running or unreachable"
			return printJSON(out)
		}

		// Try to get current page info from plugin
		channelKey, chanErr := resolveChannel("")
		if chanErr != nil {
			out["plugin"] = map[string]interface{}{
				"connected": false,
				"error":     chanErr.Error(),
			}
			return printJSON(out)
		}

		wsURL := fmt.Sprintf("ws://%s:%d", cfg.ServerHost, cfg.Port)
		wsClient := ws.NewClient(wsURL)
		if err := wsClient.Connect(channelKey); err != nil {
			out["plugin"] = map[string]interface{}{
				"connected": false,
				"error":     err.Error(),
			}
			return printJSON(out)
		}
		defer wsClient.Close()

		// Get current page info
		pageResult, pageErr := wsClient.SendCommand("page.get_all", map[string]interface{}{})
		if pageErr != nil {
			out["plugin"] = map[string]interface{}{
				"connected": true,
				"channel":   channelKey,
			}
			out["page"] = map[string]interface{}{
				"error": pageErr.Error(),
			}
		} else {
			out["plugin"] = map[string]interface{}{
				"connected": true,
				"channel":   channelKey,
			}
			var pageInfo interface{}
			if err := json.Unmarshal(pageResult, &pageInfo); err == nil {
				out["pages"] = pageInfo
			}
		}

		return printJSON(out)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
