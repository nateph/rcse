package cliconfig

import (
	"fmt"
)

// CommandResult contains various information about what a command returned.
type CommandResult struct {
	// The command that was ran.
	CommandRan string

	// stderr from the command.
	Stderr []byte
	// stdout from the command.
	Stdout []byte

	// Return code of the command executed.
	ReturnCode int

	// Host command ran on
	Host string
}

// PrintHostOutput formats the host and stdout nicely.
func (cr *CommandResult) PrintHostOutput() {
	if len(cr.Stderr) > 0 {
		fmt.Printf("Command '%s' failed, Stderr returned:\n%s", cr.CommandRan, cr.Stderr)
	} else {
		fmt.Printf("----- %s -----\n%s\n%s\n", cr.Host, cr.CommandRan, cr.Stdout)
	}
}

// CommandOptions contains information on the the command to be ran
type CommandOptions struct {
	// Command that will be ran
	CommandToRun string

	// Which host it will be ran on
	Host string

	// To execute the command as sudo or not
	Sudo bool

	// Whether or not to verify host keys
	IgnoreHostkeyCheck bool
}

// RunCommand is a wrapper around establishing the ssh connection and then
// calling the RunSSHCommand
func (co *CommandOptions) RunCommand() {
	// ignoreHostkeyCheck is a persistent flag set in the root command
	sshSession := EstablishSSHConnection(co.Host, co.IgnoreHostkeyCheck)
	defer sshSession.Close()

	result := RunSSHCommand(co.CommandToRun, co.Host, sshSession)

	result.PrintHostOutput()
}
