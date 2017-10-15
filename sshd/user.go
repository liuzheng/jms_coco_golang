package sshd

import (
	"golang.org/x/crypto/ssh"
)

// User describes an authenticable user.
type User struct {
	// The public key of the user.
	PublicKey ssh.PublicKey

	AuthKeys  string

	// The name the user will be referred to as. *NOT* the username used when
	// starting the session.
	Name      string
}