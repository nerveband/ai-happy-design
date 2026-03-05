package inspect

// SelectorFor maps inspect commands to the plugin selector contract.
func SelectorFor(command string) string {
	return "ahd.inspect"
}
