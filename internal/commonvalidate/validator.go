package commonvalidate

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
)

// Issue is a machine-readable validation warning or error.
type Issue struct {
	Code    string `json:"code"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Fix     any    `json:"fix,omitempty"`
	Applied bool   `json:"applied,omitempty"`
}

// Result is the normalized validation output.
type Result struct {
	Params   map[string]any `json:"params"`
	Warnings []Issue        `json:"warnings,omitempty"`
	Errors   []Issue        `json:"errors,omitempty"`
}

// ValidateCommand validates and normalizes params against a schema.
func ValidateCommand(command *commonschema.Command, params map[string]any, cwd string) Result {
	out := Result{Params: map[string]any{}}
	if command == nil {
		out.Errors = append(out.Errors, Issue{
			Code:    "UNSUPPORTED_COMMAND",
			Message: "unknown command",
		})
		return out
	}

	for key, value := range params {
		param := commonschema.LookupParam(command, key)
		if param == nil {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   key,
				Message: fmt.Sprintf("unknown param %q for %s", key, command.Name),
			})
			continue
		}
		canonical := param.Name
		if key != canonical {
			out.Warnings = append(out.Warnings, Issue{
				Code:    "ALIASED_PARAM",
				Field:   key,
				Message: fmt.Sprintf("normalized alias %q to %q", key, canonical),
				Fix:     canonical,
				Applied: true,
			})
		}
		normalized, warnings, errors := validateValue(*param, value, cwd)
		out.Warnings = append(out.Warnings, warnings...)
		out.Errors = append(out.Errors, errors...)
		if len(errors) == 0 {
			out.Params[canonical] = normalized
		}
	}

	for _, param := range command.Params {
		if !param.Required {
			continue
		}
		if _, ok := out.Params[param.Name]; ok {
			continue
		}
		out.Errors = append(out.Errors, Issue{
			Code:    "VALIDATION_ERROR",
			Field:   param.Name,
			Message: fmt.Sprintf("missing required param %q", param.Name),
		})
	}

	return out
}

func validateValue(param commonschema.Param, value any, cwd string) (any, []Issue, []Issue) {
	var warnings []Issue
	var errors []Issue

	switch param.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(param, value)}
		}
		if strings.ContainsRune(str, 0) || hasControlChars(str) {
			return nil, warnings, []Issue{{
				Code:    "VALIDATION_ERROR",
				Field:   param.Name,
				Message: "string contains control characters",
			}}
		}
		if param.Pattern != "" {
			re, err := regexp.Compile(param.Pattern)
			if err == nil && !re.MatchString(str) {
				errors = append(errors, Issue{
					Code:    "VALIDATION_ERROR",
					Field:   param.Name,
					Message: fmt.Sprintf("value %q does not match pattern %s", str, param.Pattern),
				})
			}
		}
		if len(param.Enum) > 0 {
			match := false
			for _, item := range param.Enum {
				if item == str {
					match = true
					break
				}
			}
			if !match {
				if param.LowRiskFuzzy {
					best := bestEnumMatch(str, param.Enum)
					if best != "" {
						warnings = append(warnings, Issue{
							Code:    "FUZZY_ENUM",
							Field:   param.Name,
							Message: fmt.Sprintf("normalized %q to %q", str, best),
							Fix:     best,
							Applied: true,
						})
						str = best
						match = true
					}
				}
				if !match {
					errors = append(errors, Issue{
						Code:    "VALIDATION_ERROR",
						Field:   param.Name,
						Message: fmt.Sprintf("value %q is not in enum %v", str, param.Enum),
					})
				}
			}
		}
		if param.SafePath {
			safe, err := SafePath(cwd, str)
			if err != nil {
				errors = append(errors, Issue{
					Code:    "VALIDATION_ERROR",
					Field:   param.Name,
					Message: err.Error(),
				})
			} else {
				str = safe
			}
		}
		return str, warnings, errors

	case "number":
		number, ok := toFloat64(value)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(param, value)}
		}
		if param.Minimum != nil && number < *param.Minimum {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   param.Name,
				Message: fmt.Sprintf("value %.2f is below minimum %.2f", number, *param.Minimum),
			})
		}
		if param.Maximum != nil && number > *param.Maximum {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   param.Name,
				Message: fmt.Sprintf("value %.2f is above maximum %.2f", number, *param.Maximum),
			})
		}
		return number, warnings, errors

	case "boolean":
		boolean, ok := value.(bool)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(param, value)}
		}
		return boolean, warnings, errors

	case "array":
		items, ok := value.([]any)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(param, value)}
		}
		return items, warnings, errors

	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(param, value)}
		}
		return object, warnings, errors

	default:
		return value, warnings, errors
	}
}

func typeMismatch(param commonschema.Param, value any) Issue {
	return Issue{
		Code:    "VALIDATION_ERROR",
		Field:   param.Name,
		Message: fmt.Sprintf("param %q must be %s, got %T", param.Name, param.Type, value),
	}
}

func hasControlChars(value string) bool {
	for _, r := range value {
		if r < 32 && r != '\n' && r != '\r' && r != '\t' {
			return true
		}
	}
	return false
}

func toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case jsonNumber:
		num, err := strconv.ParseFloat(string(v), 64)
		return num, err == nil
	case string:
		num, err := strconv.ParseFloat(v, 64)
		return num, err == nil
	default:
		return 0, false
	}
}

type jsonNumber string

func bestEnumMatch(value string, options []string) string {
	best := ""
	bestDistance := math.MaxInt
	needle := strings.ToLower(strings.TrimSpace(value))
	for _, option := range options {
		dist := levenshtein(needle, strings.ToLower(option))
		if dist < bestDistance {
			bestDistance = dist
			best = option
		}
	}
	if bestDistance <= 2 {
		return best
	}
	return ""
}

func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	if a == "" {
		return len(b)
	}
	if b == "" {
		return len(a)
	}

	previous := make([]int, len(b)+1)
	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(a); i++ {
		current := make([]int, len(b)+1)
		current[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			current[j] = min3(
				current[j-1]+1,
				previous[j]+1,
				previous[j-1]+cost,
			)
		}
		previous = current
	}
	return previous[len(b)]
}

func min3(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
