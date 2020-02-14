package cliconfig

import (
	"io/ioutil"
	"log"
	"os/user"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func getKeyFile(currentUser *user.User) (key ssh.Signer, err error) {
	IDRsaFile := currentUser.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(IDRsaFile)
	if err != nil {
		log.Fatalf("unable to read file %s", IDRsaFile)
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	return key, err
}

// EstablishSSHConnection is meant to return an ssh session from your id_rsa
func EstablishSSHConnection(host string, ignoreHostKeyCheck bool) *ssh.Session {
	currentUser, _ := user.Current()

	var hostKeyCallback ssh.HostKeyCallback

	if !ignoreHostKeyCheck {
		knownHostsCallback, err := knownhosts.New(currentUser.HomeDir + "/.ssh/known_hosts")
		if err != nil {
			log.Fatal(err)
		}
		hostKeyCallback = knownHostsCallback
	} else {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	key, err := getKeyFile(currentUser)
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: currentUser.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err.Error())
	}

	return session
}
