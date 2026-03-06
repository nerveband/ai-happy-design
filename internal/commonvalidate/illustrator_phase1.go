package commonvalidate

import (
	"fmt"
	"math"
)

func validatePreferenceSetRules(out *Result) {
	valueType := stringOrDefault(out.Params["valueType"], "")
	value, _ := out.Params["value"].(map[string]any)
	if value == nil {
		out.Errors = append(out.Errors, Issue{
			Code:    "VALIDATION_ERROR",
			Field:   "value",
			Message: "value wrapper is required",
		})
		return
	}

	switch valueType {
	case "boolean":
		if _, ok := value["boolean"]; !ok {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "value.boolean",
				Message: "boolean preference writes require value.boolean",
			})
		}
	case "integer":
		number, ok := value["number"].(float64)
		if !ok {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "value.number",
				Message: "integer preference writes require value.number",
			})
			return
		}
		if math.Round(number) != number {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "value.number",
				Message: "integer preference writes require a whole number",
			})
		}
	case "real":
		if _, ok := value["number"].(float64); !ok {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "value.number",
				Message: "real preference writes require value.number",
			})
		}
	case "string":
		if _, ok := value["string"].(string); !ok {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "value.string",
				Message: "string preference writes require value.string",
			})
		}
	default:
		out.Errors = append(out.Errors, Issue{
			Code:    "VALIDATION_ERROR",
			Field:   "valueType",
			Message: fmt.Sprintf("unsupported preference valueType %q", valueType),
		})
	}
}

func validateIllustratorColorCreationRules(out *Result) {
	colorMode := stringOrDefault(out.Params["colorMode"], "")
	switch colorMode {
	case "RGB":
		if _, ok := out.Params["hex"]; !ok {
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   "hex",
				Message: "RGB color payloads require a hex field",
			})
		}
	case "CMYK":
		for _, field := range []string{"cyan", "magenta", "yellow", "black"} {
			if _, ok := out.Params[field]; ok {
				continue
			}
			out.Errors = append(out.Errors, Issue{
				Code:    "VALIDATION_ERROR",
				Field:   field,
				Message: fmt.Sprintf("CMYK color payloads require %s", field),
			})
		}
	default:
		out.Errors = append(out.Errors, Issue{
			Code:    "VALIDATION_ERROR",
			Field:   "colorMode",
			Message: fmt.Sprintf("unsupported colorMode %q", colorMode),
		})
	}
}
