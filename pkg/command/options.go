package command

import (
	"fmt"
)

// Options contains information on the the command to be ran
type Options struct {
	// Command that will be ran, if provided
	Command string
	// if true, delete the script after execution
	Cleanup bool
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
	// A script to be ran, if provided
	Script string
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
		return Result{}, fmt.Errorf("failed to create session: %v", err.Error())
	}
	defer session.Close()

	result, err := RunSSHCommand(opts.Command, opts.Host, session)
	if err != nil {
		return Result{}, err
	}

	return result, nil
}

// RunScript is a wrapper around establishing an sftp connection, and running a script
func (opts *Options) RunScript() (Result, error) {
	conn, err := EstablishSSHConnection(opts.User, opts.Password, opts.Host, opts.IgnoreHostkeyCheck, opts.PrivateKey)
	if err != nil {
		return Result{}, err
	}
	defer conn.Close()

	// Setup an sftp connection with the ssh connection
	sftpClient, err := EstablishSFTPConnection(conn)
	if err != nil {
		return Result{}, err
	}
	defer sftpClient.Close()

	scriptPath, err := TransferScript(opts.Script, sftpClient)
	if err != nil {
		return Result{}, err
	}

	session, err := conn.NewSession()
	if err != nil {
		return Result{}, fmt.Errorf("failed to create session: %v", err.Error())
	}
	defer session.Close()

	result, err := RunSSHCommand(scriptPath, opts.Host, session)
	if err != nil {
		return Result{}, err
	}

	if opts.Cleanup {
		err := sftpClient.Remove(scriptPath)
		if err != nil {
			return Result{}, err
		}
	}

	return result, nil
}
