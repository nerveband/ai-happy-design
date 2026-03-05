package commoncli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// WriteJSON writes a single JSON payload.
func WriteJSON(value any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}

// WriteNDJSON writes one JSON record per line.
func WriteNDJSON(values []any) error {
	enc := json.NewEncoder(os.Stdout)
	for _, value := range values {
		if err := enc.Encode(value); err != nil {
			return err
		}
	}
	return nil
}

// WriteText writes a best-effort text summary for interactive use.
func WriteText(value any) error {
	switch v := value.(type) {
	case string:
		_, err := fmt.Fprintln(os.Stdout, v)
		return err
	case Envelope:
		if !v.OK && v.Error != nil {
			_, err := fmt.Fprintf(os.Stdout, "[%s] %s\n", v.Error.Code, v.Error.Message)
			return err
		}
		_, err := fmt.Fprintf(os.Stdout, "%s ok\n", strings.TrimSpace(v.Command))
		return err
	case BatchEnvelope:
		_, err := fmt.Fprintf(os.Stdout, "batch ok=%t succeeded=%d failed=%d\n", v.OK, v.Summary.Succeeded, v.Summary.Failed)
		return err
	default:
		data, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(os.Stdout, string(data))
		return err
	}
}
