package sshd

import (
	"golang.org/x/crypto/ssh"
	"coco/api"
	"coco/util/log"
)

// Session describes the current user session.
type Session struct {
	// Conn is the ssh.ServerConn associated with the connection.
	Conn      *ssh.ServerConn

	// User is the current user, or nil if unknown.
	User      *User

	// PublicKey is the public key used in this session.
	PublicKey ssh.PublicKey

	Machines  []api.Machine
}

// 关闭会话
func (s *Session)Close() {
	log.Info("Session", "Session Close")
	s.Conn.Close()
}