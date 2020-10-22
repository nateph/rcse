package command

import (
	"fmt"
)

// Options contains information on the the command to be ran
type Options struct {
	// Command that will be ran
	Command string
	// Whether or not to verify host keys
	IgnoreHostkeyCheck bool
	// Host to execute on
	Host string
	// User to execute as
	User string
	// Password for User
	Password string
	// A custom private key to use
	PrivateKey string
}

// RunCommand is a wrapper around establishing the ssh connection and then
// calling the RunSSHCommand
func (opts *Options) RunCommand() (Result, error) {
	sshClient, err := EstablishSSHConnection(opts.User, opts.Password, opts.Host, opts.IgnoreHostkeyCheck, opts.PrivateKey)
	if err != nil {
		return Result{}, err
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return Result{}, fmt.Errorf("Failed to create session: %v", err.Error())
	}
	defer session.Close()

	result, err := RunSSHCommand(opts.Command, opts.Host, session)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}
