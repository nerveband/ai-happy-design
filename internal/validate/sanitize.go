package validate

import "strings"

// promptInjectionPatterns contains strings that may indicate prompt injection
// attempts in response data flowing back to LLM agents.
var promptInjectionPatterns = []string{
	"<|system|>",
	"<|user|>",
	"<|assistant|>",
	"```system",
	"IGNORE PREVIOUS INSTRUCTIONS",
	"<system>",
	"</system>",
}

// SanitizeResponse strips potential prompt-injection markers from API response data.
// This prevents node names or text content from hijacking an LLM agent's context.
func SanitizeResponse(s string) string {
	result := s
	for _, d := range promptInjectionPatterns {
		result = strings.ReplaceAll(result, d, "")
	}
	return result
}

// SanitizeResponseMap recursively sanitizes all string values in a map.
func SanitizeResponseMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch v := v.(type) {
		case string:
			m[k] = SanitizeResponse(v)
		case map[string]interface{}:
			m[k] = SanitizeResponseMap(v)
		case []interface{}:
			m[k] = sanitizeResponseSlice(v)
		}
	}
	return m
}

func sanitizeResponseSlice(s []interface{}) []interface{} {
	for i, v := range s {
		switch v := v.(type) {
		case string:
			s[i] = SanitizeResponse(v)
		case map[string]interface{}:
			s[i] = SanitizeResponseMap(v)
		case []interface{}:
			s[i] = sanitizeResponseSlice(v)
		}
	}
	return s
}
