package command

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	terminal "golang.org/x/term"
)

// RunSSHCommand runs the command and records its information in a CommandResult
func RunSSHCommand(command string, host string, session *ssh.Session) (Result, error) {
	var stdoutBuffer bytes.Buffer
	session.Stdout = &stdoutBuffer

	var stderrBuffer bytes.Buffer
	session.Stderr = &stderrBuffer

	err := session.Run(command)
	if err != nil {
		return Result{}, fmt.Errorf("failed %s, %s %s", host, err, stderrBuffer.String())
	}

	result := Result{
		CommandRan: command,
		Host:       host,
		Stdout:     stdoutBuffer.String(),
	}
	return result, nil
}

// generateScriptName randomly generates file name, returned in an absolute path
func generateScriptName() string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	fileName := fmt.Sprintf("/tmp/rcse-%s.sh", strconv.Itoa(r1.Int()))

	return fileName
}

// TransferScript copies a script to a remote host and sets its mode to 0777
func TransferScript(scriptPath string, client *sftp.Client) (string, error) {
	// Generate a filename that likely doesn't exist
	dstFilePath := generateScriptName()
	dstFile, err := client.Create(dstFilePath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(scriptPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return "", err
	}

	err = client.Chmod(dstFilePath, 0777)
	if err != nil {
		return "", err
	}

	return dstFilePath, nil
}

// ConsumePassword will prompt the user for a password, reads it from STDIN, and returns it.
func ConsumePassword(username string, password string) (string, error) {
	fmt.Printf("Enter password for user (%s): ", username)

	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password from input: %v", err)
	}
	return string(pw), nil
}

// getKeyFile is a helper function from EstablishSSHConnection that reads in
// and parses the user's ssh id_rsa
func getKeyFile(currentUser *user.User, privateKeyPath string) (key ssh.Signer, err error) {
	var IDFile string
	if privateKeyPath == "" {
		IDFile = currentUser.HomeDir + "/.ssh/id_ed25519"
	} else {
		IDFile = privateKeyPath
	}
	buf, err := os.ReadFile(IDFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s", IDFile)
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return key, err
}

// EstablishSSHConnection returns an ssh client config from an id_rsa or username/password
func EstablishSSHConnection(username string, password string, host string, ignoreHostKeyCheck bool, privateKey string) (*ssh.Client, error) {
	var sshConfig *ssh.ClientConfig

	if username != "" {
		var hostKeyCallback ssh.HostKeyCallback

		if !ignoreHostKeyCheck {
			currentUser, _ := user.Current()
			knownHostsCallback, err := knownhosts.New(currentUser.HomeDir + "/.ssh/known_hosts")
			if err != nil {
				return nil, err
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
				return nil, err
			}
			hostKeyCallback = knownHostsCallback
		} else {
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		}

		key, err := getKeyFile(currentUser, privateKey)
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
		return nil, fmt.Errorf("failed to connect to host: %s\nerror: %s", host, err)
	}

	return client, nil
}

// EstablishSFTPConnection returns an sftp client for use
func EstablishSFTPConnection(sshClient *ssh.Client) (*sftp.Client, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}
