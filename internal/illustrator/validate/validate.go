package validate

import (
	"os"

	"github.com/nerveband/ai-happy-design/internal/commonschema"
	"github.com/nerveband/ai-happy-design/internal/commonvalidate"
)

// ValidateCommand applies common schema validation using the current working directory.
func ValidateCommand(command *commonschema.Command, params map[string]any) commonvalidate.Result {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return commonvalidate.ValidateCommand(command, params, cwd)
}
