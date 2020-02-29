package cliconfig

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
		fmt.Printf("Command '%s' failed stderr below:\n%s", cr.CommandRan, cr.Stderr)
	} else {
		fmt.Printf("----- %s -----\n%s\n%s\n", cr.Host, cr.CommandRan, cr.Stdout)
	}
}

// CommandOptions contains information on the the command to be ran
type CommandOptions struct {
	// Command that will be ran
	CommandsToRun []string
	// Which host it will be ran on
	Host string
	// Whether or not to verify host keys
	IgnoreHostkeyCheck bool
	// To execute the command as sudo or not
	Sudo bool
	// User to execute as
	User string
	// Password for User
	Password string
}

// RunCommands is a wrapper around establishing the ssh connection and then
// calling the RunSSHCommand
func (co *CommandOptions) RunCommands() []CommandResult {
	sshClient := EstablishSSHConnection(co.User, co.Password, co.Host, co.IgnoreHostkeyCheck)
	defer sshClient.Close()

	var CommandResults []CommandResult

	for _, command := range co.CommandsToRun {
		session, err := sshClient.NewSession()
		if err != nil {
			logrus.Fatalf("Failed to create session: %v", err.Error())
		}
		defer session.Close()
		result := RunSSHCommand(command, co.Host, session)
		CommandResults = append(CommandResults, result)
	}
	return CommandResults
}
