package cliconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
)

// RunSSHCommand runs the passed command and records its information in the
// CommandResult struct
func RunSSHCommand(command string, host string, session *ssh.Session) CommandResult {
	var stdoutBuffer bytes.Buffer
	session.Stdout = &stdoutBuffer

	var stderrBuffer bytes.Buffer
	session.Stderr = &stderrBuffer

	sessionErr := session.Run(command)

	if sessionErr != nil {
		logrus.Errorf("Failed on %s, %s\n", host, sessionErr)
	}
	result := CommandResult{
		CommandRan: command,
		Host:       host,
		Stderr:     stderrBuffer.Bytes(),
		Stdout:     stdoutBuffer.Bytes(),
	}
	return result
}

// CheckAndConsumePassword will prompt the user for a password, read it from STDIN,
// and set that password in viper.
func CheckAndConsumePassword() {
	fmt.Printf("Enter password for user '%s': ", viper.GetViper().GetString("user"))
	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logrus.Fatalf("Couldn't read password, error was: %v", err)
	}
	userPassword := string(pw)
	fmt.Println(userPassword)
	viper.Set("password", userPassword)
}

// getKeyFile is a helper function from EstablishSSHConnection that reads in
// and parses the users ssh id_rsa
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

// EstablishSSHConnection is meant to return an ssh session from your id_rsa
func EstablishSSHConnection(host string, ignoreHostKeyCheck bool) *ssh.Client {
	var sshConfig *ssh.ClientConfig

	if viper.IsSet("user") {
		sshConfig = &ssh.ClientConfig{
			User: viper.GetString("user"),
			Auth: []ssh.AuthMethod{
				ssh.Password(viper.GetString("password")),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
		logrus.Fatalf("Failed to dial: %v", err.Error())
	}

	return client
}
