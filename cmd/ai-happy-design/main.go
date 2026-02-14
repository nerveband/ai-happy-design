package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/nerveband/ai-happy-design/internal/mcp"
	pluginpkg "github.com/nerveband/ai-happy-design/internal/plugin"
	relaymgr "github.com/nerveband/ai-happy-design/internal/relay"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

var version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "ai-happy-design",
	Short: "AI Happy Design - Figma MCP server and CLI",
	Long: `A single binary that works as both an MCP server for AI editors
and a CLI for direct Figma manipulation via a WebSocket relay.`,
	Version: version,
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	Long: `Starts the MCP server on stdio for use with AI editors (Claude Code,
Cursor, Windsurf, etc). Also starts a WebSocket relay server in the
background for communicating with the Figma plugin.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		wsServer := ws.NewServer(cfg.Port)
		go wsServer.Start()
		return mcp.StartServer(wsServer)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [channel-key]",
	Short: "Connect to a Figma plugin channel",
	Long: `Connect to the WebSocket relay and join a specific channel.
The channel key should match what the Figma plugin is using
(e.g. happy-unicorn-42). If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}
		cfg := config.Load()
		url := fmt.Sprintf("ws://%s:%d", cfg.ServerHost, cfg.Port)
		client := ws.NewClient(url)
		return client.JoinChannel(args[0])
	},
}

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "Start WebSocket relay server only",
	Long:  `Starts just the WebSocket relay server without the MCP server. Useful for debugging or running the relay separately.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		wsServer := ws.NewServer(cfg.Port)
		fmt.Printf("Starting WebSocket relay on port %d...\n", cfg.Port)
		return wsServer.Start()
	},
}

var setupPath string
var setupForce bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Extract and install the Figma plugin",
	Long:  `Extracts the embedded Figma plugin files to a local directory so you can import it into Figma.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := setupPath
		if dir == "" {
			var err error
			dir, err = pluginpkg.DefaultPluginDir()
			if err != nil {
				return fmt.Errorf("failed to determine plugin directory: %w", err)
			}
		}

		if !setupForce && !pluginpkg.NeedsUpdate(dir, version) {
			fmt.Printf("Plugin is already up to date (version %s) at:\n  %s\n", version, dir)
			fmt.Println("Use --force to re-extract.")
			return nil
		}

		fmt.Printf("Extracting plugin to %s...\n", dir)
		if err := pluginpkg.ExtractTo(dir, version); err != nil {
			return fmt.Errorf("failed to extract plugin: %w", err)
		}

		manifestPath := filepath.Join(dir, "manifest.json")
		fmt.Printf("Plugin extracted successfully!\n\n")
		fmt.Printf("Files:\n  %s\n\n", dir)

		if runtime.GOOS == "darwin" {
			exec.Command("open", "-R", manifestPath).Run()
			fmt.Println("Finder opened to the plugin directory.")
		}

		fmt.Println("To load in Figma:")
		fmt.Println("  1. Open Figma Desktop")
		fmt.Println("  2. Plugins > Development > Import plugin from manifest")
		fmt.Printf("  3. Select: %s\n", manifestPath)
		fmt.Println("  4. Run: Plugins > AI Happy Design")

		return nil
	},
}

var commandParams string
var commandLive bool
var commandChannel string

var commandCmd = &cobra.Command{
	Use:   "command [channel-key] <command>",
	Short: "Execute a command against a connected Figma channel",
	Long: `Executes one command through the WebSocket relay for scripting and LLM-driven CLI flows.
Command supports both legacy command names (e.g. set_fill_color) and domain actions (e.g. paint.set_solid).
Channel resolution order: positional arg, --channel, AHD_CHANNEL env, relay preferred/active channel.
If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		var explicitChannel string
		var command string
		switch len(args) {
		case 1:
			explicitChannel = commandChannel
			command = args[0]
		case 2:
			if commandChannel != "" {
				return fmt.Errorf("provide channel either as positional argument or --channel, not both")
			}
			explicitChannel = args[0]
			command = args[1]
		default:
			return fmt.Errorf("invalid command arguments")
		}

		channelKey, err := resolveChannel(explicitChannel)
		if err != nil {
			return err
		}

		params := map[string]interface{}{}
		if commandParams != "" {
			if err := json.Unmarshal([]byte(commandParams), &params); err != nil {
				return fmt.Errorf("invalid --params JSON: %w", err)
			}
		}

		client, err := newConnectedClient(channelKey, commandLive)
		if err != nil {
			return err
		}
		defer client.Close()

		result, err := client.SendCommand(command, params)
		if err != nil {
			return err
		}

		var pretty interface{}
		if err := json.Unmarshal(result, &pretty); err == nil {
			return printJSON(pretty)
		}
		fmt.Println(string(result))
		return nil
	},
}

var batchOperations string
var batchOperationsFile string
var batchLive bool
var batchChannel string
var batchFailFast bool
var batchRetries int
var batchRetryDelayMs int
var batchInterpolation bool

type batchOperation struct {
	Name    string                 `json:"name,omitempty"`
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

var batchCmd = &cobra.Command{
	Use:   "batch [channel-key]",
	Short: "Execute a sequence of commands in order",
	Long: `Runs multiple commands over one connection to support tool chaining workflows.
Provide operations as JSON array: [{"name":"createCard","command":"shape.create_rectangle","params":{"x":40,"y":40,"width":220,"height":120}}].
Supports interpolation placeholders like ${{steps.0.result.id}} or ${{steps.createCard.result.id}}.
Channel resolution order: positional arg, --channel, AHD_CHANNEL env, relay preferred/active channel.
If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		explicitChannel := batchChannel
		if len(args) == 1 {
			if batchChannel != "" {
				return fmt.Errorf("provide channel either as positional argument or --channel, not both")
			}
			explicitChannel = args[0]
		}
		channelKey, err := resolveChannel(explicitChannel)
		if err != nil {
			return err
		}

		ops, err := loadBatchOperations(batchOperations, batchOperationsFile)
		if err != nil {
			return err
		}
		if len(ops) == 0 {
			return fmt.Errorf("operations array is empty")
		}
		if batchRetries < 0 {
			return fmt.Errorf("--retries must be >= 0")
		}
		if batchRetryDelayMs < 0 {
			return fmt.Errorf("--retry-delay-ms must be >= 0")
		}

		client, err := newConnectedClient(channelKey, batchLive)
		if err != nil {
			return err
		}
		defer client.Close()

		results := make([]map[string]interface{}, 0, len(ops))
		stepStates := make([]batchutil.StepState, 0, len(ops))
		retryDelay := time.Duration(batchRetryDelayMs) * time.Millisecond
		succeeded := 0
		failed := 0
		retriesUsed := 0
		stoppedEarly := false

		for i, op := range ops {
			params := op.Params
			if params == nil {
				params = map[string]interface{}{}
			}

			if batchInterpolation {
				interpolatedParams, err := batchutil.InterpolateParams(params, stepStates)
				if err != nil {
					failed++
					entry := map[string]interface{}{
						"index":    i,
						"name":     op.Name,
						"command":  op.Command,
						"ok":       false,
						"error":    fmt.Sprintf("interpolation error: %v", err),
						"attempts": 0,
					}
					results = append(results, entry)
					stepStates = append(stepStates, batchutil.StepState{
						Index:   i,
						Name:    op.Name,
						Command: op.Command,
						OK:      false,
						Error:   entry["error"].(string),
					})
					if batchFailFast {
						stoppedEarly = true
						break
					}
					continue
				}
				params = interpolatedParams
			}

			maxAttempts := batchRetries + 1
			attempts := 0
			var result json.RawMessage
			var sendErr error
			for attempts = 1; attempts <= maxAttempts; attempts++ {
				result, sendErr = client.SendCommand(op.Command, params)
				if sendErr == nil {
					break
				}
				if attempts < maxAttempts && retryDelay > 0 {
					time.Sleep(retryDelay)
				}
			}
			if attempts > 1 {
				retriesUsed += attempts - 1
			}

			if sendErr != nil {
				failed++
				entry := map[string]interface{}{
					"index":    i,
					"name":     op.Name,
					"command":  op.Command,
					"ok":       false,
					"attempts": attempts,
					"error":    sendErr.Error(),
				}
				results = append(results, entry)
				stepStates = append(stepStates, batchutil.StepState{
					Index:   i,
					Name:    op.Name,
					Command: op.Command,
					OK:      false,
					Error:   sendErr.Error(),
				})
				if batchFailFast {
					stoppedEarly = true
					break
				}
				continue
			}

			var parsed interface{}
			if err := json.Unmarshal(result, &parsed); err != nil {
				parsed = string(result)
			}
			succeeded++
			entry := map[string]interface{}{
				"index":    i,
				"name":     op.Name,
				"command":  op.Command,
				"ok":       true,
				"attempts": attempts,
				"result":   parsed,
			}
			results = append(results, entry)
			stepStates = append(stepStates, batchutil.StepState{
				Index:   i,
				Name:    op.Name,
				Command: op.Command,
				OK:      true,
				Result:  parsed,
			})
		}

		processed := len(results)
		pending := len(ops) - processed
		out := map[string]interface{}{
			"ok": failed == 0 && pending == 0,
			"summary": map[string]interface{}{
				"total":         len(ops),
				"processed":     processed,
				"succeeded":     succeeded,
				"failed":        failed,
				"pending":       pending,
				"retriesUsed":   retriesUsed,
				"failFast":      batchFailFast,
				"interpolation": batchInterpolation,
			},
			"stoppedEarly": stoppedEarly,
			"results":      results,
		}
		return printJSON(out)
	},
}

var catalogJSON bool
var catalogLLM bool
var noAutoRelay bool
var relayLogsLines int
var actionsJSON bool

var actionCatalog = map[string][]string{
	"document":  {"get_info", "get_selection", "set_selection", "scan_text", "scan_by_type", "get_styles", "find_by_name", "find_by_type", "focus", "zoom_to"},
	"node":      {"get_info", "create_frame", "move", "resize", "rotate", "set_opacity", "set_blend_mode", "set_visibility", "set_locked", "rename", "delete", "clone", "set_corner_radius", "get_tree"},
	"layer":     {"set_order", "bring_forward", "send_backward", "bring_to_front", "send_to_back", "group", "ungroup", "move_to_parent", "insert_child"},
	"layout":    {"set_auto_layout", "set_padding", "set_spacing", "set_alignment", "set_sizing", "set_constraints", "set_layout_wrap", "set_wrap", "remove_auto_layout"},
	"export":    {"image", "svg", "pdf", "json"},
	"variable":  {"create", "get_all", "set_value", "bind", "unbind", "create_collection", "get_collections", "delete"},
	"style":     {"create_paint", "create_text", "create_effect", "apply", "get_all", "remove"},
	"component": {"create", "create_instance", "create_set", "get_local", "get_remote", "get_overrides", "set_overrides", "detach_instance", "reset_instance", "swap_instance"},
	"boolean":   {"union", "subtract", "intersect", "exclude", "flatten"},
	"paint":     {"set_solid", "set_gradient", "set_image_fill", "set_image", "set_image_fill_from_url", "set_image_url", "add_fill", "remove_fill", "get_fills", "set_stroke"},
	"effect":    {"set_effects", "add_shadow", "add_blur", "apply_style", "remove", "remove_effect", "get_effects"},
	"page":      {"create", "delete", "rename", "set_current", "get_all", "get_current", "duplicate"},
	"shape":     {"create_rectangle", "create_ellipse", "create_polygon", "create_star", "create_line", "create_from_svg"},
	"text":      {"create", "set_content", "set_font", "set_size", "set_weight", "set_color", "set_align", "set_spacing", "set_line_height", "set_letter_spacing", "set_decoration", "set_case", "set_paragraph_spacing", "get_content", "get_segments", "load_font", "set_style_id"},
}

var actionsCmd = &cobra.Command{
	Use:   "actions [domain]",
	Short: "List all available domain.action pairs",
	Long: `Lists all available domain.action pairs that can be used with the command and batch subcommands.
With no arguments, lists all domains and their actions.
With a domain argument, lists only that domain's actions.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			domain := args[0]
			actions, ok := actionCatalog[domain]
			if !ok {
				available := make([]string, 0, len(actionCatalog))
				for d := range actionCatalog {
					available = append(available, d)
				}
				sort.Strings(available)
				return fmt.Errorf("unknown domain %q. Available domains: %s", domain, strings.Join(available, ", "))
			}

			if actionsJSON {
				out := map[string]interface{}{
					"domains": map[string][]string{
						domain: actions,
					},
				}
				return printJSON(out)
			}

			fmt.Printf("%s:\n", domain)
			printActionsWrapped(actions)
			fmt.Println()
			return nil
		}

		// All domains
		domains := make([]string, 0, len(actionCatalog))
		for d := range actionCatalog {
			domains = append(domains, d)
		}
		sort.Strings(domains)

		if actionsJSON {
			out := map[string]interface{}{
				"domains": actionCatalog,
			}
			return printJSON(out)
		}

		for _, domain := range domains {
			fmt.Printf("%s:\n", domain)
			printActionsWrapped(actionCatalog[domain])
			fmt.Println()
		}
		return nil
	},
}

func printActionsWrapped(actions []string) {
	const maxWidth = 72
	line := " "
	for _, a := range actions {
		entry := " " + a
		if len(line)+len(entry) > maxWidth && line != " " {
			fmt.Println(line)
			line = "  " + a
		} else {
			line += entry
		}
	}
	if line != " " {
		fmt.Println(line)
	}
}

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Print discoverable tool and action catalog",
	Long: `Outputs the tool/action catalog so LLMs and scripts can discover capabilities before planning command chains.
Use --llm --json for enriched examples and recommended execution playbook.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if catalogLLM {
			return printJSON(tools.LLMCatalog())
		}

		catalog := tools.ToolCatalog()
		if catalogJSON {
			return printJSON(catalog)
		}

		for toolName, actions := range catalog {
			fmt.Printf("%s\n", toolName)
			for action, desc := range actions {
				fmt.Printf("  - %s: %s\n", action, desc)
			}
		}
		return nil
	},
}

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Manage local relay lifecycle",
	Long:  "Controls local WebSocket relay process used by CLI/plugin communication.",
}

var relayStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start local relay (idempotent)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		result, err := relaymgr.Ensure(relaymgr.EnsureOptions{
			Host:      cfg.ServerHost,
			Port:      cfg.Port,
			AutoStart: true,
		})
		if err != nil {
			return err
		}
		out := map[string]interface{}{
			"ok":             true,
			"started":        result.Started,
			"alreadyHealthy": result.AlreadyHealthy,
			"host":           cfg.ServerHost,
			"port":           cfg.Port,
			"status":         result.Status,
			"state":          result.State,
		}
		return printJSON(out)
	},
}

var relayStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop managed local relay",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := relaymgr.Stop()
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var relayStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show relay health and process metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		out := map[string]interface{}{
			"host":          cfg.ServerHost,
			"port":          cfg.Port,
			"statusURL":     fmt.Sprintf("http://%s:%d/status", cfg.ServerHost, cfg.Port),
			"autoStartHost": relaymgr.IsLocalHost(cfg.ServerHost),
		}

		status, healthy, statusErr := relaymgr.ProbeStatus(cfg.ServerHost, cfg.Port, 1200*time.Millisecond)
		out["healthy"] = healthy
		if statusErr != nil {
			out["statusError"] = statusErr.Error()
		} else {
			out["status"] = status
		}

		state, stateErr := relaymgr.LoadState()
		if stateErr != nil {
			out["stateError"] = stateErr.Error()
		}
		out["state"] = state
		if statePath, err := relaymgr.StateFilePath(); err == nil {
			out["stateFile"] = statePath
		}
		if logPath, err := relaymgr.LogFilePath(); err == nil {
			out["logFile"] = logPath
		}

		inUse, inUseErr := relaymgr.IsPortInUse(cfg.Port)
		if inUseErr != nil {
			out["portInUseError"] = inUseErr.Error()
		} else {
			out["portInUse"] = inUse
		}
		if owner, _ := relaymgr.PortOwner(cfg.Port); owner != "" {
			out["portOwner"] = owner
		}

		return printJSON(out)
	},
}

var relayLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show relay log tail",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, logPath, err := relaymgr.TailLogs(relayLogsLines)
		if err != nil {
			return fmt.Errorf("unable to read relay logs (%s): %w", logPath, err)
		}
		out := map[string]interface{}{
			"logPath": logPath,
			"lines":   relayLogsLines,
			"content": content,
		}
		return printJSON(out)
	},
}

var relayInstallAgentCmd = &cobra.Command{
	Use:   "install-agent",
	Short: "Install macOS launch agent for persistent relay (optional)",
	Long: `Installs a user launch agent with KeepAlive so relay auto-starts at login.
This is optional and never enabled automatically.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Load()
		result, err := relaymgr.InstallLaunchAgent(cfg.ServerHost, cfg.Port)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func ensureRelayIfNeeded() error {
	cfg := config.Load()
	result, err := relaymgr.Ensure(relaymgr.EnsureOptions{
		Host:      cfg.ServerHost,
		Port:      cfg.Port,
		AutoStart: !noAutoRelay,
	})
	if err != nil {
		return err
	}
	if result.Started {
		fmt.Fprintf(os.Stderr, "[relay] started local relay on %s:%d\n", cfg.ServerHost, cfg.Port)
	}
	return nil
}

func newConnectedClient(channelKey string, live bool) (*ws.Client, error) {
	cfg := config.Load()
	url := fmt.Sprintf("ws://%s:%d", cfg.ServerHost, cfg.Port)
	client := ws.NewClient(url)

	if live {
		client.SetProgressHandler(func(msg *ws.Message) {
			out := map[string]interface{}{
				"type": "progress_update",
			}
			if len(msg.Result) > 0 {
				var parsed interface{}
				if err := json.Unmarshal(msg.Result, &parsed); err == nil {
					out["progress"] = parsed
				} else {
					out["progressRaw"] = string(msg.Result)
				}
			}
			_ = printJSON(out)
		})
	}

	if err := client.Connect(channelKey); err != nil {
		return nil, err
	}
	return client, nil
}

func printJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func loadBatchOperations(operationsJSON, operationsFile string) ([]batchOperation, error) {
	if operationsJSON == "" && operationsFile == "" {
		return nil, fmt.Errorf("either --operations or --operations-file is required")
	}
	if operationsJSON != "" && operationsFile != "" {
		return nil, fmt.Errorf("use only one of --operations or --operations-file")
	}

	var raw []byte
	var err error
	if operationsFile != "" {
		raw, err = os.ReadFile(operationsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read operations file: %w", err)
		}
	} else {
		raw = []byte(operationsJSON)
	}

	var ops []batchOperation
	if err := json.Unmarshal(raw, &ops); err != nil {
		return nil, fmt.Errorf("invalid operations JSON: %w", err)
	}
	return ops, nil
}

type relayStatus struct {
	Status           string         `json:"status"`
	Channels         map[string]int `json:"channels"`
	PreferredChannel string         `json:"preferredChannel"`
}

func resolveChannel(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	if envChannel := os.Getenv("AHD_CHANNEL"); envChannel != "" {
		return envChannel, nil
	}

	cfg := config.Load()
	statusURL := fmt.Sprintf("http://%s:%d/status", cfg.ServerHost, cfg.Port)
	httpClient := &http.Client{Timeout: 2 * time.Second}
	resp, err := httpClient.Get(statusURL)
	if err != nil {
		return "", fmt.Errorf("unable to resolve channel automatically: relay status request failed (%v). Pass channel explicitly or set AHD_CHANNEL", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to resolve channel automatically: relay status endpoint returned %d", resp.StatusCode)
	}

	var status relayStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", fmt.Errorf("unable to parse relay status response: %w", err)
	}

	active := make([]string, 0, len(status.Channels))
	for channel, count := range status.Channels {
		if count > 0 {
			active = append(active, channel)
		}
	}
	sort.Strings(active)

	if len(active) == 0 {
		return "", fmt.Errorf("no active plugin channel found. Open the plugin and connect, or pass a channel key")
	}

	if status.PreferredChannel != "" {
		for _, channel := range active {
			if channel == status.PreferredChannel {
				return channel, nil
			}
		}
	}

	if len(active) == 1 {
		return active[0], nil
	}

	return "", fmt.Errorf(
		"multiple active channels found: %s. Pass a channel key or set AHD_CHANNEL",
		strings.Join(active, ", "),
	)
}

func main() {
	rootCmd.PersistentFlags().BoolVar(&noAutoRelay, "no-auto-relay", false, "Disable CLI auto-start of local relay for connect/command/batch")

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		switch cmd.Name() {
		case "upgrade", "mcp":
			return
		}
		notifyUpdateAvailable()
	}

	setupCmd.Flags().StringVar(&setupPath, "path", "", "Custom directory to extract plugin into")
	setupCmd.Flags().BoolVar(&setupForce, "force", false, "Force re-extraction even if up to date")

	commandCmd.Flags().StringVarP(&commandParams, "params", "p", "", "JSON object passed as command params")
	commandCmd.Flags().BoolVar(&commandLive, "live", false, "Print live progress events while command is running")
	commandCmd.Flags().StringVar(&commandChannel, "channel", "", "Channel key override (optional)")

	batchCmd.Flags().StringVarP(&batchOperations, "operations", "o", "", "JSON array of operations")
	batchCmd.Flags().StringVarP(&batchOperationsFile, "operations-file", "f", "", "Path to JSON file containing operations array")
	batchCmd.Flags().BoolVar(&batchLive, "live", false, "Print live progress events while batch is running")
	batchCmd.Flags().StringVar(&batchChannel, "channel", "", "Channel key override (optional)")
	batchCmd.Flags().BoolVar(&batchFailFast, "fail-fast", false, "Stop at first failed operation")
	batchCmd.Flags().IntVar(&batchRetries, "retries", 1, "Retry count per operation after first attempt")
	batchCmd.Flags().IntVar(&batchRetryDelayMs, "retry-delay-ms", 250, "Delay between retries in milliseconds")
	batchCmd.Flags().BoolVar(&batchInterpolation, "interpolate", true, "Enable placeholder interpolation from prior step results")

	toolsCmd.Flags().BoolVar(&catalogJSON, "json", true, "Output as JSON for machine-readable discovery")
	toolsCmd.Flags().BoolVar(&catalogLLM, "llm", false, "Output enriched LLM-focused catalog with examples and playbook")
	relayLogsCmd.Flags().IntVar(&relayLogsLines, "lines", 80, "Number of log lines to show")
	actionsCmd.Flags().BoolVar(&actionsJSON, "json", false, "Output as JSON for machine-readable discovery")

	relayCmd.AddCommand(relayStartCmd)
	relayCmd.AddCommand(relayStopCmd)
	relayCmd.AddCommand(relayStatusCmd)
	relayCmd.AddCommand(relayLogsCmd)
	relayCmd.AddCommand(relayInstallAgentCmd)

	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(wsCmd)
	rootCmd.AddCommand(commandCmd)
	rootCmd.AddCommand(batchCmd)
	rootCmd.AddCommand(toolsCmd)
	rootCmd.AddCommand(actionsCmd)
	rootCmd.AddCommand(relayCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(upgradeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
