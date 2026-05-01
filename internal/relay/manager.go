package relay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	stateDirName  = ".ai-happy-design"
	stateFileName = "relay.json"
	logFileName   = "relay.log"
	agentLabel    = "com.ai-happy-design.relay"
)

// State stores local relay process metadata for lifecycle management.
type State struct {
	Version    int      `json:"version"`
	PID        int      `json:"pid"`
	Port       int      `json:"port"`
	Host       string   `json:"host"`
	StartedAt  string   `json:"startedAt"`
	Executable string   `json:"executable"`
	Command    []string `json:"command"`
	LogPath    string   `json:"logPath"`
}

// EnsureOptions control relay auto-start behavior.
type EnsureOptions struct {
	Host         string
	Port         int
	AutoStart    bool
	Wait         time.Duration
	PollInterval time.Duration
}

// EnsureResult describes relay availability outcome.
type EnsureResult struct {
	Started        bool                   `json:"started"`
	AlreadyHealthy bool                   `json:"alreadyHealthy"`
	Status         map[string]interface{} `json:"status,omitempty"`
	State          *State                 `json:"state,omitempty"`
}

// StopResult summarizes stop outcome.
type StopResult struct {
	Stopped   bool   `json:"stopped"`
	PID       int    `json:"pid,omitempty"`
	Message   string `json:"message"`
	StatePath string `json:"statePath,omitempty"`
}

// InstallAgentResult summarizes launch agent installation.
type InstallAgentResult struct {
	Label     string `json:"label"`
	PlistPath string `json:"plistPath"`
	Loaded    bool   `json:"loaded"`
	LogPath   string `json:"logPath"`
}

// IsLocalHost reports whether host maps to local machine.
func IsLocalHost(host string) bool {
	h := strings.TrimSpace(strings.ToLower(host))
	return h == "" || h == "localhost" || h == "127.0.0.1" || h == "::1"
}

// StateFilePath returns the managed relay state file path.
func StateFilePath() (string, error) {
	dir, err := stateDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, stateFileName), nil
}

// LogFilePath returns the managed relay log file path.
func LogFilePath() (string, error) {
	dir, err := stateDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, logFileName), nil
}

func stateDirPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve home directory: %w", err)
	}
	return filepath.Join(home, stateDirName), nil
}

func ensureStateDir() (string, error) {
	dir, err := stateDirPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("unable to create state directory %s: %w", dir, err)
	}
	return dir, nil
}

// LoadState reads relay state metadata if present.
func LoadState() (*State, error) {
	path, err := StateFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to read relay state: %w", err)
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("invalid relay state json: %w", err)
	}
	return &state, nil
}

func saveState(state *State) error {
	if state == nil {
		return fmt.Errorf("cannot save nil relay state")
	}
	if _, err := ensureStateDir(); err != nil {
		return err
	}
	path, err := StateFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode relay state: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("failed to write relay state: %w", err)
	}
	return nil
}

func removeStateFile() error {
	path, err := StateFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove relay state: %w", err)
	}
	return nil
}

// ProbeStatus calls relay /status and returns parsed payload.
func ProbeStatus(host string, port int, timeout time.Duration) (map[string]interface{}, bool, error) {
	if timeout <= 0 {
		timeout = 1500 * time.Millisecond
	}
	url := fmt.Sprintf("http://%s:%d/status", host, port)
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("status endpoint returned %d", resp.StatusCode)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, fmt.Errorf("status payload decode failed: %w", err)
	}
	status, _ := payload["status"].(string)
	return payload, status == "ok", nil
}

// Ensure verifies relay availability and optionally auto-starts it.
func Ensure(opts EnsureOptions) (*EnsureResult, error) {
	host := opts.Host
	if host == "" {
		host = "localhost"
	}
	if opts.Port <= 0 {
		opts.Port = 3055
	}
	if opts.Wait <= 0 {
		opts.Wait = 6 * time.Second
	}
	if opts.PollInterval <= 0 {
		opts.PollInterval = 200 * time.Millisecond
	}

	status, healthy, err := ProbeStatus(host, opts.Port, 1200*time.Millisecond)
	if err == nil && healthy {
		return &EnsureResult{
			Started:        false,
			AlreadyHealthy: true,
			Status:         status,
		}, nil
	}

	if !opts.AutoStart {
		return nil, fmt.Errorf("relay is not running on %s:%d and auto-start is disabled (--no-auto-relay)", host, opts.Port)
	}
	if !IsLocalHost(host) {
		return nil, fmt.Errorf("relay is not running on %s:%d and auto-start only supports localhost/127.0.0.1", host, opts.Port)
	}

	inUse, useErr := IsPortInUse(opts.Port)
	if useErr != nil {
		return nil, useErr
	}
	if inUse {
		owner, ownerIsAHD := PortOwner(opts.Port)
		if owner != "" && !ownerIsAHD {
			return nil, fmt.Errorf("port %d is already in use by %s; refusing to start relay", opts.Port, owner)
		}

		if waitForHealthy(host, opts.Port, opts.Wait, opts.PollInterval) {
			status, healthy, _ := ProbeStatus(host, opts.Port, 1200*time.Millisecond)
			if healthy {
				return &EnsureResult{
					Started:        false,
					AlreadyHealthy: true,
					Status:         status,
				}, nil
			}
		}

		if owner != "" {
			return nil, fmt.Errorf("port %d is occupied by %s but relay /status is not healthy", opts.Port, owner)
		}
		return nil, fmt.Errorf("port %d is in use but relay /status is not healthy", opts.Port)
	}

	state, err := Start(host, opts.Port)
	if err != nil {
		return nil, err
	}

	if !waitForHealthy(host, opts.Port, opts.Wait, opts.PollInterval) {
		return nil, fmt.Errorf("started relay (pid %d) but /status did not become healthy in %s; check logs at %s", state.PID, opts.Wait, state.LogPath)
	}

	status, _, _ = ProbeStatus(host, opts.Port, 1200*time.Millisecond)
	return &EnsureResult{
		Started:        true,
		AlreadyHealthy: false,
		Status:         status,
		State:          state,
	}, nil
}

func waitForHealthy(host string, port int, wait, poll time.Duration) bool {
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		if _, ok, _ := ProbeStatus(host, port, 900*time.Millisecond); ok {
			return true
		}
		time.Sleep(poll)
	}
	return false
}

// Start spawns local relay process in background and persists metadata.
func Start(host string, port int) (*State, error) {
	if !IsLocalHost(host) {
		return nil, fmt.Errorf("relay start supports localhost hosts only, got %q", host)
	}
	if port <= 0 {
		return nil, fmt.Errorf("invalid relay port %d", port)
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("unable to locate current executable: %w", err)
	}

	logPath, err := LogFilePath()
	if err != nil {
		return nil, err
	}
	if _, err := ensureStateDir(); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("unable to open relay log file: %w", err)
	}

	cmd := exec.Command(exe, "ws")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		fmt.Sprintf("SERVER_HOST=%s", host),
	)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	detachProcess(cmd)

	if err := cmd.Start(); err != nil {
		_ = logFile.Close()
		return nil, fmt.Errorf("failed to start relay process: %w", err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()
	_ = logFile.Close()

	state := &State{
		Version:    1,
		PID:        pid,
		Port:       port,
		Host:       host,
		StartedAt:  time.Now().UTC().Format(time.RFC3339),
		Executable: exe,
		Command:    []string{exe, "ws"},
		LogPath:    logPath,
	}
	if err := saveState(state); err != nil {
		return nil, err
	}

	return state, nil
}

// Stop terminates relay process tracked in state file.
func Stop() (*StopResult, error) {
	statePath, err := StateFilePath()
	if err != nil {
		return nil, err
	}

	state, err := LoadState()
	if err != nil {
		return nil, err
	}
	if state == nil {
		return &StopResult{
			Stopped:   false,
			Message:   "no managed relay state found",
			StatePath: statePath,
		}, nil
	}

	if state.PID <= 0 {
		_ = removeStateFile()
		return &StopResult{
			Stopped:   false,
			Message:   "state file had no valid pid; state cleared",
			StatePath: statePath,
		}, nil
	}

	if cmdline, cmdErr := processCommand(state.PID); cmdErr == nil {
		lower := strings.ToLower(cmdline)
		if cmdline != "" && !looksLikeAHDProcess(lower) {
			return nil, fmt.Errorf("refusing to stop pid %d because it does not look like ai-happy-design/ahd-figma: %s", state.PID, cmdline)
		}
	}

	proc, err := os.FindProcess(state.PID)
	if err != nil {
		return nil, fmt.Errorf("unable to find relay process %d: %w", state.PID, err)
	}
	if err := proc.Kill(); err != nil {
		// If it is already gone, clear state and report.
		if strings.Contains(strings.ToLower(err.Error()), "process already finished") {
			_ = removeStateFile()
			return &StopResult{
				Stopped:   false,
				PID:       state.PID,
				Message:   "relay process already exited; state cleared",
				StatePath: statePath,
			}, nil
		}
		return nil, fmt.Errorf("failed to stop relay pid %d: %w", state.PID, err)
	}

	_ = removeStateFile()
	return &StopResult{
		Stopped:   true,
		PID:       state.PID,
		Message:   "relay stopped",
		StatePath: statePath,
	}, nil
}

// IsPortInUse checks whether localhost port is already bound.
func IsPortInUse(port int) (bool, error) {
	if port <= 0 {
		return false, fmt.Errorf("invalid port %d", port)
	}

	checks := []struct {
		network string
		addr    string
	}{
		{network: "tcp", addr: fmt.Sprintf(":%d", port)},
		{network: "tcp4", addr: fmt.Sprintf("127.0.0.1:%d", port)},
		{network: "tcp6", addr: fmt.Sprintf("[::1]:%d", port)},
	}

	for _, check := range checks {
		ln, err := net.Listen(check.network, check.addr)
		if err == nil {
			_ = ln.Close()
			continue
		}
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "address already in use") {
			return true, nil
		}
	}

	return false, nil
}

// PortOwner returns best-effort listener owner and whether it looks like ai-happy-design/ahd-figma.
func PortOwner(port int) (string, bool) {
	if port <= 0 {
		return "", false
	}
	lsof, err := exec.LookPath("lsof")
	if err != nil {
		return "", false
	}
	cmd := exec.Command(lsof, "-nP", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN")
	out, err := cmd.Output()
	if err != nil {
		return "", false
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return "", false
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return "", false
	}

	command := fields[0]
	pid := fields[1]
	summary := fmt.Sprintf("%s (pid %s)", command, pid)
	lower := strings.ToLower(command)
	return summary, looksLikeAHDProcess(lower)
}

func looksLikeAHDProcess(lower string) bool {
	return strings.Contains(lower, "ai-happy-design") || strings.Contains(lower, "ahd-figma")
}

func processCommand(pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("invalid pid %d", pid)
	}
	if runtime.GOOS == "windows" {
		task := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		out, err := task.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	ps := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=")
	out, err := ps.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// TailLogs returns last N lines from managed relay log file.
func TailLogs(lines int) (string, string, error) {
	if lines <= 0 {
		lines = 80
	}
	logPath, err := LogFilePath()
	if err != nil {
		return "", "", err
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		return "", logPath, err
	}
	all := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	// Drop trailing empty line for cleaner output.
	if len(all) > 0 && all[len(all)-1] == "" {
		all = all[:len(all)-1]
	}
	if len(all) > lines {
		all = all[len(all)-lines:]
	}
	return strings.Join(all, "\n"), logPath, nil
}

// InstallLaunchAgent installs and loads a user launch agent (macOS).
func InstallLaunchAgent(host string, port int) (*InstallAgentResult, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("relay install-agent is only supported on macOS")
	}
	if !IsLocalHost(host) {
		return nil, fmt.Errorf("relay install-agent supports localhost hosts only, got %q", host)
	}
	if port <= 0 {
		return nil, fmt.Errorf("invalid port %d", port)
	}

	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("unable to locate executable: %w", err)
	}
	logPath, err := LogFilePath()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to resolve home directory: %w", err)
	}
	agentsDir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		return nil, fmt.Errorf("unable to create launch agents directory: %w", err)
	}
	plistPath := filepath.Join(agentsDir, fmt.Sprintf("%s.plist", agentLabel))
	plist := buildLaunchdPlist(exe, host, port, logPath)

	if err := os.WriteFile(plistPath, []byte(plist), 0o644); err != nil {
		return nil, fmt.Errorf("failed to write launch agent plist: %w", err)
	}

	// Try reloading idempotently.
	_ = exec.Command("launchctl", "unload", plistPath).Run()
	loadCmd := exec.Command("launchctl", "load", plistPath)
	if out, err := loadCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("launchctl load failed: %v (%s)", err, strings.TrimSpace(string(out)))
	}

	return &InstallAgentResult{
		Label:     agentLabel,
		PlistPath: plistPath,
		Loaded:    true,
		LogPath:   logPath,
	}, nil
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(value)
}

func buildLaunchdPlist(executable, host string, port int, logPath string) string {
	exe := xmlEscape(executable)
	escapedHost := xmlEscape(host)
	portValue := strconv.Itoa(port)
	log := xmlEscape(logPath)

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>%s</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>
    <string>ws</string>
  </array>
  <key>EnvironmentVariables</key>
  <dict>
    <key>PORT</key>
    <string>%s</string>
    <key>SERVER_HOST</key>
    <string>%s</string>
  </dict>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>%s</string>
  <key>StandardErrorPath</key>
  <string>%s</string>
</dict>
</plist>
`, agentLabel, exe, portValue, escapedHost, log, log)
}
