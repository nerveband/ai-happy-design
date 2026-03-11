package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
	"github.com/nerveband/ai-happy-design/internal/benchmark"
	"github.com/nerveband/ai-happy-design/internal/config"
	"github.com/nerveband/ai-happy-design/internal/designlint"
	"github.com/nerveband/ai-happy-design/internal/extract"
	"github.com/nerveband/ai-happy-design/internal/imgutil"
	pluginpkg "github.com/nerveband/ai-happy-design/internal/plugin"
	relaymgr "github.com/nerveband/ai-happy-design/internal/relay"
	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/nerveband/ai-happy-design/internal/validate"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

var version = "0.0.0-dev"

var rootCmd = &cobra.Command{
	Use:   "ahd-figma",
	Short: "AHD Figma - Figma MCP server and CLI",
	Long: `A single binary that works as both an MCP server for AI editors
and a CLI for direct Figma manipulation via a WebSocket relay.

Discovery-first flow for LLM agents:
  1) ahd-figma tools --llm --json
  2) ahd-figma actions [domain]
  3) ahd-figma batch --help`,
	Version:       version,
	SilenceErrors: true,
	SilenceUsage:  true,
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
				return fmt.Errorf("port %d is already in use by %s (not ahd-figma)", cfg.Port, owner)
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
			fmt.Fprintf(os.Stderr, "   CLI commands: ahd-figma --port %d command ...\n\n", cfg.Port)
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
		fmt.Println("  4. Run: Plugins > AHD Figma")

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
For image commands, imageData supports base64/data URLs plus file paths and HTTP(S) URLs (auto-resolved).
Use --compress-images to reduce transfer size for large image payloads.
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

		// Resolve file:// and path references in imageData
		if resolved, changed, resolveErr := imgutil.ResolveParamsImageData(params); resolveErr != nil {
			return fmt.Errorf("imageData resolution failed: %w", resolveErr)
		} else if changed {
			params = resolved
		}

		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		channelKey, err := resolveChannel(commandChannel)
		if err != nil {
			return err
		}

		// Check plugin connectivity before sending command
		if err := checkPluginConnected(channelKey); err != nil {
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
					if r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|' || r == ' ' {
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
				outPath = fmt.Sprintf("%s/ahd-export-%s-%d%s", os.TempDir(), name, time.Now().Unix(), ext)
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
var batchNoFix bool

type batchOperation struct {
	Name    string                 `json:"name,omitempty"`
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

var batchParallel bool
var batchCompact bool
var batchLint bool
var batchNoLint bool
var batchStrictQuality bool

var batchCmd = &cobra.Command{
	Use:   "batch [operations-json | file1.json file2.json ... | directory/]",
	Short: "Execute a sequence of commands in order",
	Long: `Runs multiple commands over one connection to support tool chaining workflows.

LLM discovery quickstart:
  1) ahd-figma tools --llm --json
  2) ahd-figma actions [domain]
  3) ahd-figma batch --help

Operations can be provided in multiple ways (in priority order):
  1. Positional arg (JSON):  batch '[{"command":"node.create_frame","params":{...}}]'
  2. Positional arg (file):  batch ops.json
  3. Multiple files:         batch file1.json file2.json file3.json
  4. Directory:              batch ./carousels/  (runs all .json files)
  5. Glob:                   batch ops/*.json
  6. Flag (inline JSON):     batch -o '[...]'
  7. Flag (file path):       batch -f ops.json
  8. Stdin (piped):          cat ops.json | batch

Use --parallel to run multiple batch files concurrently (max 4).
Each file gets its own connection and auto-placement.

Image prep behavior (automatic):
  - imageData accepts base64/data URLs, file paths, and HTTP(S) URLs.
  - Batch resolves and preps unique image payloads in parallel before execution.
  - Repeated sources are de-duplicated to avoid rework.
  - --compress-images applies optional ImageMagick compression.
  - Progress is reported to stderr with [image-prep] lines.

Quality + telemetry:
  - --lint is enabled by default and checks overlaps/overflow/default names/text sizing.
  - --strict-quality fails the run if lint reports any warning/error issue.
  - JSON output always includes summary, timing, and imagePrep for agent feedback loops.

Supports interpolation placeholders like ${{steps.0.result.id}} or ${{steps.createCard.result.id}}.
Short interpolation also works: $createCard (id), $createCard.width, $last.
Compact aliases are supported (examples: frame, rect, text, fill, stroke, parent).
Command-aware shorthand params are supported in batch/bulk (examples: w/h/pid, sz/ff/lh/ls, sw, bg).
Channel resolution: --channel flag, AHD_CHANNEL env, relay preferred/active channel.
If relay is not running, CLI auto-starts it unless --no-auto-relay is set.`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)
		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}
		if batchStrictQuality && (!batchLint || batchNoLint) {
			return fmt.Errorf("--strict-quality requires lint checks. Remove --no-lint (or set --lint=true)")
		}

		// Collect batch files: multiple positional args, directory, inline JSON, or single file
		batchFiles, inlineJSON, err := collectBatchInputs(args, batchOperations, batchOperationsFile)
		if err != nil {
			return err
		}

		// Multi-file mode
		if len(batchFiles) > 1 || (len(batchFiles) == 1 && inlineJSON == "") {
			channelKey, chErr := resolveChannel(batchChannel)
			if chErr != nil {
				return chErr
			}
			if err := checkPluginConnected(channelKey); err != nil {
				return err
			}
			return runMultiBatch(batchFiles, channelKey)
		}

		// Single batch mode (legacy path)
		if inlineJSON != "" {
			batchOperations = inlineJSON
			batchOperationsFile = ""
		} else if len(batchFiles) == 1 {
			batchOperationsFile = batchFiles[0]
			batchOperations = ""
		}

		channelKey, err := resolveChannel(batchChannel)
		if err != nil {
			return err
		}

		// Check plugin connectivity
		if err := checkPluginConnected(channelKey); err != nil {
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
		ops, imagePrep := preprocessBatchImageData(ops, batchCompressImages, "", true)

		// Schema validation (warn+fix mode)
		var schemaResult validate.Result
		if !batchNoFix {
			validationOps := make([]map[string]interface{}, len(ops))
			for i, op := range ops {
				validationOps[i] = map[string]interface{}{
					"command": op.Command,
					"params":  op.Params,
					"name":    op.Name,
				}
			}
			schemaResult = validate.ValidateBatch(validationOps)
			// Write fixed values back to ops
			for i, vo := range validationOps {
				if cmd, ok := vo["command"].(string); ok {
					ops[i].Command = cmd
				}
				if p, ok := vo["params"].(map[string]interface{}); ok {
					ops[i].Params = p
				}
			}
			if schemaResult.Fixed > 0 || len(schemaResult.Warnings) > 0 {
				fmt.Fprintf(os.Stderr, "[schema] %d fixed, %d warnings, %d blocked\n",
					schemaResult.Fixed, len(schemaResult.Warnings), schemaResult.Blocked)
			}
			if batchStrictQuality && schemaResult.Blocked > 0 {
				out := map[string]interface{}{
					"ok": false,
					"preValidation": map[string]interface{}{
						"schema": map[string]interface{}{
							"warnings": schemaResult.Warnings,
							"errors":   schemaResult.Errors,
							"fixed":    schemaResult.Fixed,
							"blocked":  schemaResult.Blocked,
						},
					},
				}
				printJSONErr(out)
				return fmt.Errorf("schema validation blocked %d issues", schemaResult.Blocked)
			}
		}

		// Design lint (pre-execution quality checks)
		var lintResult designlint.Result
		if batchLint && !batchNoLint {
			lintOps := make([]map[string]interface{}, len(ops))
			for i, op := range ops {
				lintOps[i] = map[string]interface{}{
					"command": op.Command,
					"params":  op.Params,
					"name":    op.Name,
				}
			}
			lintResult = designlint.Check(lintOps)
			// Write fixes back
			for i, lo := range lintOps {
				if p, ok := lo["params"].(map[string]interface{}); ok {
					ops[i].Params = p
				}
			}
			if lintResult.Fixed > 0 || len(lintResult.Warnings) > 0 {
				fmt.Fprintf(os.Stderr, "[design-lint] %d fixed, %d warnings, score: %.1f/10\n",
					lintResult.Fixed, len(lintResult.Warnings), lintResult.Score.Overall)
			}
			if batchStrictQuality && lintResult.Score.Overall < 7.0 {
				out := map[string]interface{}{
					"ok": false,
					"preValidation": map[string]interface{}{
						"designLint": map[string]interface{}{
							"canvas":   lintResult.Canvas,
							"warnings": lintResult.Warnings,
							"fixed":    lintResult.Fixed,
							"score":    lintResult.Score,
						},
					},
				}
				printJSONErr(out)
				return fmt.Errorf("design quality score %.1f/10 below threshold 7.0", lintResult.Score.Overall)
			}
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

			// Runtime fallback for interpolated imageData values that could not be preprocessed.
			if rawImage, ok := params["imageData"].(string); ok && needsRuntimeImagePrep(rawImage) {
				if resolved, changed, resolveErr := imgutil.ResolveParamsImageData(params); resolveErr != nil {
					fmt.Fprintf(os.Stderr, "[resolve] step %d warning: %v\n", i, resolveErr)
				} else if changed {
					params = resolved
				}
				if batchCompressImages {
					if compressed, changed, compErr := imgutil.CompressParamsImageData(params, imgutil.DefaultOptions()); compErr != nil {
						fmt.Fprintf(os.Stderr, "[compress] step %d warning: %v\n", i, compErr)
					} else if changed {
						params = compressed
					}
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
					"errorCode": classifyError(sendErr.Error()),
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
		summaryMap := map[string]interface{}{
			"total":         len(ops),
			"processed":     processed,
			"succeeded":     succeeded,
			"failed":        failed,
			"pending":       pending,
			"retriesUsed":   retriesUsed,
			"failFast":      batchFailFast,
			"interpolation": batchInterpolation,
		}
		timingMap := map[string]interface{}{
			"totalMs":     totalMs,
			"avgMs":       int(avgMs),
			"opsPerSec":   roundTo(opsPerSec, 2),
			"imagePrepMs": imagePrep.TotalMs,
		}
		var out map[string]interface{}
		if batchCompact {
			// Compact mode: minimal output for LLM token efficiency
			compactSteps := make([]map[string]interface{}, len(results))
			for j, r := range results {
				cs := map[string]interface{}{
					"name": r["name"],
					"ok":   r["ok"],
				}
				if r["ok"] == true {
					cs["result"] = r["result"]
				} else {
					cs["error"] = r["error"]
					cs["errorCode"] = r["errorCode"]
				}
				compactSteps[j] = cs
			}
			out = map[string]interface{}{
				"ok":        failed == 0 && pending == 0,
				"summary":   summaryMap,
				"timing":    timingMap,
				"imagePrep": imagePrepSummaryMap(imagePrep),
				"steps":     compactSteps,
			}
		} else {
			out = map[string]interface{}{
				"ok":           failed == 0 && pending == 0,
				"summary":      summaryMap,
				"timing":       timingMap,
				"imagePrep":    imagePrepSummaryMap(imagePrep),
				"stoppedEarly": stoppedEarly,
				"steps":        results,
			}
		}

		// Include pre-validation results (schema + design lint)
		preVal := map[string]interface{}{}
		if schemaResult.Fixed > 0 || len(schemaResult.Warnings) > 0 {
			preVal["schema"] = map[string]interface{}{
				"warnings": schemaResult.Warnings,
				"fixed":    schemaResult.Fixed,
				"blocked":  schemaResult.Blocked,
			}
		}
		if lintResult.Fixed > 0 || len(lintResult.Warnings) > 0 {
			preVal["designLint"] = map[string]interface{}{
				"canvas":   lintResult.Canvas,
				"tokens":   lintResult.Tokens,
				"warnings": lintResult.Warnings,
				"fixed":    lintResult.Fixed,
				"score":    lintResult.Score,
			}
		}
		if len(preVal) > 0 {
			out["preValidation"] = preVal
		}

		// Post-batch lint: check created root frames for design issues.
		var lintInfo lintSummary
		if batchLint && !batchNoLint && failed == 0 {
			rootNodeIDs := collectCreatedRootFrameIDs(ops, results)
			lintInfo = runPostBatchLint(client, rootNodeIDs, "")
			lintGuides := lintGuidance(lintInfo)
			lintSamples := lintSamplesForOutput(lintInfo.Samples, 5)
			if !batchCompact {
				lintMap := map[string]interface{}{
					"issues":   lintInfo.Issues,
					"warnings": lintInfo.Warnings,
					"errors":   lintInfo.Errors,
					"byType":   lintInfo.ByType,
				}
				if len(lintGuides) > 0 {
					lintMap["guidance"] = lintGuides
				}
				if len(lintSamples) > 0 {
					lintMap["samples"] = lintSamples
				}
				out["lint"] = lintMap
			}
			if batchStrictQuality && lintInfo.Issues > 0 {
				out["ok"] = false
				gate := map[string]interface{}{
					"mode":   "strict",
					"status": "failed",
					"issues": lintInfo.Issues,
				}
				if len(lintGuides) > 0 {
					gate["guidance"] = lintGuides
				}
				out["qualityGate"] = gate
				if summary, ok := out["summary"].(map[string]interface{}); ok {
					summary["qualityGate"] = "failed"
					summary["qualityIssues"] = lintInfo.Issues
					if len(lintGuides) > 0 {
						summary["qualityGuidance"] = lintGuides
					}
				}
				if err := printJSON(out); err != nil {
					return err
				}
				return fmt.Errorf("strict quality gate failed: %d lint warning/error issue(s)", lintInfo.Issues)
			}
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
	"document":      {"get_info", "get_selection", "set_selection", "scan_text", "scan_by_type", "get_styles", "find_by_name", "find_by_type", "focus", "zoom_to", "find_free_space"},
	"node":          {"get_info", "create_frame", "move", "resize", "rotate", "set_opacity", "set_blend_mode", "set_visibility", "set_locked", "rename", "delete", "clone", "set_corner_radius", "get_tree"},
	"layer":         {"set_order", "bring_forward", "send_backward", "bring_to_front", "send_to_back", "group", "ungroup", "move_to_parent", "insert_child"},
	"layout":        {"set_auto_layout", "set_padding", "set_spacing", "set_alignment", "set_sizing", "set_constraints", "set_layout_wrap", "set_wrap", "remove_auto_layout"},
	"export":        {"image", "svg", "pdf", "json"},
	"design_system": {"analyze"},
	"variable":      {"create", "get_all", "set_value", "bind", "unbind", "create_collection", "get_collections", "delete"},
	"style":         {"create_paint", "create_text", "create_effect", "apply", "get_all", "remove"},
	"component":     {"create", "create_instance", "create_set", "get_local", "get_remote", "get_overrides", "set_overrides", "detach_instance", "reset_instance", "swap_instance"},
	"boolean":       {"union", "subtract", "intersect", "exclude", "flatten"},
	"paint":         {"set_solid", "set_gradient", "set_image_fill", "set_image", "set_image_fill_from_url", "set_image_url", "add_fill", "remove_fill", "get_fills", "set_stroke"},
	"effect":        {"set_effects", "add_shadow", "add_blur", "apply_style", "remove", "remove_effect", "get_effects"},
	"page":          {"create", "delete", "rename", "set_current", "get_all", "get_current", "duplicate"},
	"shape":         {"create_rectangle", "create_ellipse", "create_polygon", "create_star", "create_line", "create_from_svg", "create_image"},
	"text":          {"create", "set_content", "set_font", "set_size", "set_weight", "set_color", "set_align", "set_spacing", "set_line_height", "set_letter_spacing", "set_decoration", "set_case", "set_paragraph_spacing", "get_content", "get_segments", "load_font", "set_style_id"},
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
Use --llm --json for enriched examples and recommended execution playbook.

Suggested discovery flow:
  1) ahd-figma tools --llm --json
  2) ahd-figma actions [domain]
  3) ahd-figma batch --help`,
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
		return true, printJSON(tokens)
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
			hint = fmt.Sprintf("\nNote: using non-default port %d. Make sure the relay is running on this port:\n  ahd-figma --port %d ws\nAnd update the Figma plugin Relay URL to: ws://localhost:%d/ws", cfg.Port, cfg.Port, cfg.Port)
		} else {
			hint = "\nTip: if the relay is on a different port, use --port <port> or set PORT env var."
		}
		return nil, fmt.Errorf("%w%s", err, hint)
	}
	return client, nil
}

// classifyError returns a machine-readable error code from an error message.
func classifyError(msg string) string {
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "not found") || strings.Contains(lower, "node with id"):
		return "NODE_NOT_FOUND"
	case strings.Contains(lower, "timed out") || strings.Contains(lower, "timeout"):
		return "TIMEOUT"
	case strings.Contains(lower, "no active plugin") || strings.Contains(lower, "no figma plugin"):
		return "NO_PLUGIN"
	case strings.Contains(lower, "cannot be resized") || strings.Contains(lower, "does not support"):
		return "UNSUPPORTED_OPERATION"
	case strings.Contains(lower, "unknown") && (strings.Contains(lower, "action") || strings.Contains(lower, "command")):
		return "UNKNOWN_COMMAND"
	case strings.Contains(lower, "invalid") || strings.Contains(lower, "parse"):
		return "INVALID_PARAMS"
	case strings.Contains(lower, "font") && (strings.Contains(lower, "not available") || strings.Contains(lower, "missing")):
		return "FONT_NOT_FOUND"
	case strings.Contains(lower, "failed to send") || strings.Contains(lower, "connection"):
		return "CONNECTION_ERROR"
	case strings.Contains(lower, "fill index") || strings.Contains(lower, "out of range"):
		return "INDEX_OUT_OF_RANGE"
	default:
		return "UNKNOWN_ERROR"
	}
}

type imagePrepSummary struct {
	Candidates  int
	Unique      int
	CacheHits   int
	Prepared    int
	Changed     int
	Failed      int
	Resolved    int
	Compressed  int
	InputBytes  int
	OutputBytes int
	TotalMs     int
	AvgMs       int
	SlowestMs   int
}

type imagePrepTaskResult struct {
	Value       string
	Changed     bool
	Resolved    bool
	Compressed  bool
	Err         error
	ElapsedMs   int
	InputBytes  int
	OutputBytes int
}

func hasInterpolationToken(raw string) bool {
	return strings.Contains(raw, "${{") || strings.Contains(raw, "$steps.") || strings.Contains(raw, "$last")
}

func needsRuntimeImagePrep(raw string) bool {
	s := strings.TrimSpace(raw)
	if s == "" {
		return false
	}
	if hasInterpolationToken(s) {
		return true
	}
	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, "~/") || strings.HasPrefix(s, "file://") || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return true
	}
	return false
}

func shortImageSource(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "imageData"
	}
	if strings.HasPrefix(s, "data:") {
		return "data-uri"
	}
	if strings.HasPrefix(s, "file://") {
		s = strings.TrimPrefix(s, "file://")
	}
	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, "~/") {
		return filepath.Base(s)
	}
	if len(s) > 96 {
		return s[:93] + "..."
	}
	return s
}

func imagePrepSummaryMap(summary imagePrepSummary) map[string]interface{} {
	return map[string]interface{}{
		"candidates":  summary.Candidates,
		"unique":      summary.Unique,
		"cacheHits":   summary.CacheHits,
		"prepared":    summary.Prepared,
		"changed":     summary.Changed,
		"failed":      summary.Failed,
		"resolved":    summary.Resolved,
		"compressed":  summary.Compressed,
		"inputBytes":  summary.InputBytes,
		"outputBytes": summary.OutputBytes,
		"totalMs":     summary.TotalMs,
		"avgMs":       summary.AvgMs,
		"slowestMs":   summary.SlowestMs,
	}
}

// preprocessBatchImageData resolves and optionally compresses imageData fields ahead of execution.
// It deduplicates repeated payloads and processes unique items in parallel for throughput.
func preprocessBatchImageData(ops []batchOperation, compress bool, label string, showProgress bool) ([]batchOperation, imagePrepSummary) {
	summary := imagePrepSummary{}
	if len(ops) == 0 {
		return ops, summary
	}

	unique := make(map[string]struct{})
	for _, op := range ops {
		raw, ok := op.Params["imageData"].(string)
		if !ok {
			continue
		}
		raw = strings.TrimSpace(raw)
		if raw == "" || hasInterpolationToken(raw) {
			continue
		}
		summary.Candidates++
		unique[raw] = struct{}{}
	}
	summary.Unique = len(unique)
	summary.CacheHits = summary.Candidates - summary.Unique
	if summary.Unique == 0 {
		return ops, summary
	}

	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}
	if workers > 8 {
		workers = 8
	}
	if workers > summary.Unique {
		workers = summary.Unique
	}

	prefix := "[image-prep]"
	if strings.TrimSpace(label) != "" {
		prefix = fmt.Sprintf("[image-prep %s]", label)
	}
	if showProgress {
		fmt.Fprintf(os.Stderr, "%s preparing %d unique image payload(s) from %d operation(s) using %d worker(s)\n",
			prefix, summary.Unique, summary.Candidates, workers)
	}

	start := time.Now()
	results := make(map[string]imagePrepTaskResult, summary.Unique)
	tasks := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	completed := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for raw := range tasks {
				taskStart := time.Now()
				result := imagePrepTaskResult{
					Value:       raw,
					InputBytes:  len(raw),
					OutputBytes: len(raw),
				}

				resolved, err := imgutil.ResolveImageData(raw)
				if err != nil {
					result.Err = err
				} else {
					result.Value = resolved
					result.OutputBytes = len(resolved)
					if resolved != raw {
						result.Changed = true
						result.Resolved = true
					}

					if compress {
						compressed, compErr := imgutil.CompressBase64(resolved, imgutil.DefaultOptions())
						if compErr != nil {
							result.Err = compErr
						} else {
							result.Value = compressed
							result.OutputBytes = len(compressed)
							if compressed != resolved {
								result.Changed = true
								result.Compressed = true
							}
						}
					}
				}

				result.ElapsedMs = int(time.Since(taskStart).Milliseconds())

				mu.Lock()
				results[raw] = result
				completed++
				done := completed
				mu.Unlock()

				if showProgress {
					status := "ok"
					if result.Err != nil {
						status = "error"
					}
					fmt.Fprintf(os.Stderr, "%s %d/%d %s %dms %s\n",
						prefix, done, summary.Unique, status, result.ElapsedMs, shortImageSource(raw))
				}
			}
		}()
	}

	for raw := range unique {
		tasks <- raw
	}
	close(tasks)
	wg.Wait()

	elapsedSum := 0
	for _, result := range results {
		elapsedSum += result.ElapsedMs
		summary.InputBytes += result.InputBytes
		summary.OutputBytes += result.OutputBytes
		if result.ElapsedMs > summary.SlowestMs {
			summary.SlowestMs = result.ElapsedMs
		}
		if result.Err != nil {
			summary.Failed++
			continue
		}
		summary.Prepared++
		if result.Changed {
			summary.Changed++
		}
		if result.Resolved {
			summary.Resolved++
		}
		if result.Compressed {
			summary.Compressed++
		}
	}
	if summary.Unique > 0 {
		summary.AvgMs = elapsedSum / summary.Unique
	}
	summary.TotalMs = int(time.Since(start).Milliseconds())

	out := make([]batchOperation, len(ops))
	copy(out, ops)
	warned := make(map[string]bool)
	for i, op := range out {
		raw, ok := op.Params["imageData"].(string)
		if !ok {
			continue
		}
		raw = strings.TrimSpace(raw)
		if raw == "" || hasInterpolationToken(raw) {
			continue
		}
		result, ok := results[raw]
		if !ok {
			continue
		}
		if result.Err != nil {
			if showProgress && !warned[raw] {
				fmt.Fprintf(os.Stderr, "%s warning: %v (%s)\n", prefix, result.Err, shortImageSource(raw))
				warned[raw] = true
			}
			continue
		}
		if result.Value == op.Params["imageData"] {
			continue
		}
		cloned := make(map[string]interface{}, len(op.Params))
		for k, v := range op.Params {
			cloned[k] = v
		}
		cloned["imageData"] = result.Value
		out[i].Params = cloned
	}

	if showProgress {
		fmt.Fprintf(os.Stderr, "%s done in %dms (unique=%d, cacheHits=%d, changed=%d, failed=%d)\n",
			prefix, summary.TotalMs, summary.Unique, summary.CacheHits, summary.Changed, summary.Failed)
	}

	return out, summary
}

func collectCreatedRootFrameIDs(ops []batchOperation, results []map[string]interface{}) []string {
	seen := map[string]bool{}
	rootNodeIDs := make([]string, 0)
	for _, r := range results {
		ok, _ := r["ok"].(bool)
		if !ok {
			continue
		}
		idx := resultIndex(r["index"])
		if idx < 0 || idx >= len(ops) {
			continue
		}
		if !isRootFrameCreateOp(ops[idx]) {
			continue
		}
		resultMap, _ := r["result"].(map[string]interface{})
		if resultMap == nil {
			continue
		}
		nodeID, _ := resultMap["id"].(string)
		nodeID = strings.TrimSpace(nodeID)
		if nodeID == "" || seen[nodeID] {
			continue
		}
		seen[nodeID] = true
		rootNodeIDs = append(rootNodeIDs, nodeID)
	}
	return rootNodeIDs
}

func resultIndex(raw interface{}) int {
	switch v := raw.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return -1
	}
}

func isRootFrameCreateOp(op batchOperation) bool {
	cmd := strings.ToLower(strings.TrimSpace(op.Command))
	switch cmd {
	case "frame", "node.create_frame", "create_frame":
	default:
		return false
	}
	parentID, _ := op.Params["parentId"].(string)
	pid, _ := op.Params["pid"].(string)
	return strings.TrimSpace(parentID) == "" && strings.TrimSpace(pid) == ""
}

type lintSummary struct {
	Issues   int
	Warnings int
	Errors   int
	ByType   map[string]int
	Samples  []lintIssueSample
}

type lintIssueSample struct {
	Type     string
	NodeName string
	Message  string
}

func runPostBatchLint(client *ws.Client, rootNodeIDs []string, label string) lintSummary {
	summary := lintSummary{
		ByType: make(map[string]int),
	}
	if len(rootNodeIDs) == 0 {
		return summary
	}

	prefix := "[lint]"
	if strings.TrimSpace(label) != "" {
		prefix = fmt.Sprintf("[lint %s]", label)
	}

	fmt.Fprintf(os.Stderr, "%s checking %d root frame(s)\n", prefix, len(rootNodeIDs))

	const maxShown = 10
	total := 0
	shown := 0

	for _, nodeID := range rootNodeIDs {
		lintResult, lintErr := client.SendCommand("document.lint", map[string]interface{}{"nodeId": nodeID})
		if lintErr != nil {
			fmt.Fprintf(os.Stderr, "%s lint failed for %s: %s\n", prefix, nodeID, lintErr.Error())
			continue
		}

		var lintData map[string]interface{}
		if err := json.Unmarshal(lintResult, &lintData); err != nil {
			fmt.Fprintf(os.Stderr, "%s invalid lint response for %s\n", prefix, nodeID)
			continue
		}
		warnings, _ := lintData["warnings"].([]interface{})
		for _, w := range warnings {
			wm, _ := w.(map[string]interface{})
			if wm == nil {
				continue
			}
			sev, _ := wm["severity"].(string)
			sev = strings.ToLower(strings.TrimSpace(sev))
			if sev != "warning" && sev != "error" {
				continue
			}
			if sev == "error" {
				summary.Errors++
			} else {
				summary.Warnings++
			}
			summary.Issues++
			total++

			wtype, _ := wm["type"].(string)
			wtype = strings.TrimSpace(wtype)
			if wtype == "" {
				wtype = "lint"
			}
			summary.ByType[wtype]++
			msg, _ := wm["message"].(string)
			msg = strings.TrimSpace(msg)
			nodeName, _ := wm["nodeName"].(string)
			nodeName = strings.TrimSpace(nodeName)
			if len(summary.Samples) < 12 {
				summary.Samples = append(summary.Samples, lintIssueSample{
					Type:     wtype,
					NodeName: nodeName,
					Message:  msg,
				})
			}
			if shown >= maxShown {
				continue
			}
			lintNodeID, _ := wm["nodeId"].(string)

			icon := "WARN"
			if sev == "error" {
				icon = "ERROR"
			}

			location := strings.TrimSpace(lintNodeID)
			if strings.TrimSpace(nodeName) != "" && location != "" {
				location = nodeName + " (" + location + ")"
			}
			if location != "" {
				fmt.Fprintf(os.Stderr, "%s %s [%s] %s - %s\n", prefix, icon, wtype, location, msg)
			} else {
				fmt.Fprintf(os.Stderr, "%s %s [%s] %s\n", prefix, icon, wtype, msg)
			}

			if wtype == "default_name" && strings.TrimSpace(lintNodeID) != "" {
				fmt.Fprintf(os.Stderr, "%s fix: ahd-figma command node.modify '{\"nodeId\":\"%s\",\"name\":\"<semantic-name>\"}'\n", prefix, lintNodeID)
			}
			shown++
		}
	}

	if total == 0 {
		fmt.Fprintf(os.Stderr, "%s no warning/error issues found\n", prefix)
		return summary
	}
	if total > shown {
		fmt.Fprintf(os.Stderr, "%s %d warning/error issue(s) total; showing %d\n", prefix, total, shown)
		for _, hint := range lintGuidance(summary) {
			fmt.Fprintf(os.Stderr, "%s hint: %s\n", prefix, hint)
		}
		return summary
	}
	fmt.Fprintf(os.Stderr, "%s %d warning/error issue(s)\n", prefix, total)
	for _, hint := range lintGuidance(summary) {
		fmt.Fprintf(os.Stderr, "%s hint: %s\n", prefix, hint)
	}
	return summary
}

func lintGuidance(summary lintSummary) []string {
	out := make([]string, 0, 6)
	hasOverlap := summary.ByType["overlap"] > 0
	hasOverflow := summary.ByType["overflow"] > 0
	hasDefaultNames := summary.ByType["default_name"] > 0
	hasTextLarge := summary.ByType["text_too_large"] > 0
	hasTextSmall := summary.ByType["text_too_small"] > 0
	hasOversizedChild := summary.ByType["oversized_child"] > 0
	hasAbsoluteBad := summary.ByType["absolute_child_non_autolayout"] > 0
	hasAbsoluteOverflow := summary.ByType["absolute_overflow"] > 0

	if hasOverlap {
		out = append(out, "Overlap issues: keep content in auto-layout flow and reserve absolute positioning for decorative elements only.")
		for _, s := range summary.Samples {
			if strings.Contains(strings.ToLower(s.NodeName), "stat value") || strings.Contains(strings.ToLower(s.NodeName), "stat label") {
				out = append(out, "Stats overlap: enforce single-line stat values and increase value-to-label gap before placing labels.")
				break
			}
		}
		for _, s := range summary.Samples {
			msgLower := strings.ToLower(s.Message)
			if strings.Contains(msgLower, "parent \"banner") || (strings.Contains(strings.ToLower(s.NodeName), "headline") && strings.Contains(msgLower, "subtitle")) {
				out = append(out, "Banner overlap: use adaptive headline sizing and a smaller subtitle tier so subtitle placement is based on wrapped headline height.")
				break
			}
		}
	}
	if hasOverflow || hasOversizedChild {
		out = append(out, "Overflow/oversized issues: increase parent size or reduce child dimensions; prefer FILL/HUG sizing in auto-layout frames.")
	}
	if hasAbsoluteBad {
		out = append(out, "Invalid ABSOLUTE usage: layoutPositioning:ABSOLUTE is only valid on children of auto-layout parents. Remove ABSOLUTE for manual x/y parents.")
	}
	if hasAbsoluteOverflow {
		out = append(out, "Absolute child overflow: keep absolute overlays inside parent bounds or move them into normal auto-layout flow.")
	}
	if hasTextLarge {
		out = append(out, "Text too large: lower headline tier or enable auto-fit/adaptive tiering for long copy.")
	}
	if hasTextSmall {
		out = append(out, "Text too small: raise minimum font sizes and avoid caption tiers for primary body content.")
	}
	if hasDefaultNames {
		out = append(out, "Default naming issues: assign semantic names to every frame/layer so interpolation and debugging stay stable.")
	}
	if len(out) == 0 && summary.Issues > 0 {
		out = append(out, "Lint issues detected. Re-run with --strict-quality, inspect lint.byType, and iterate the batch JSON before final export.")
	}
	return uniqueStrings(out)
}

func lintSamplesForOutput(samples []lintIssueSample, limit int) []map[string]interface{} {
	if limit <= 0 {
		limit = 3
	}
	out := make([]map[string]interface{}, 0, limit)
	for _, sample := range samples {
		if len(out) >= limit {
			break
		}
		entry := map[string]interface{}{
			"type": sample.Type,
		}
		if strings.TrimSpace(sample.NodeName) != "" {
			entry["nodeName"] = sample.NodeName
		}
		if strings.TrimSpace(sample.Message) != "" {
			entry["message"] = sample.Message
		}
		out = append(out, entry)
	}
	return out
}

func uniqueStrings(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, item := range in {
		key := strings.TrimSpace(item)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, key)
	}
	return out
}

// stdoutIsTTY returns true when stdout is connected to a terminal (not piped).
func stdoutIsTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// marshalJSON returns v as JSON bytes using TTY-aware formatting.
// When stdout is a TTY (interactive terminal), it uses indented formatting
// for readability. When piped (non-TTY), it uses compact single-line JSON
// to save tokens for LLM agents.
func marshalJSON(v interface{}) ([]byte, error) {
	if stdoutIsTTY() {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// printJSON outputs v as JSON to stdout with TTY-aware formatting.
func printJSON(v interface{}) error {
	data, err := marshalJSON(v)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// printJSONErr outputs v as JSON to stderr with TTY-aware formatting.
// Used for structured error envelopes so agents can separate data (stdout)
// from errors (stderr).
func printJSONErr(v interface{}) {
	data, err := marshalJSON(v)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, string(data))
}

func roundTo(v float64, decimals int) float64 {
	if decimals < 0 {
		return v
	}
	pow := math.Pow(10, float64(decimals))
	return math.Round(v*pow) / pow
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
		X             float64 `json:"x"`
		Y             float64 `json:"y"`
		ExistingCount int     `json:"existingCount"`
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

	// Collect all root frame indices and their widths
	type rootFrame struct {
		idx   int
		width float64
	}
	var rootFrames []rootFrame
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
		w := getFloatParam(ops[i].Params, "width", getFloatParam(ops[i].Params, "w", 0))
		rootFrames = append(rootFrames, rootFrame{idx: i, width: w})
	}

	if len(rootFrames) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "⚠ Auto-placing: shifted %d root frame(s) to avoid overlap with %d existing frame(s).\n", len(rootFrames), freeSpace.ExistingCount)
	fmt.Fprintf(os.Stderr, "  Use --allow-overlap to place at exact coordinates.\n")

	// Place root frames side by side starting at freeSpace.X, freeSpace.Y
	// Each subsequent frame is offset by the previous frame's width + gap
	gap := float64(100)
	currentX := freeSpace.X
	for _, rf := range rootFrames {
		ops[rf.idx].Params["x"] = currentX
		ops[rf.idx].Params["y"] = freeSpace.Y
		currentX += rf.width + gap
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

	// Auto-normalize LLM output (strip fences, fix type→command, hoist top-level params)
	// unless --no-fix is set. Fixes are printed to stderr so callers can see corrections.
	if !batchNoFix {
		fixed, fixes, fixErr := batchutil.FixBatchOps(raw)
		if fixErr != nil {
			return nil, fmt.Errorf("cannot parse operations JSON: %w", fixErr)
		}
		if len(fixes) > 0 {
			fmt.Fprintf(os.Stderr, "✎ Auto-normalized %d issue(s) in batch input:\n", len(fixes))
			for _, f := range fixes {
				fmt.Fprintf(os.Stderr, "  • %s\n", f)
			}
		}
		raw = fixed
	}

	// Expand composite commands (slide, banner) into primitive ops
	if !batchNoFix {
		var rawOps []map[string]interface{}
		if unmErr := json.Unmarshal(raw, &rawOps); unmErr == nil {
			before := len(rawOps)
			expanded, expErr := batchutil.ExpandAllComposites(rawOps)
			if expErr != nil {
				return nil, fmt.Errorf("composite expansion failed: %w", expErr)
			}
			after := len(expanded)
			if after != before {
				fmt.Fprintf(os.Stderr, "Expanded %d composite commands → %d operations\n", before, after)
				remarshaled, mErr := json.Marshal(expanded)
				if mErr != nil {
					return nil, fmt.Errorf("failed to re-marshal expanded ops: %w", mErr)
				}
				raw = remarshaled
			}
		}
	}

	// Normalize CSS properties to Figma properties in all ops
	if !batchNoFix {
		var cssOps []map[string]interface{}
		if unmErr := json.Unmarshal(raw, &cssOps); unmErr == nil {
			for _, op := range cssOps {
				if params, ok := op["params"].(map[string]interface{}); ok {
					batchutil.NormalizeCSSProps(params)
				}
			}
			remarshaled, mErr := json.Marshal(cssOps)
			if mErr == nil {
				raw = remarshaled
			}
		}
	}

	// Resolve semantic token aliases (e.g. sz:"hero", padding:"side", w:"content")
	// against the detected root frame size. This is additive and backward compatible.
	{
		var tokenOps []map[string]interface{}
		if unmErr := json.Unmarshal(raw, &tokenOps); unmErr == nil {
			applied, rootWidth := batchutil.ResolveTokenAliases(tokenOps)
			if applied > 0 {
				fmt.Fprintf(os.Stderr, "Resolved %d token alias(es) using root width %dpx\n", applied, rootWidth)
				if remarshaled, mErr := json.Marshal(tokenOps); mErr == nil {
					raw = remarshaled
				}
			}
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

// checkPluginConnected checks if a Figma plugin is connected on the target channel.
// If not, it polls for up to 30s with user feedback before failing.
func checkPluginConnected(channelKey string) error {
	cfg := loadConfig()
	statusURL := fmt.Sprintf("http://%s:%d/status", cfg.ServerHost, cfg.Port)
	httpClient := &http.Client{Timeout: 2 * time.Second}

	check := func() (bool, error) {
		resp, err := httpClient.Get(statusURL)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()
		var status relayStatus
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return false, err
		}
		if count, ok := status.Channels[channelKey]; ok && count > 0 {
			return true, nil
		}
		// Also check if ANY channel has clients (channel might not be resolved yet)
		for _, count := range status.Channels {
			if count > 0 {
				return true, nil
			}
		}
		return false, nil
	}

	connected, err := check()
	if err != nil || connected {
		return nil // Can't check or already connected — proceed normally
	}

	// No plugin connected — give feedback and poll
	fmt.Fprintf(os.Stderr, "Relay running on localhost:%d but no Figma plugin connected on channel %q.\n", cfg.Port, channelKey)
	fmt.Fprintf(os.Stderr, "Plugin auto-reconnects every ~5s. You can also:\n")
	fmt.Fprintf(os.Stderr, "  - Press \"Connect\" in the Figma plugin UI\n")
	fmt.Fprintf(os.Stderr, "  - Pass --channel <key> if using a different channel\n")
	fmt.Fprintf(os.Stderr, "Waiting up to 30s for plugin...\n")

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(2 * time.Second)
		connected, _ = check()
		if connected {
			fmt.Fprintf(os.Stderr, "Plugin connected!\n")
			return nil
		}
	}

	return fmt.Errorf("no Figma plugin connected after 30s. Start the Figma plugin and click Connect, then retry")
}

// collectBatchInputs processes positional args for batch command.
// Returns a list of file paths and/or inline JSON string.
func collectBatchInputs(args []string, flagOps, flagFile string) (files []string, inlineJSON string, err error) {
	if len(args) == 0 {
		// Legacy: use flags or stdin
		if flagOps != "" {
			return nil, flagOps, nil
		}
		if flagFile != "" {
			return []string{flagFile}, "", nil
		}
		// Will fall through to stdin in loadBatchOperations
		return nil, "", nil
	}

	if len(args) == 1 {
		arg := strings.TrimSpace(args[0])
		if strings.HasPrefix(arg, "[") {
			if flagOps != "" {
				return nil, "", fmt.Errorf("provide operations as positional arg OR --operations, not both")
			}
			return nil, arg, nil
		}
		// Could be a file, directory, or glob
		expanded, err := expandBatchPath(arg)
		if err != nil {
			return nil, "", err
		}
		return expanded, "", nil
	}

	// Multiple args — all must be files/directories/globs
	for _, arg := range args {
		expanded, err := expandBatchPath(arg)
		if err != nil {
			return nil, "", err
		}
		files = append(files, expanded...)
	}
	if len(files) == 0 {
		return nil, "", fmt.Errorf("no .json files found in provided paths")
	}
	return files, "", nil
}

// expandBatchPath expands a single path arg into file paths.
// Handles directories (all .json inside), globs, and plain files.
func expandBatchPath(path string) ([]string, error) {
	// Check if it's a glob pattern
	if strings.ContainsAny(path, "*?[") {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern %q: %w", path, err)
		}
		// Filter to .json only
		var jsonFiles []string
		for _, m := range matches {
			if strings.HasSuffix(strings.ToLower(m), ".json") {
				jsonFiles = append(jsonFiles, m)
			}
		}
		if len(jsonFiles) == 0 {
			return nil, fmt.Errorf("no .json files matched pattern %q", path)
		}
		sort.Strings(jsonFiles)
		return jsonFiles, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access %q: %w", path, err)
	}

	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return nil, fmt.Errorf("cannot read directory %q: %w", path, err)
		}
		var jsonFiles []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".json") {
				jsonFiles = append(jsonFiles, filepath.Join(path, e.Name()))
			}
		}
		if len(jsonFiles) == 0 {
			return nil, fmt.Errorf("no .json files found in directory %q", path)
		}
		sort.Strings(jsonFiles)
		return jsonFiles, nil
	}

	return []string{path}, nil
}

// runMultiBatch runs one or more batch files, optionally in parallel.
func runMultiBatch(files []string, channelKey string) error {
	if batchParallel && len(files) > 1 {
		return runMultiBatchParallel(files, channelKey)
	}
	return runMultiBatchSequential(files, channelKey)
}

type batchFileResult struct {
	File    string      `json:"file"`
	OK      bool        `json:"ok"`
	Steps   interface{} `json:"steps,omitempty"`
	Summary interface{} `json:"summary,omitempty"`
	Timing  interface{} `json:"timing,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func runSingleBatchFile(filePath, channelKey string) batchFileResult {
	ops, err := loadBatchOperations("", filePath)
	if err != nil {
		return batchFileResult{File: filePath, OK: false, Error: err.Error()}
	}
	if len(ops) == 0 {
		return batchFileResult{File: filePath, OK: false, Error: "operations array is empty"}
	}

	client, err := newConnectedClient(channelKey, batchLive)
	if err != nil {
		return batchFileResult{File: filePath, OK: false, Error: err.Error()}
	}
	defer client.Close()

	// Auto-place root frames
	if !batchAllowOverlap {
		autoPlaceRootFrames(client, ops)
	}
	ops, imagePrep := preprocessBatchImageData(ops, batchCompressImages, filepath.Base(filePath), true)

	results := make([]map[string]interface{}, 0, len(ops))
	stepStates := make([]batchutil.StepState, 0, len(ops))
	retryDelay := time.Duration(batchRetryDelayMs) * time.Millisecond
	succeeded, failed, retriesUsed := 0, 0, 0
	stoppedEarly := false
	batchStart := time.Now()

	for i, op := range ops {
		opStart := time.Now()
		op.Name = batchutil.SanitizeStepName(op.Name)
		params := batchutil.NormalizeBatchParams(op.Command, op.Params)

		if batchInterpolation {
			interpolatedParams, iErr := batchutil.InterpolateParams(params, stepStates)
			if iErr != nil {
				failed++
				entry := map[string]interface{}{
					"index": i, "name": op.Name, "command": op.Command,
					"ok": false, "error": fmt.Sprintf("interpolation error: %v", iErr),
					"attempts": 0, "elapsedMs": int(time.Since(opStart).Milliseconds()),
				}
				results = append(results, entry)
				stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: false, Error: entry["error"].(string)})
				if batchFailFast {
					stoppedEarly = true
					break
				}
				continue
			}
			params = interpolatedParams
		}

		// Runtime fallback for interpolated imageData values that could not be preprocessed.
		if rawImage, ok := params["imageData"].(string); ok && needsRuntimeImagePrep(rawImage) {
			if resolved, changed, resolveErr := imgutil.ResolveParamsImageData(params); resolveErr != nil {
				fmt.Fprintf(os.Stderr, "[resolve] %s step %d warning: %v\n", filepath.Base(filePath), i, resolveErr)
			} else if changed {
				params = resolved
			}
			if batchCompressImages {
				if compressed, changed, compErr := imgutil.CompressParamsImageData(params, imgutil.DefaultOptions()); compErr != nil {
					fmt.Fprintf(os.Stderr, "[compress] %s step %d warning: %v\n", filepath.Base(filePath), i, compErr)
				} else if changed {
					params = compressed
				}
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
				"index": i, "name": op.Name, "command": op.Command,
				"ok": false, "attempts": attempts, "error": sendErr.Error(),
				"elapsedMs": int(time.Since(opStart).Milliseconds()),
			}
			results = append(results, entry)
			stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: false, Error: sendErr.Error()})
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
			"index": i, "name": op.Name, "command": op.Command,
			"ok": true, "attempts": attempts, "result": parsed,
			"elapsedMs": int(time.Since(opStart).Milliseconds()),
		}
		results = append(results, entry)
		stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: true, Result: parsed})
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

	var lintInfo lintSummary
	if batchLint && !batchNoLint && failed == 0 {
		rootNodeIDs := collectCreatedRootFrameIDs(ops, results)
		lintInfo = runPostBatchLint(client, rootNodeIDs, filepath.Base(filePath))
		if batchStrictQuality && lintInfo.Issues > 0 {
			failed++
		}
	}

	_ = stoppedEarly
	summary := map[string]interface{}{
		"total":         len(ops),
		"processed":     processed,
		"succeeded":     succeeded,
		"failed":        failed,
		"pending":       pending,
		"retriesUsed":   retriesUsed,
		"failFast":      batchFailFast,
		"interpolation": batchInterpolation,
		"imagePrep":     imagePrepSummaryMap(imagePrep),
	}
	if batchLint && !batchNoLint {
		lintMap := map[string]interface{}{
			"issues":   lintInfo.Issues,
			"warnings": lintInfo.Warnings,
			"errors":   lintInfo.Errors,
			"byType":   lintInfo.ByType,
		}
		lintGuides := lintGuidance(lintInfo)
		if len(lintGuides) > 0 {
			lintMap["guidance"] = lintGuides
		}
		lintSamples := lintSamplesForOutput(lintInfo.Samples, 5)
		if len(lintSamples) > 0 {
			lintMap["samples"] = lintSamples
		}
		summary["lint"] = lintMap
		if batchStrictQuality {
			if lintInfo.Issues > 0 {
				summary["qualityGate"] = "failed"
				summary["qualityIssues"] = lintInfo.Issues
				if len(lintGuides) > 0 {
					summary["qualityGuidance"] = lintGuides
				}
			} else {
				summary["qualityGate"] = "passed"
			}
		}
	}
	resultOK := failed == 0 && pending == 0
	resultErr := ""
	if batchStrictQuality && lintInfo.Issues > 0 {
		resultOK = false
		resultErr = fmt.Sprintf("strict quality gate failed: %d lint warning/error issue(s)", lintInfo.Issues)
	}
	return batchFileResult{
		File:    filePath,
		OK:      resultOK,
		Steps:   results,
		Summary: summary,
		Timing: map[string]interface{}{
			"totalMs": totalMs, "avgMs": int(avgMs), "opsPerSec": roundTo(opsPerSec, 2), "imagePrepMs": imagePrep.TotalMs,
		},
		Error: resultErr,
	}
}

func runMultiBatchSequential(files []string, channelKey string) error {
	fileResults := make([]batchFileResult, 0, len(files))
	allOK := true
	for _, f := range files {
		fmt.Fprintf(os.Stderr, "[batch] running %s...\n", filepath.Base(f))
		result := runSingleBatchFile(f, channelKey)
		if !result.OK {
			allOK = false
		}
		fileResults = append(fileResults, result)
	}
	out := map[string]interface{}{
		"ok":    allOK,
		"files": fileResults,
	}
	if err := printJSON(out); err != nil {
		return err
	}
	if batchStrictQuality && !allOK {
		return fmt.Errorf("strict quality gate failed for one or more batch files")
	}
	return nil
}

func runMultiBatchParallel(files []string, channelKey string) error {
	maxConcurrency := 4
	if len(files) < maxConcurrency {
		maxConcurrency = len(files)
	}
	sem := make(chan struct{}, maxConcurrency)
	fileResults := make([]batchFileResult, len(files))
	var wg sync.WaitGroup

	for i, f := range files {
		wg.Add(1)
		go func(idx int, filePath string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			fmt.Fprintf(os.Stderr, "[batch] starting %s (parallel)...\n", filepath.Base(filePath))
			fileResults[idx] = runSingleBatchFile(filePath, channelKey)
			fmt.Fprintf(os.Stderr, "[batch] finished %s (ok=%v)\n", filepath.Base(filePath), fileResults[idx].OK)
		}(i, f)
	}
	wg.Wait()

	allOK := true
	for _, r := range fileResults {
		if !r.OK {
			allOK = false
		}
	}
	out := map[string]interface{}{
		"ok":       allOK,
		"parallel": true,
		"files":    fileResults,
	}
	if err := printJSON(out); err != nil {
		return err
	}
	if batchStrictQuality && !allOK {
		return fmt.Errorf("strict quality gate failed for one or more batch files")
	}
	return nil
}

// --- Benchmark commands ---

// ── Extract command ──────────────────────────────────────────────────

var extractCmd = &cobra.Command{
	Use:   "extract [file.html]",
	Short: "Convert HTML/CSS to batch JSON",
	Long:  "Parse an HTML file and output Figma batch operations as composite commands (slide/banner).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		canvasWidth, _ := cmd.Flags().GetInt("width")
		canvasHeight, _ := cmd.Flags().GetInt("height")
		bannerWidth, _ := cmd.Flags().GetInt("banner-width")
		bannerHeight, _ := cmd.Flags().GetInt("banner-height")
		outputPath, _ := cmd.Flags().GetString("output")

		f, ferr := os.Open(inputPath)
		if ferr != nil {
			return ferr
		}
		defer f.Close()

		ops, err := extract.FromHTML(f, extract.Options{
			CanvasWidth:  canvasWidth,
			CanvasHeight: canvasHeight,
			BannerWidth:  bannerWidth,
			BannerHeight: bannerHeight,
			BaseDir:      filepath.Dir(inputPath),
		})
		if err != nil {
			return err
		}

		if outputPath != "" {
			data, _ := json.MarshalIndent(ops, "", "  ")
			return os.WriteFile(outputPath, data, 0644)
		}
		return printJSON(ops)
	},
}

// ── Benchmark command ────────────────────────────────────────────────

var benchmarkRuns int
var benchmarkAllowOverlap bool
var benchmarkPhaseAMs int
var benchmarkWidth int
var benchmarkHeight int

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Provider-agnostic performance benchmarking",
	Long: `Benchmark AHD Figma workflows without calling any LLM API directly.
Times stdin/stdout pipes so you can wrap ANY LLM curl call externally.

Subcommands:
  exec     Time batch execution against Figma
  pipe     Accept batch JSON from stdin with external LLM timing
  compare  Compare extraction methods (stub)`,
}

var benchmarkExecCmd = &cobra.Command{
	Use:   "exec [file.json]",
	Short: "Time batch execution only",
	Long: `Runs a batch JSON file N times and reports aggregate timing.
Connects to Figma for each run, executes the batch, and measures Phase B (CLI execution).

Example:
  benchmark exec ops.json --runs 5
  benchmark exec ops.json --runs 3 --allow-overlap`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)
		filePath := args[0]

		if benchmarkRuns < 1 {
			return fmt.Errorf("--runs must be >= 1")
		}

		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		channelKey, err := resolveChannel(batchChannel)
		if err != nil {
			return err
		}
		if err := checkPluginConnected(channelKey); err != nil {
			return err
		}

		runs := make([]benchmark.RunResult, 0, benchmarkRuns)

		for i := 0; i < benchmarkRuns; i++ {
			fmt.Fprintf(os.Stderr, "[benchmark] run %d/%d...\n", i+1, benchmarkRuns)

			result := benchExecOnce(filePath, channelKey)
			runs = append(runs, result)

			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "[benchmark] run %d error: %v\n", i+1, result.Error)
			} else {
				fmt.Fprintf(os.Stderr, "[benchmark] run %d: %d ops in %s (%.1f ops/s)\n",
					i+1, result.PhaseB.OpsCount, result.PhaseB.Duration.Round(time.Millisecond),
					float64(result.PhaseB.OpsCount)/result.PhaseB.Duration.Seconds())
			}
		}

		agg := benchmark.Aggregate(runs)
		fmt.Println(agg.String())
		return nil
	},
}

// benchExecOnce runs a single batch file and returns a RunResult with Phase B timing.
func benchExecOnce(filePath, channelKey string) benchmark.RunResult {
	ops, err := loadBatchOperations("", filePath)
	if err != nil {
		return benchmark.RunResult{Error: err}
	}
	if len(ops) == 0 {
		return benchmark.RunResult{Error: fmt.Errorf("operations array is empty")}
	}

	client, err := newConnectedClient(channelKey, false)
	if err != nil {
		return benchmark.RunResult{Error: err}
	}
	defer client.Close()

	if !benchmarkAllowOverlap {
		autoPlaceRootFrames(client, ops)
	}
	ops, _ = preprocessBatchImageData(ops, batchCompressImages, "benchmark", false)

	retryDelay := time.Duration(batchRetryDelayMs) * time.Millisecond
	maxAttempts := batchRetries + 1
	opsCount := 0
	errors := 0

	start := time.Now()
	stepStates := make([]batchutil.StepState, 0, len(ops))

	for i, op := range ops {
		op.Name = batchutil.SanitizeStepName(op.Name)
		params := batchutil.NormalizeBatchParams(op.Command, op.Params)

		if batchInterpolation {
			interpolatedParams, iErr := batchutil.InterpolateParams(params, stepStates)
			if iErr != nil {
				errors++
				stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: false, Error: iErr.Error()})
				continue
			}
			params = interpolatedParams
		}

		// Runtime fallback for interpolated imageData values that could not be preprocessed.
		if rawImage, ok := params["imageData"].(string); ok && needsRuntimeImagePrep(rawImage) {
			if resolved, changed, resolveErr := imgutil.ResolveParamsImageData(params); resolveErr != nil {
				fmt.Fprintf(os.Stderr, "[resolve] step %d warning: %v\n", i, resolveErr)
			} else if changed {
				params = resolved
			}
			if batchCompressImages {
				if compressed, changed, compErr := imgutil.CompressParamsImageData(params, imgutil.DefaultOptions()); compErr != nil {
					fmt.Fprintf(os.Stderr, "[compress] step %d warning: %v\n", i, compErr)
				} else if changed {
					params = compressed
				}
			}
		}

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

		opsCount++
		if sendErr != nil {
			errors++
			stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: false, Error: sendErr.Error()})
			continue
		}

		var parsed interface{}
		if err := json.Unmarshal(result, &parsed); err != nil {
			parsed = string(result)
		}
		stepStates = append(stepStates, batchutil.StepState{Index: i, Name: op.Name, Command: op.Command, OK: true, Result: parsed})
	}

	elapsed := time.Since(start)

	return benchmark.RunResult{
		PhaseB: benchmark.PhaseTiming{
			Label:    "CLI Exec",
			Duration: elapsed,
			OpsCount: opsCount,
			Errors:   errors,
		},
	}
}

var benchmarkPipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: "Accept batch JSON from stdin with external LLM timing",
	Long: `Reads batch JSON from stdin, executes it against Figma, and reports
combined Phase A (LLM generation, user-provided) + Phase B (CLI execution) timing.

User provides their LLM generation time via --phase-a-ms.

Example:
  cat ops.json | benchmark pipe --phase-a-ms 4200
  curl ... | jq .ops | benchmark pipe --phase-a-ms 3500`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetOutput(io.Discard)

		if err := ensureRelayIfNeeded(); err != nil {
			return err
		}

		channelKey, err := resolveChannel(batchChannel)
		if err != nil {
			return err
		}
		if err := checkPluginConnected(channelKey); err != nil {
			return err
		}

		// Read batch JSON from stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return fmt.Errorf("no input on stdin. Pipe batch JSON, e.g.: cat ops.json | benchmark pipe --phase-a-ms 4200")
		}

		raw, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		// Write to a temp file and use benchExecOnce
		tmpFile, err := os.CreateTemp("", "ahd-bench-*.json")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write(raw); err != nil {
			tmpFile.Close()
			return fmt.Errorf("failed to write temp file: %w", err)
		}
		tmpFile.Close()

		fmt.Fprintf(os.Stderr, "[benchmark] executing batch from stdin...\n")

		result := benchExecOnce(tmpFile.Name(), channelKey)
		result.PhaseA = benchmark.PhaseTiming{
			Label:    "LLM Gen",
			Duration: time.Duration(benchmarkPhaseAMs) * time.Millisecond,
		}

		runs := []benchmark.RunResult{result}
		agg := benchmark.Aggregate(runs)
		fmt.Println(agg.String())
		return nil
	},
}

var benchmarkCompareCmd = &cobra.Command{
	Use:   "compare [file.html]",
	Short: "Compare extraction methods (stub)",
	Long: `Compare lightweight extraction vs computed extraction on an HTML file.
Runs both extraction methods, executes the resulting batches, and prints
a side-by-side comparison table.

NOTE: This command requires the extract package (Task 5) which is being
built in parallel. Currently stubbed.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement once the extract package (Task 5) is available.
		// The plan:
		// 1. Run lightweight extraction on file.html -> batch JSON A
		// 2. Run computed extraction on file.html -> batch JSON B (if Chrome available)
		// 3. Execute batch A N times, time it
		// 4. Execute batch B N times, time it
		// 5. Print side-by-side comparison table
		return fmt.Errorf("benchmark compare is not yet implemented (waiting for extract package, Task 5)")
	},
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
	commandCmd.Flags().BoolVar(&commandCompressImages, "compress-images", false, "Compress resolved imageData before sending (requires ImageMagick)")

	batchCmd.Flags().StringVarP(&batchOperations, "operations", "o", "", "JSON array of operations")
	batchCmd.Flags().StringVarP(&batchOperationsFile, "operations-file", "f", "", "Path to JSON file containing operations array")
	batchCmd.Flags().BoolVar(&batchLive, "live", false, "Print live command progress while batch runs (image prep progress is always on stderr)")
	batchCmd.Flags().StringVar(&batchChannel, "channel", "", "Channel key override (optional)")
	batchCmd.Flags().BoolVar(&batchFailFast, "fail-fast", false, "Stop at first failed operation")
	batchCmd.Flags().IntVar(&batchRetries, "retries", 1, "Retry count per operation after first attempt")
	batchCmd.Flags().IntVar(&batchRetryDelayMs, "retry-delay-ms", 250, "Delay between retries in milliseconds")
	batchCmd.Flags().BoolVar(&batchInterpolation, "interpolate", true, "Enable placeholder interpolation from prior step results")
	batchCmd.Flags().BoolVar(&batchCompressImages, "compress-images", false, "Compress prepared imageData before sending (requires ImageMagick)")
	batchCmd.Flags().BoolVar(&batchAllowOverlap, "allow-overlap", false, "Skip auto-placement and place frames at exact coordinates (may overlap existing work)")
	batchCmd.Flags().BoolVar(&batchParallel, "parallel", false, "Run multiple batch files concurrently (max 4 parallel)")
	batchCmd.Flags().BoolVar(&batchNoFix, "no-fix", false, "Skip automatic LLM output normalization (use if your JSON is already valid)")
	batchCmd.Flags().BoolVar(&batchCompact, "compact", false, "Minimal output: ok + summary/timing/imagePrep + compact step results (saves tokens for LLM agents)")
	batchCmd.Flags().BoolVar(&batchLint, "lint", true, "Auto-check created frames for overlaps, overflow, naming, and text sizing issues")
	batchCmd.Flags().BoolVar(&batchNoLint, "no-lint", false, "Disable post-batch design lint checks")
	batchCmd.Flags().BoolVar(&batchStrictQuality, "strict-quality", false, "Fail the batch if lint reports any warning/error issue (quality gate)")

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

	// Benchmark subcommands
	benchmarkExecCmd.Flags().IntVar(&benchmarkRuns, "runs", 3, "Number of benchmark runs")
	benchmarkExecCmd.Flags().BoolVar(&benchmarkAllowOverlap, "allow-overlap", false, "Skip auto-placement (may overlap existing work)")
	benchmarkExecCmd.Flags().StringVar(&batchChannel, "channel", "", "Channel key override")

	benchmarkPipeCmd.Flags().IntVar(&benchmarkPhaseAMs, "phase-a-ms", 0, "LLM generation time in milliseconds (user-reported)")
	benchmarkPipeCmd.Flags().BoolVar(&benchmarkAllowOverlap, "allow-overlap", false, "Skip auto-placement (may overlap existing work)")
	benchmarkPipeCmd.Flags().StringVar(&batchChannel, "channel", "", "Channel key override")

	benchmarkCompareCmd.Flags().IntVar(&benchmarkRuns, "runs", 3, "Number of benchmark runs per method")
	benchmarkCompareCmd.Flags().IntVar(&benchmarkWidth, "width", 1080, "Canvas width for extraction")
	benchmarkCompareCmd.Flags().IntVar(&benchmarkHeight, "height", 1350, "Canvas height for extraction")

	extractCmd.Flags().Int("width", 1080, "Canvas width")
	extractCmd.Flags().Int("height", 1080, "Canvas height")
	extractCmd.Flags().Int("banner-width", 1200, "Banner canvas width")
	extractCmd.Flags().Int("banner-height", 400, "Banner canvas height")
	extractCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	rootCmd.AddCommand(extractCmd)

	benchmarkCmd.AddCommand(benchmarkExecCmd)
	benchmarkCmd.AddCommand(benchmarkPipeCmd)
	benchmarkCmd.AddCommand(benchmarkCompareCmd)
	rootCmd.AddCommand(benchmarkCmd)

	configInitCmd.Flags().BoolVar(&configInitForce, "force", false, "Overwrite existing config file")
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)

	if err := rootCmd.Execute(); err != nil {
		printJSONErr(map[string]interface{}{
			"error": err.Error(),
			"code":  classifyTopLevelError(err),
		})
		os.Exit(1)
	}
}

// classifyTopLevelError categorises a cobra/command error into a machine-readable code.
func classifyTopLevelError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "not connected"), strings.Contains(msg, "disconnected"):
		return "CONNECTION_ERROR"
	case strings.Contains(msg, "timed out"), strings.Contains(msg, "timeout"):
		return "TIMEOUT"
	case strings.Contains(msg, "unknown command"), strings.Contains(msg, "not found"):
		return "UNKNOWN_COMMAND"
	case strings.Contains(msg, "required"):
		return "VALIDATION_ERROR"
	default:
		return "ERROR"
	}
}
