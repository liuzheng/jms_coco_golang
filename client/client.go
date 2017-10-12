package client

import (
	"golang.org/x/crypto/ssh"
	"net"
	"coco/api"
	"fmt"
	"coco/util/log"
)

type Client struct {
	Host   api.Machine
	Credit api.LoginCredit
	Signer ssh.Signer
	Config *ssh.ClientConfig

	Client   *ssh.Client
	Sessions []*Session

	SSHForward string // TODO: ssh forward internal
	Proxy      string // TODO: ssh use proxy
}

// 建立连接到后端服务器
func New(host api.Machine, credit api.LoginCredit) (client Client, err error) {
	client = Client{
		Host:   host,
		Credit: credit,
	}
	log.Debug("client", "New credit.PrivateKey : %v", credit.PrivateKey)
	client.Signer, err = ssh.ParsePrivateKey([]byte(credit.PrivateKey))
	if err != nil {
		log.Error("client", "unable to parse private key: %v", err)
	}
	client.Config = &ssh.ClientConfig{
		User: credit.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(client.Signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Ip, host.Port), client.Config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	return
}

// 关闭
func (c *Client) Close() {
	for _, session := range c.Sessions {
		session.Close()
	}
	c.Client.Close()
}
