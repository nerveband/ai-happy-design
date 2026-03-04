package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/schema"
)

// Issue represents a single validation warning or error.
type Issue struct {
	Step     int         `json:"step"`
	Name     string      `json:"name,omitempty"`
	Phase    string      `json:"phase"`
	Code     string      `json:"code"`
	Param    string      `json:"param,omitempty"`
	Message  string      `json:"message"`
	Got      interface{} `json:"got,omitempty"`
	Expected interface{} `json:"expected,omitempty"`
	Fix      interface{} `json:"fix,omitempty"`
	Applied  bool        `json:"applied"`
}

// Result holds all validation results for a batch.
type Result struct {
	Warnings []Issue `json:"warnings,omitempty"`
	Errors   []Issue `json:"errors,omitempty"`
	Fixed    int     `json:"fixed"`
	Blocked  int     `json:"blocked"`
}

// ValidateBatch validates all operations against registered schemas.
// It mutates ops in-place when auto-fixes are applied (warn+fix mode).
func ValidateBatch(ops []map[string]interface{}) Result {
	var result Result
	commands := schema.Commands()

	for i, op := range ops {
		cmdRaw, _ := op["command"].(string)
		name, _ := op["name"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			params = map[string]interface{}{}
		}

		// Check command exists
		s := schema.Lookup(cmdRaw)
		if s == nil {
			fix, found := FuzzyMatchCommand(cmdRaw, commands)
			issue := Issue{
				Step: i, Name: name, Phase: "schema", Code: "UNKNOWN_COMMAND",
				Message: fmt.Sprintf("unknown command: %s", cmdRaw),
				Got:     cmdRaw,
			}
			if found {
				issue.Fix = fix
				issue.Message += fmt.Sprintf(". Did you mean: %s?", fix)
				issue.Applied = true
				op["command"] = fix
				s = schema.Lookup(fix)
				result.Fixed++
			} else {
				result.Blocked++
			}
			result.Warnings = append(result.Warnings, issue)
			if s == nil {
				continue
			}
		}

		// Validate each param against schema
		for key, val := range params {
			p := schema.LookupParam(s, key)
			if p == nil {
				// Unknown param — try fuzzy match
				bestMatch := ""
				bestDist := 4
				for _, sp := range s.Params {
					d := levenshtein(strings.ToLower(key), strings.ToLower(sp.Name))
					if d < bestDist {
						bestDist = d
						bestMatch = sp.Name
					}
				}
				issue := Issue{
					Step: i, Name: name, Phase: "schema", Code: "UNKNOWN_PARAM",
					Param: key, Got: key,
					Message: fmt.Sprintf("unknown param '%s' for %s", key, s.Command),
				}
				if bestDist <= 3 {
					issue.Fix = bestMatch
					issue.Message += fmt.Sprintf(". Did you mean: %s?", bestMatch)
				}
				result.Warnings = append(result.Warnings, issue)
				continue
			}

			// Type checking + bounds + enum + pattern
			issues := validateParam(i, name, p, val, params)
			for _, issue := range issues {
				if issue.Applied {
					result.Fixed++
				} else if issue.Code != "" && issue.Fix == nil {
					result.Blocked++
				}
				result.Warnings = append(result.Warnings, issue)
			}
		}

		// Check required params
		for _, p := range s.Params {
			if !p.Required {
				continue
			}
			if _, exists := params[p.Name]; exists {
				continue
			}
			// Check aliases too
			found := false
			for _, alias := range p.Aliases {
				if _, exists := params[alias]; exists {
					found = true
					break
				}
			}
			if !found {
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "schema", Code: "REQUIRED_MISSING",
					Param:   p.Name,
					Message: fmt.Sprintf("required param '%s' missing", p.Name),
				})
				result.Blocked++
			}
		}

		// Check auto-fix dependencies
		for _, p := range s.Params {
			if p.AutoFix == "" {
				continue
			}
			if _, exists := params[p.Name]; !exists {
				continue // param not set, skip dependency check
			}
			// Parse "key:value"
			parts := strings.SplitN(p.AutoFix, ":", 2)
			if len(parts) != 2 {
				continue
			}
			depKey, depVal := parts[0], parts[1]
			if _, exists := params[depKey]; !exists {
				params[depKey] = depVal
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "schema", Code: "DEPENDENCY_MISSING",
					Param:   depKey,
					Message: fmt.Sprintf("auto-set %s=%s (required when %s is set)", depKey, depVal, p.Name),
					Fix:     depVal,
					Applied: true,
				})
				result.Fixed++
			}
		}

		// Write back params (may have been mutated)
		op["params"] = params
	}
	return result
}

func validateParam(step int, name string, p *schema.Param, val interface{}, params map[string]interface{}) []Issue {
	var issues []Issue

	switch p.Type {
	case "number":
		num, ok := toFloat64(val)
		if !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a number, got %T", p.Name, val),
			})
			return issues
		}

		if p.Min != nil && num < *p.Min {
			params[p.Name] = *p.Min
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "BELOW_MIN",
				Param: p.Name, Got: num,
				Expected: map[string]interface{}{"min": *p.Min, "max": p.Max},
				Fix:      *p.Min,
				Applied:  true,
				Message:  fmt.Sprintf("%s must be >= %.0f, got %.0f", p.Name, *p.Min, num),
			})
		} else if p.Max != nil && num > *p.Max {
			params[p.Name] = *p.Max
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "ABOVE_MAX",
				Param: p.Name, Got: num,
				Expected: map[string]interface{}{"min": p.Min, "max": *p.Max},
				Fix:      *p.Max,
				Applied:  true,
				Message:  fmt.Sprintf("%s must be <= %.0f, got %.0f", p.Name, *p.Max, num),
			})
		}

	case "string":
		str, ok := val.(string)
		if !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a string, got %T", p.Name, val),
			})
			return issues
		}

		// Enum check with fuzzy matching
		if len(p.Enum) > 0 {
			found := false
			for _, e := range p.Enum {
				if e == str {
					found = true
					break
				}
			}
			if !found {
				fix, matched := FuzzyMatchEnum(str, p.Enum)
				issue := Issue{
					Step: step, Name: name, Phase: "schema", Code: "ENUM_INVALID",
					Param: p.Name, Got: str,
					Expected: map[string]interface{}{"enum": p.Enum},
					Message:  fmt.Sprintf("%s must be one of %v, got '%s'", p.Name, p.Enum, str),
				}
				if matched {
					issue.Fix = fix
					issue.Applied = true
					params[p.Name] = fix
				}
				issues = append(issues, issue)
			}
		}

		// Pattern check (hex colors, node IDs)
		if p.Pattern != "" {
			re, err := regexp.Compile(p.Pattern)
			if err == nil && !re.MatchString(str) {
				issue := Issue{
					Step: step, Name: name, Phase: "schema", Code: "PATTERN_MISMATCH",
					Param: p.Name, Got: str,
					Expected: map[string]interface{}{"pattern": p.Pattern},
					Message:  fmt.Sprintf("%s doesn't match expected pattern, got '%s'", p.Name, str),
				}
				// Try named color resolution
				if hex := ResolveNamedColor(str); hex != "" {
					issue.Fix = hex
					issue.Applied = true
					params[p.Name] = hex
				}
				issues = append(issues, issue)
			}
		}

	case "boolean":
		if _, ok := val.(bool); !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a boolean, got %T", p.Name, val),
			})
		}
	}
	return issues
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
