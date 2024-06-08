package utils

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConnector struct {
	client  *ssh.Client
	timeout time.Duration
}

func NewSSHConnector() *SSHConnector {
	return &SSHConnector{
		timeout: 30 * time.Second,
	}
}

func (s *SSHConnector) Connect(serverIP string, username string, sshKeyPath string) (*ssh.Client, error) {
	if serverIP == "" || username == "" || sshKeyPath == "" {
		return nil, fmt.Errorf("serverIP, username, or sshKeyPath is empty")
	}

	key, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         s.timeout,
	}

	client, err := ssh.Dial("tcp", serverIP+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	s.client = client
	return client, nil
}

func (s *SSHConnector) Disconnect() {
	if s.client != nil {
		s.client.Close()
	}
}
