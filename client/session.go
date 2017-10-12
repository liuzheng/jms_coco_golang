package client

import (
	"golang.org/x/crypto/ssh"
	"net"
	"coco/api"
	"fmt"
)

//建立连接到后端服务器的Session
func New(host api.Machine, credit api.LoginCredit) (session *ssh.Session) {
	signer, err := ssh.ParsePrivateKey([]byte(credit.PrivateKey))
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: credit.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Ip, host.Port), config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err = client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	return
}
