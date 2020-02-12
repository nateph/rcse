package common

// CommandResult contains various information about what a command returned.
type CommandResult struct {
	// stdout from the command.
	Stdout []byte

	// Return code of the command executed.
	ReturnCode int

	// Host command ran on
	Host string
}
