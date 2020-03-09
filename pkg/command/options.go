package command

import "fmt"

// Options contains information on the the command to be ran
type Options struct {
	// Command that will be ran
	CommandToRun string
	// Which host it will be ran on
	Host string
	// Whether or not to verify host keys
	IgnoreHostkeyCheck bool
	// User to execute as
	User string
	// Password for User
	Password string
}

// RunCommand is a wrapper around establishing the ssh connection and then
// calling the RunSSHCommand
func (o *Options) RunCommand() (Result, error) {
	sshClient, err := EstablishSSHConnection(o.User, o.Password, o.Host, o.IgnoreHostkeyCheck)
	if err != nil {
		return Result{}, err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return Result{}, fmt.Errorf("Failed to create session: %v", err.Error())
	}
	defer session.Close()

	// If 'ping' is being ran, no command is ran and on successfull SSH connection,
	// 'pong' is returned to signify success.
	if o.CommandToRun == "" {
		return Result{
			Host:       o.Host,
			Stdout:     "pong",
			CommandRan: "",
		}, nil
	}

	result, err := RunSSHCommand(o.CommandToRun, o.Host, session)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}
