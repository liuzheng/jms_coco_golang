package sshd

import (
	"errors"
	"golang.org/x/crypto/ssh"
)

type publicKey struct {
	publicKey     []byte
	publicKeyType string
}

func (p *publicKey) Marshal() []byte {
	//b := make([]byte, len(p.publicKey))
	//copy(b, p.publicKey)
	return p.publicKey
}

func (p *publicKey) Type() string {
	return p.publicKeyType
}

func (p *publicKey) Verify([]byte, *ssh.Signature) error {
	return errors.New("verify not implemented")
}
