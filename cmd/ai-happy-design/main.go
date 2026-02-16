package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/nerveband/ai-happy-design/internal/imgutil"
	"github.com/nerveband/ai-happy-design/internal/mcp"
	pluginpkg "github.com/nerveband/ai-happy-design/internal/plugin"
	relaymgr "github.com/nerveband/ai-happy-design/internal/relay"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

var version = "0.0.0-dev"

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
background for communicating with the Figma plugin. If a relay is already
running on the configured port, it connects as a client instead.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadConfig()

		// Try to start embedded relay, but if port is in use, connect as client
		wsServer := ws.NewServer(cfg.Port)
		go func() {
			err := wsServer.Start()
			if err != nil {
				// Port likely in use by an existing relay — connect as client instead
				url := fmt.Sprintf("ws://%s:%d/ws", cfg.ServerHost, cfg.Port)
				log.Printf("[mcp] relay port %d in use, connecting as client to %s", cfg.Port, url)
				if clientErr := wsServer.ConnectAsClient(url); clientErr != nil {
					log.Printf("[mcp] client connect failed: %v", clientErr)
				}
			}
		}()

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
		cfg := loadConfig()
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
		cfg := loadConfig()

		// Auto-kill stale relay process holding the port
		if inUse, _ := relaymgr.IsPortInUse(cfg.Port); inUse {
			owner, isAHD := relaymgr.PortOwner(cfg.Port)
			if isAHD {
				fmt.Fprintf(os.Stderr, "Port %d held by stale relay (%s) — killing it.\n", cfg.Port, owner)
				_, _ = relaymgr.Stop()
				// Also force-kill by port in case state file was stale
				if killCmd, err := exec.LookPath("lsof"); err == nil {
					out, _ := exec.Command(killCmd, "-ti", fmt.Sprintf(":%d", cfg.Port)).Output()
					for _, pidStr := range strings.Fields(strings.TrimSpace(string(out))) {
						if pid, err := strconv.Atoi(pidStr); err == nil && pid != os.Getpid() {
							if p, err := os.FindProcess(pid); err == nil {
								_ = p.Kill()
							}
						}
					}
				}
				time.Sleep(500 * time.Millisecond)
			} else if owner != "" {
				return fmt.Errorf("port %d is already in use by %s (not ai-happy-design)", cfg.Port, owner)
			}
		}

		wsServer := ws.NewServer(cfg.Port)
		wsServer.SetIdleTimeout(cfg.IdleTimeout)
		if cfg.IdleTimeout > 0 {
			fmt.Printf("Starting WebSocket relay on port %d (idle shutdown: %s)...\n", cfg.Port, cfg.IdleTimeout)
		} else {
			fmt.Printf("Starting WebSocket relay on port %d...\n", cfg.Port)
		}
		if cfg.Port != config.DefaultPort {
			fmt.Fprintf(os.Stderr, "\n⚠  Non-default port %d. Update the Figma plugin:\n", cfg.Port)
			fmt.Fprintf(os.Stderr, "   1. Open the plugin in Figma\n")
			fmt.Fprintf(os.Stderr, "   2. Change the Relay URL to: ws://localhost:%d/ws  (or just type %d)\n", cfg.Port, cfg.Port)
			fmt.Fprintf(os.Stderr, "   3. Click Connect\n")
			fmt.Fprintf(os.Stderr, "   CLI commands: ai-happy-design --port %d command ...\n\n", cfg.Port)
		}
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
var commandOutput string
var commandBase64 bool
var commandCompressImages bool

var commandCmd = &cobra.Command{
	Use:   "command <command> [params-json]",
	Short: "Execute a command against a connected Figma channel",
	Long: `Executes one command through the WebSocket relay for scripting and LLM-driven CLI flows.

Usage:
  command node.create_frame '{"x":0,"y":0,"width":1080,"height":1350}'
  command document.find_free_space
  command node.create_frame -p '{"x":0,"y":0}'

Command supports both legacy names (e.g. set_fill_color) and domain.action (e.g. paint.set_solid).
Params can be passed as a second positional arg (JSON) or via -p/--params flag.
Channel resolution: --channel flag, AHD_CHANNEL env, relay preferred/active channel.
If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.RangeArgs(1, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)
		var command string
		command = args[0]

		// Parse remaining positional args: could be params JSON, channel, or both
		for _, arg := range args[1:] {
			trimmed := strings.TrimSpace(arg)
			if strings.HasPrefix(trimmed, "{") {
				// It's JSON params
				if commandParams != "" {
					return fmt.Errorf("provide params as positional arg OR --params, not both")
				}
				commandParams = trimmed
			} else {
				// It's a channel key
				if commandChannel != "" {
					return fmt.Errorf("provide channel either as positional argument or --channel, not both")
				}
				commandChannel = trimmed
			}
		}

		params := map[string]interface{}{}
		if commandParams != "" {
			if err := json.Unmarshal([]byte(commandParams), &params); err != nil {
				return fmt.Errorf("invalid params JSON: %w", err)
			}
		}
		params = batchutil.NormalizeBatchParams(command, params)

		// Handle local-only commands that don't need Figma connection
		if handled, err := handleLocalCommand(command, params); handled {
			return err
		}

		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		channelKey, err := resolveChannel(commandChannel)
		if err != nil {
			return err
		}

		// Opt-in image compression
		if commandCompressImages {
			if compressed, changed, compErr := imgutil.CompressParamsImageData(params, imgutil.DefaultOptions()); compErr != nil {
				fmt.Fprintf(os.Stderr, "[compress] warning: %v\n", compErr)
			} else if changed {
				params = compressed
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

		// Check if response has base64 "data" field (e.g. export commands)
		var parsed map[string]interface{}
		hasData := false
		if err := json.Unmarshal(result, &parsed); err == nil {
			_, hasData = parsed["data"].(string)
		}

		// For responses with data field (exports): default to file save, --base64 for raw
		if hasData && !commandBase64 {
			dataStr := parsed["data"].(string)
			format, _ := parsed["format"].(string)

			// SVG and JSON exports return raw text; PNG/JPG/PDF return base64
			var fileBytes []byte
			switch strings.ToUpper(format) {
			case "SVG", "JSON":
				fileBytes = []byte(dataStr)
			default:
				var decErr error
				fileBytes, decErr = base64.StdEncoding.DecodeString(dataStr)
				if decErr != nil {
					return fmt.Errorf("failed to decode base64 data: %w", decErr)
				}
			}

			outPath := commandOutput
			if outPath == "" {
				// Auto-generate filename from node name and format
				name, _ := parsed["name"].(string)
				if name == "" {
					name = "export"
				}
				// Sanitize name for filename
				name = strings.Map(func(r rune) rune {
					if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' {
						return '-'
					}
					return r
				}, name)
				ext := ".png"
				switch strings.ToUpper(format) {
				case "JPG":
					ext = ".jpg"
				case "SVG":
					ext = ".svg"
				case "PDF":
					ext = ".pdf"
				case "JSON":
					ext = ".json"
				}
				outPath = name + ext
			}

			if err := os.WriteFile(outPath, fileBytes, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			// Print metadata without the bulky data field
			delete(parsed, "data")
			parsed["savedTo"] = outPath
			parsed["fileSize"] = len(fileBytes)
			return printJSON(parsed)
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
var batchCompressImages bool
var batchAllowOverlap bool

type batchOperation struct {
	Name    string                 `json:"name,omitempty"`
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

var batchCmd = &cobra.Command{
	Use:   "batch [operations-json | operations-file]",
	Short: "Execute a sequence of commands in order",
	Long: `Runs multiple commands over one connection to support tool chaining workflows.

Operations can be provided in multiple ways (in priority order):
  1. Positional arg (JSON):  batch '[{"command":"node.create_frame","params":{...}}]'
  2. Positional arg (file):  batch ops.json
  3. Flag (inline JSON):     batch -o '[...]'
  4. Flag (file path):       batch -f ops.json
  5. Stdin (piped):          cat ops.json | batch

	Supports interpolation placeholders like ${{steps.0.result.id}} or ${{steps.createCard.result.id}}.
	Short interpolation also works: $createCard (id), $createCard.width, $last.
	Compact aliases are supported (examples: frame, rect, text, fill, stroke, parent).
	Command-aware shorthand params are supported in batch/bulk (examples: w/h/pid, sz/ff/lh/ls, sw, bg).
	Channel resolution: --channel flag, AHD_CHANNEL env, relay preferred/active channel.
	If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)
		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		// Positional arg: if starts with '[' it's inline JSON, otherwise a file path
		if len(args) == 1 {
			arg := strings.TrimSpace(args[0])
			if strings.HasPrefix(arg, "[") {
				if batchOperations != "" {
					return fmt.Errorf("provide operations as positional arg OR --operations, not both")
				}
				batchOperations = arg
			} else {
				if batchOperationsFile != "" {
					return fmt.Errorf("provide operations file as positional arg OR --operations-file, not both")
				}
				batchOperationsFile = arg
			}
		}

		channelKey, err := resolveChannel(batchChannel)
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
		batchStart := time.Now()

		// Auto-place root frames to avoid overlap (unless --allow-overlap)
		if !batchAllowOverlap {
			autoPlaceRootFrames(client, ops)
		}

		for i, op := range ops {
			opStart := time.Now()
			op.Name = batchutil.SanitizeStepName(op.Name)
			params := batchutil.NormalizeBatchParams(op.Command, op.Params)

			if batchInterpolation {
				interpolatedParams, err := batchutil.InterpolateParams(params, stepStates)
				if err != nil {
					failed++
					entry := map[string]interface{}{
						"index":     i,
						"name":      op.Name,
						"command":   op.Command,
						"ok":        false,
						"error":     fmt.Sprintf("interpolation error: %v", err),
						"attempts":  0,
						"elapsedMs": int(time.Since(opStart).Milliseconds()),
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

			// Opt-in image compression
			if batchCompressImages {
				if compressed, changed, compErr := imgutil.CompressParamsImageData(params, imgutil.DefaultOptions()); compErr != nil {
					fmt.Fprintf(os.Stderr, "[compress] step %d warning: %v\n", i, compErr)
				} else if changed {
					params = compressed
				}
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
					"index":     i,
					"name":      op.Name,
					"command":   op.Command,
					"ok":        false,
					"attempts":  attempts,
					"error":     sendErr.Error(),
					"elapsedMs": int(time.Since(opStart).Milliseconds()),
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
				"index":     i,
				"name":      op.Name,
				"command":   op.Command,
				"ok":        true,
				"attempts":  attempts,
				"result":    parsed,
				"elapsedMs": int(time.Since(opStart).Milliseconds()),
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
		totalElapsed := time.Since(batchStart)
		totalMs := int(totalElapsed.Milliseconds())
		opsPerSec := float64(0)
		avgMs := float64(0)
		if processed > 0 {
			opsPerSec = float64(processed) / totalElapsed.Seconds()
			avgMs = float64(totalMs) / float64(processed)
		}
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
			"timing": map[string]interface{}{
				"totalMs":   totalMs,
				"avgMs":     int(avgMs),
				"opsPerSec": int(opsPerSec),
			},
			"stoppedEarly": stoppedEarly,
			"steps":        results,
		}
		return printJSON(out)
	},
}

var catalogJSON bool
var catalogLLM bool
var noAutoRelay bool
var globalPort int

func loadConfig() *config.Config {
	cfg := config.Load()
	if globalPort > 0 {
		cfg.Port = globalPort
	}
	return cfg
}
var relayLogsLines int
var actionsJSON bool

var actionCatalog = map[string][]string{
	"document":  {"get_info", "get_selection", "set_selection", "scan_text", "scan_by_type", "get_styles", "find_by_name", "find_by_type", "focus", "zoom_to", "find_free_space"},
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
	"shape":     {"create_rectangle", "create_ellipse", "create_polygon", "create_star", "create_line", "create_from_svg", "create_image"},
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
		cfg := loadConfig()
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
		cfg := loadConfig()
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
		cfg := loadConfig()
		result, err := relaymgr.InstallLaunchAgent(cfg.ServerHost, cfg.Port)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func ensureRelayIfNeeded() error {
	cfg := loadConfig()
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

// handleLocalCommand handles commands that can be processed locally without a Figma connection.
// Returns (true, err) if handled, (false, nil) if should be routed to Figma.
func handleLocalCommand(command string, params map[string]interface{}) (bool, error) {
	switch command {
	case "design.compute_tokens":
		w, _ := params["width"].(float64)
		h, _ := params["height"].(float64)
		if w <= 0 || h <= 0 {
			return true, fmt.Errorf("width and height must be positive numbers")
		}
		dpi, _ := params["dpi"].(float64)
		tokens := tools.ComputeDesignTokens(w, h, dpi)
		out, _ := json.MarshalIndent(tokens, "", "  ")
		fmt.Println(string(out))
		return true, nil
	default:
		return false, nil
	}
}

func newConnectedClient(channelKey string, live bool) (*ws.Client, error) {
	cfg := loadConfig()
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
		hint := ""
		if cfg.Port != config.DefaultPort {
			hint = fmt.Sprintf("\nNote: using non-default port %d. Make sure the relay is running on this port:\n  ai-happy-design --port %d ws\nAnd update the Figma plugin Relay URL to: ws://localhost:%d/ws", cfg.Port, cfg.Port, cfg.Port)
		} else {
			hint = "\nTip: if the relay is on a different port, use --port <port> or set PORT env var."
		}
		return nil, fmt.Errorf("%w%s", err, hint)
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

// autoPlaceRootFrames checks if the batch has root-level frame creation (no parentId)
// without a find_free_space step. If so, it calls find_free_space via the relay and
// offsets the first root frame's x/y to a safe position. Subsequent root frames are
// placed relative to the offset. This prevents accidentally drawing on top of existing work.
func autoPlaceRootFrames(client *ws.Client, ops []batchOperation) {
	// Check if batch already has a find_free_space step
	for _, op := range ops {
		if strings.Contains(strings.ToLower(op.Command), "find_free_space") {
			return // User is already handling placement
		}
	}

	// Find the first root-level frame creation
	firstRootIdx := -1
	var rootW, rootH float64
	for i, op := range ops {
		cmd := strings.ToLower(op.Command)
		if cmd != "frame" && cmd != "node.create_frame" {
			continue
		}
		pid, _ := op.Params["parentId"].(string)
		if pid == "" {
			pid, _ = op.Params["pid"].(string)
		}
		if pid != "" {
			continue
		}
		firstRootIdx = i
		rootW = getFloatParam(op.Params, "width", getFloatParam(op.Params, "w", 0))
		rootH = getFloatParam(op.Params, "height", getFloatParam(op.Params, "h", 0))
		break
	}

	if firstRootIdx < 0 || (rootW == 0 && rootH == 0) {
		return // No root frame or no dimensions to check
	}

	// Call find_free_space to get safe coordinates
	findParams := map[string]interface{}{"width": rootW, "height": rootH}
	result, err := client.SendCommand("document.find_free_space", findParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠ Could not check free space: %v (proceeding without offset)\n", err)
		return
	}

	var freeSpace struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		ExistingCount  int     `json:"existingCount"`
	}
	if err := json.Unmarshal(result, &freeSpace); err != nil {
		return
	}

	if freeSpace.ExistingCount == 0 {
		return // Empty page, no need to offset
	}

	// Calculate offset from original position to safe position
	origX := getFloatParam(ops[firstRootIdx].Params, "x", 0)
	origY := getFloatParam(ops[firstRootIdx].Params, "y", 0)
	dx := freeSpace.X - origX
	dy := freeSpace.Y - origY

	if dx == 0 && dy == 0 {
		return // Already in a safe spot
	}

	fmt.Fprintf(os.Stderr, "⚠ Auto-placing: shifted root frames by (%+.0f, %+.0f) to avoid overlap with %d existing frame(s).\n", dx, dy, freeSpace.ExistingCount)
	fmt.Fprintf(os.Stderr, "  Use --allow-overlap to place at exact coordinates.\n")

	// Offset all root-level frame creations
	for i := range ops {
		cmd := strings.ToLower(ops[i].Command)
		if cmd != "frame" && cmd != "node.create_frame" {
			continue
		}
		pid, _ := ops[i].Params["parentId"].(string)
		if pid == "" {
			pid, _ = ops[i].Params["pid"].(string)
		}
		if pid != "" {
			continue
		}
		ops[i].Params["x"] = getFloatParam(ops[i].Params, "x", 0) + dx
		ops[i].Params["y"] = getFloatParam(ops[i].Params, "y", 0) + dy
	}
}

// getFloatParam extracts a numeric value from params, checking both the key and common aliases.
func getFloatParam(params map[string]interface{}, key string, def float64) float64 {
	if v, ok := params[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		case json.Number:
			f, _ := n.Float64()
			return f
		}
	}
	return def
}

func loadBatchOperations(operationsJSON, operationsFile string) ([]batchOperation, error) {
	if operationsJSON != "" && operationsFile != "" {
		return nil, fmt.Errorf("use only one of --operations or --operations-file")
	}

	var raw []byte
	var err error
	switch {
	case operationsFile != "":
		raw, err = os.ReadFile(operationsFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read operations file: %w", err)
		}
	case operationsJSON != "":
		raw = []byte(operationsJSON)
	default:
		// Try stdin (for piped input)
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			raw, err = io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("failed to read stdin: %w", err)
			}
		} else {
			return nil, fmt.Errorf("no operations provided. Usage:\n  batch '[{\"command\":\"node.create_frame\",\"params\":{...}}]'\n  batch ops.json\n  batch -f ops.json\n  cat ops.json | batch")
		}
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

	cfg := loadConfig()
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
		hint := "no active plugin channel found. Open the Figma plugin and click Connect."
		if cfg.Port != config.DefaultPort {
			hint += fmt.Sprintf("\n\nYou're using port %d (non-default). In the Figma plugin, set Relay URL to:\n  ws://localhost:%d/ws  (or just type: %d)\nThen click Connect.", cfg.Port, cfg.Port, cfg.Port)
		}
		return "", fmt.Errorf("%s", hint)
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
	rootCmd.PersistentFlags().IntVar(&globalPort, "port", 0, "Override relay port (default 3055, or PORT env var)")

	rootCmd.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		switch cmd.Name() {
		case "upgrade", "mcp", "command", "batch", "ws", "start", "stop", "status", "logs":
			return
		}
		notifyUpdateAvailable()
	}

	setupCmd.Flags().StringVar(&setupPath, "path", "", "Custom directory to extract plugin into")
	setupCmd.Flags().BoolVar(&setupForce, "force", false, "Force re-extraction even if up to date")

	commandCmd.Flags().StringVarP(&commandParams, "params", "p", "", "JSON object passed as command params")
	commandCmd.Flags().BoolVar(&commandLive, "live", false, "Print live progress events while command is running")
	commandCmd.Flags().StringVar(&commandChannel, "channel", "", "Channel key override (optional)")
	commandCmd.Flags().StringVarP(&commandOutput, "output", "o", "", "Save binary data (e.g. exported image) to this file path")
	commandCmd.Flags().BoolVar(&commandBase64, "base64", false, "Output raw base64 JSON instead of saving to file (for export commands)")
	commandCmd.Flags().BoolVar(&commandCompressImages, "compress-images", false, "Compress imageData before sending (requires ImageMagick)")

	batchCmd.Flags().StringVarP(&batchOperations, "operations", "o", "", "JSON array of operations")
	batchCmd.Flags().StringVarP(&batchOperationsFile, "operations-file", "f", "", "Path to JSON file containing operations array")
	batchCmd.Flags().BoolVar(&batchLive, "live", false, "Print live progress events while batch is running")
	batchCmd.Flags().StringVar(&batchChannel, "channel", "", "Channel key override (optional)")
	batchCmd.Flags().BoolVar(&batchFailFast, "fail-fast", false, "Stop at first failed operation")
	batchCmd.Flags().IntVar(&batchRetries, "retries", 1, "Retry count per operation after first attempt")
	batchCmd.Flags().IntVar(&batchRetryDelayMs, "retry-delay-ms", 250, "Delay between retries in milliseconds")
	batchCmd.Flags().BoolVar(&batchInterpolation, "interpolate", true, "Enable placeholder interpolation from prior step results")
	batchCmd.Flags().BoolVar(&batchCompressImages, "compress-images", false, "Compress imageData before sending (requires ImageMagick)")
	batchCmd.Flags().BoolVar(&batchAllowOverlap, "allow-overlap", false, "Skip auto-placement and place frames at exact coordinates (may overlap existing work)")

	toolsCmd.Flags().BoolVar(&catalogJSON, "json", true, "Output as JSON for machine-readable discovery")
	toolsCmd.Flags().BoolVar(&catalogLLM, "llm", false, "Output enriched LLM-focused catalog with examples and playbook")
	relayLogsCmd.Flags().IntVar(&relayLogsLines, "lines", 80, "Number of log lines to show")
	actionsCmd.Flags().BoolVar(&actionsJSON, "json", false, "Output as JSON for machine-readable discovery")

	relayCmd.AddCommand(relayStartCmd)
	relayCmd.AddCommand(relayStopCmd)
	relayCmd.AddCommand(relayStatusCmd)
	relayCmd.AddCommand(relayLogsCmd)
	relayCmd.AddCommand(relayInstallAgentCmd)

	registerCmd.Flags().StringVar(&registerEditor, "editor", "", "Register with a specific editor only (e.g. 'Claude Code', 'Cursor')")
	registerCmd.Flags().BoolVar(&registerForce, "force", false, "Force re-registration even if already configured")

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
	rootCmd.AddCommand(registerCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
