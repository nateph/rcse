package cliconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
)

// RunSSHCommand runs the passed command and records its information in the
// CommandResult struct
func RunSSHCommand(command string, host string, session *ssh.Session) (CommandResult, error) {
	var stdoutBuffer bytes.Buffer
	session.Stdout = &stdoutBuffer

	var stderrBuffer bytes.Buffer
	session.Stderr = &stderrBuffer

	sessionErr := session.Start(command)
	if sessionErr != nil {
		return CommandResult{}, fmt.Errorf("Failed on %s, %s\n", host, sessionErr)
	}
	err := session.Wait()
	if err != nil {
		return CommandResult{}, fmt.Errorf("Failed on %s, %s\n", host, sessionErr)
	}

	result := CommandResult{
		CommandRan: command,
		Host:       host,
		Stderr:     stderrBuffer.Bytes(),
		Stdout:     stdoutBuffer.Bytes(),
	}
	return result, nil
}

// CheckAndConsumePassword will prompt the user for a password, read it from STDIN,
// and return it.
func CheckAndConsumePassword(username string, password string) string {
	fmt.Printf("Enter password for user '%s': ", username)
	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logrus.Fatalf("Couldn't read password, error was: %v", err)
	}
	return string(pw)
}

// getKeyFile is a helper function from EstablishSSHConnection that reads in
// and parses the user's ssh id_rsa
func getKeyFile(currentUser *user.User) (key ssh.Signer, err error) {
	IDRsaFile := currentUser.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(IDRsaFile)
	if err != nil {
		logrus.Fatalf("unable to read file %s", IDRsaFile)
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		logrus.Fatalf("unable to parse private key: %v", err)
	}

	return key, err
}

// EstablishSSHConnection returns an ssh session from your id_rsa or username/password
func EstablishSSHConnection(username string, password string, host string, ignoreHostKeyCheck bool) (*ssh.Client, error) {
	var sshConfig *ssh.ClientConfig

	if username != "" {

		var hostKeyCallback ssh.HostKeyCallback

		if !ignoreHostKeyCheck {
			currentUser, _ := user.Current()
			knownHostsCallback, err := knownhosts.New(currentUser.HomeDir + "/.ssh/known_hosts")
			if err != nil {
				logrus.Fatal(err)
			}
			hostKeyCallback = knownHostsCallback
		} else {
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		sshConfig = &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: hostKeyCallback,
		}
	} else {
		currentUser, _ := user.Current()

		var hostKeyCallback ssh.HostKeyCallback

		if !ignoreHostKeyCheck {
			knownHostsCallback, err := knownhosts.New(currentUser.HomeDir + "/.ssh/known_hosts")
			if err != nil {
				logrus.Fatal(err)
			}
			hostKeyCallback = knownHostsCallback
		} else {
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		key, err := getKeyFile(currentUser)
		if err != nil {
			panic(err)
		}

		sshConfig = &ssh.ClientConfig{
			User: currentUser.Username,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
			HostKeyCallback: hostKeyCallback,
		}
	}
	client, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to host: %s", host)
	}

	return client, nil
}
