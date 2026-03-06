package commands

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestPhase1SchemaCommandsHaveExecutorCases(t *testing.T) {
	t.Parallel()

	schemaSource, err := os.ReadFile("../schema/phase1.go")
	if err != nil {
		t.Fatalf("read phase1 schema source: %v", err)
	}
	commandSource, err := os.ReadFile("phase1.go")
	if err != nil {
		t.Fatalf("read phase1 command source: %v", err)
	}

	schemaNames := extractMatches(string(schemaSource), regexp.MustCompile(`Name:\s+"([^"]+)"`))
	caseNames := extractMatches(string(commandSource), regexp.MustCompile(`case "([^"]+)"`))

	missing := make([]string, 0)
	for name := range schemaNames {
		if _, ok := caseNames[name]; !ok {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("phase1 schema commands missing executor cases: %v", missing)
	}
}

func TestPhase1SchemaExcludesUnstableRuntimeCommands(t *testing.T) {
	t.Parallel()

	schemaSource, err := os.ReadFile("../schema/phase1.go")
	if err != nil {
		t.Fatalf("read phase1 schema source: %v", err)
	}
	schemaNames := extractMatches(string(schemaSource), regexp.MustCompile(`Name:\s+"([^"]+)"`))
	for _, name := range []string{
		"document.write_as_library",
		"page_item.bring_in_perspective",
		"trace.preset.store",
	} {
		if _, ok := schemaNames[name]; ok {
			t.Fatalf("phase1 schema unexpectedly still exposes %s", name)
		}
	}
}

func extractMatches(source string, pattern *regexp.Regexp) map[string]struct{} {
	matches := pattern.FindAllStringSubmatch(source, -1)
	out := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		if !strings.Contains(match[1], ".") {
			continue
		}
		out[match[1]] = struct{}{}
	}
	return out
}
