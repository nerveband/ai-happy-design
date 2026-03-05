package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
)

var interpolationPattern = regexp.MustCompile(`\$\{\{steps\.([A-Za-z0-9_-]+)\.result(?:\.([A-Za-z0-9_.-]+))?\}\}`)

// ResolveBatchOp resolves step interpolation references against completed steps.
func ResolveBatchOp(op commoncli.BatchOp, completed []commoncli.BatchStep) (commoncli.BatchOp, error) {
	if len(op.Params) == 0 {
		return op, nil
	}
	index := map[string]commoncli.BatchStep{}
	for _, step := range completed {
		if step.Name != "" {
			index[step.Name] = step
		}
	}
	resolved, err := resolveValue(op.Params, index)
	if err != nil {
		return op, err
	}
	params, ok := resolved.(map[string]any)
	if !ok {
		return op, fmt.Errorf("resolved batch params are not an object")
	}
	op.Params = params
	return op, nil
}

func resolveValue(value any, steps map[string]commoncli.BatchStep) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, nested := range typed {
			resolved, err := resolveValue(nested, steps)
			if err != nil {
				return nil, err
			}
			out[key] = resolved
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for i, nested := range typed {
			resolved, err := resolveValue(nested, steps)
			if err != nil {
				return nil, err
			}
			out[i] = resolved
		}
		return out, nil
	case string:
		return resolveString(typed, steps)
	default:
		return value, nil
	}
}

func resolveString(value string, steps map[string]commoncli.BatchStep) (any, error) {
	matches := interpolationPattern.FindAllStringSubmatchIndex(value, -1)
	if len(matches) == 0 {
		return value, nil
	}
	if len(matches) == 1 && matches[0][0] == 0 && matches[0][1] == len(value) {
		match := interpolationPattern.FindStringSubmatch(value)
		return lookupStepValue(match[1], match[2], steps)
	}
	var builder strings.Builder
	last := 0
	for _, loc := range matches {
		builder.WriteString(value[last:loc[0]])
		match := interpolationPattern.FindStringSubmatch(value[loc[0]:loc[1]])
		resolved, err := lookupStepValue(match[1], match[2], steps)
		if err != nil {
			return nil, err
		}
		switch scalar := resolved.(type) {
		case string:
			builder.WriteString(scalar)
		case float64, bool, int:
			builder.WriteString(fmt.Sprint(scalar))
		default:
			return nil, fmt.Errorf("cannot embed non-scalar result from step %q into string", match[1])
		}
		last = loc[1]
	}
	builder.WriteString(value[last:])
	return builder.String(), nil
}

func lookupStepValue(name, path string, steps map[string]commoncli.BatchStep) (any, error) {
	step, ok := steps[name]
	if !ok {
		return nil, fmt.Errorf("unknown step reference %q", name)
	}
	if path == "" {
		return step.Result, nil
	}
	current := step.Result
	for _, part := range strings.Split(path, ".") {
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("step %q has no result path %q", name, path)
		}
		current, ok = obj[part]
		if !ok {
			return nil, fmt.Errorf("step %q has no result path %q", name, path)
		}
	}
	return current, nil
}
