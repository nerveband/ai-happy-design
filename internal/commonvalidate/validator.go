package commonvalidate

import (
	"fmt"
	"math"
	"net/url"
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
		normalized, warnings, errors := validateValue(*param, value, cwd, canonical)
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

	applyCommandSpecificRules(command, &out)
	return out
}

func applyCommandSpecificRules(command *commonschema.Command, out *Result) {
	if command == nil {
		return
	}
	switch command.Name {
	case "document.new":
		validateDocumentNewRules(out)
	}
}

func validateDocumentNewRules(out *Result) {
	artboards := int(numberOrDefault(out.Params["artboards"], 1))
	if artboards < 1 {
		return
	}
	rowsOrCols := int(numberOrDefault(out.Params["artboardRowsOrCols"], 1))
	layout := stringOrDefault(out.Params["artboardLayout"], "GridByRow")
	if rowsOrCols < 1 {
		out.Errors = append(out.Errors, Issue{
			Code:    "VALIDATION_ERROR",
			Field:   "artboardRowsOrCols",
			Message: "artboardRowsOrCols must be at least 1",
		})
		return
	}

	switch layout {
	case "Row", "Column", "RLRow":
		if rowsOrCols != 1 {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "artboardRowsOrCols",
				Message: fmt.Sprintf("artboardRowsOrCols must be 1 for %s layout", layout),
			})
		}
	case "GridByRow", "GridByCol", "RLGridByRow", "RLGridByCol":
		maxRowsOrCols := artboards - 1
		if maxRowsOrCols < 1 {
			maxRowsOrCols = 1
		}
		if rowsOrCols > maxRowsOrCols {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "artboardRowsOrCols",
				Message: fmt.Sprintf("artboardRowsOrCols must be between 1 and %d for %d artboards using %s layout", maxRowsOrCols, artboards, layout),
			})
		}
	}
}

func validateValue(param commonschema.Param, value any, cwd, fieldPath string) (any, []Issue, []Issue) {
	var warnings []Issue
	var errors []Issue

	switch param.Type {
	case "string":
		str, ok := value.(string)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(fieldPath, param.Type, value)}
		}
		if strings.ContainsRune(str, 0) || hasControlChars(str) {
			return nil, warnings, []Issue{{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: "string contains control characters",
			}}
		}
		if param.OpaqueIdentifier {
			errors = append(errors, validateOpaqueIdentifier(fieldPath, str)...)
		}
		if param.Pattern != "" {
			re, err := regexp.Compile(param.Pattern)
			if err == nil && !re.MatchString(str) {
				errors = append(errors, Issue{
					Code:    "VALIDATION_ERROR",
					Field:   fieldPath,
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
							Field:   fieldPath,
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
						Field:   fieldPath,
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
					Field:   fieldPath,
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
			return nil, warnings, []Issue{typeMismatch(fieldPath, param.Type, value)}
		}
		if param.Minimum != nil && number < *param.Minimum {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: fmt.Sprintf("value %.2f is below minimum %.2f", number, *param.Minimum),
			})
		}
		if param.Maximum != nil && number > *param.Maximum {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: fmt.Sprintf("value %.2f is above maximum %.2f", number, *param.Maximum),
			})
		}
		return number, warnings, errors

	case "boolean":
		boolean, ok := value.(bool)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(fieldPath, param.Type, value)}
		}
		return boolean, warnings, errors

	case "array":
		items, ok := value.([]any)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(fieldPath, param.Type, value)}
		}
		if param.MinItems != nil && len(items) < *param.MinItems {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: fmt.Sprintf("array must contain at least %d item(s)", *param.MinItems),
			})
		}
		if param.MaxItems != nil && len(items) > *param.MaxItems {
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: fmt.Sprintf("array must contain at most %d item(s)", *param.MaxItems),
			})
		}
		normalized := make([]any, 0, len(items))
		for index, item := range items {
			itemField := fmt.Sprintf("%s[%d]", fieldPath, index)
			if param.Items != nil {
				value, childWarnings, childErrors := validateValue(*param.Items, item, cwd, itemField)
				warnings = append(warnings, childWarnings...)
				errors = append(errors, childErrors...)
				if len(childErrors) == 0 {
					normalized = append(normalized, value)
				}
				continue
			}
			errors = append(errors, validateGenericValue(itemField, item)...)
			normalized = append(normalized, item)
		}
		if len(errors) > 0 {
			return nil, warnings, errors
		}
		return normalized, warnings, errors

	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			return nil, warnings, []Issue{typeMismatch(fieldPath, param.Type, value)}
		}
		if len(param.Fields) == 0 {
			errors = append(errors, validateGenericValue(fieldPath, object)...)
			if len(errors) > 0 {
				return nil, warnings, errors
			}
			return object, warnings, errors
		}

		normalized := map[string]any{}
		for key, childValue := range object {
			child := lookupField(param.Fields, key)
			if child == nil {
				if !param.AllowUnknown {
					errors = append(errors, Issue{
						Code:    "VALIDATION_ERROR",
						Field:   fmt.Sprintf("%s.%s", fieldPath, key),
						Message: fmt.Sprintf("unknown field %q", key),
					})
				}
				continue
			}
			childPath := fmt.Sprintf("%s.%s", fieldPath, child.Name)
			value, childWarnings, childErrors := validateValue(*child, childValue, cwd, childPath)
			warnings = append(warnings, childWarnings...)
			errors = append(errors, childErrors...)
			if len(childErrors) == 0 {
				normalized[child.Name] = value
			}
		}

		for _, child := range param.Fields {
			if !child.Required {
				continue
			}
			if _, ok := normalized[child.Name]; ok {
				continue
			}
			errors = append(errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   fmt.Sprintf("%s.%s", fieldPath, child.Name),
				Message: fmt.Sprintf("missing required field %q", child.Name),
			})
		}

		if len(errors) > 0 {
			return nil, warnings, errors
		}
		return normalized, warnings, errors

	default:
		errors = append(errors, validateGenericValue(fieldPath, value)...)
		if len(errors) > 0 {
			return nil, warnings, errors
		}
		return value, warnings, errors
	}
}

func typeMismatch(field, want string, value any) Issue {
	return Issue{
		Code:    "VALIDATION_ERROR",
		Field:   field,
		Message: fmt.Sprintf("param %q must be %s, got %T", field, want, value),
	}
}

func validateGenericValue(fieldPath string, value any) []Issue {
	switch typed := value.(type) {
	case string:
		if strings.ContainsRune(typed, 0) || hasControlChars(typed) {
			return []Issue{{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: "string contains control characters",
			}}
		}
	case []any:
		var errors []Issue
		for index, item := range typed {
			errors = append(errors, validateGenericValue(fmt.Sprintf("%s[%d]", fieldPath, index), item)...)
		}
		return errors
	case map[string]any:
		var errors []Issue
		for key, item := range typed {
			errors = append(errors, validateGenericValue(fmt.Sprintf("%s.%s", fieldPath, key), item)...)
		}
		return errors
	}
	return nil
}

func validateOpaqueIdentifier(fieldPath, value string) []Issue {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return []Issue{{
			Code:    "VALIDATION_ERROR",
			Field:   fieldPath,
			Message: "identifier cannot be blank",
		}}
	}
	if trimmed != value {
		return []Issue{{
			Code:    "VALIDATION_ERROR",
			Field:   fieldPath,
			Message: "identifier must not have leading or trailing whitespace",
		}}
	}
	lower := strings.ToLower(value)
	if strings.Contains(lower, "://") || strings.HasPrefix(lower, "file:") {
		return []Issue{{
			Code:    "VALIDATION_ERROR",
			Field:   fieldPath,
			Message: "identifier must not be a URL",
		}}
	}
	if strings.ContainsAny(value, "?#") {
		return []Issue{{
			Code:    "VALIDATION_ERROR",
			Field:   fieldPath,
			Message: "identifier must not contain query or fragment syntax",
		}}
	}
	for _, snippet := range []string{"../", "..\\", "?=", "&", "=%", "%2e", "%2f", "%5c", "%00", "%3f", "%26", "%3d"} {
		if strings.Contains(lower, snippet) {
			return []Issue{{
				Code:    "VALIDATION_ERROR",
				Field:   fieldPath,
				Message: "identifier contains encoded or query-like syntax",
			}}
		}
	}
	return nil
}

func lookupField(fields []commonschema.Param, name string) *commonschema.Param {
	needle := strings.ToLower(strings.TrimSpace(name))
	for i := range fields {
		if strings.ToLower(fields[i].Name) == needle {
			return &fields[i]
		}
		for _, alias := range fields[i].Aliases {
			if strings.ToLower(alias) == needle {
				return &fields[i]
			}
		}
	}
	return nil
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

func normalizePotentialEncodedPath(value string) (string, bool) {
	unescaped, err := url.PathUnescape(value)
	if err != nil {
		return value, false
	}
	return unescaped, unescaped != value
}

func numberOrDefault(value any, fallback float64) float64 {
	number, ok := toFloat64(value)
	if !ok {
		return fallback
	}
	return number
}

func stringOrDefault(value any, fallback string) string {
	typed, ok := value.(string)
	if !ok || strings.TrimSpace(typed) == "" {
		return fallback
	}
	return typed
}
