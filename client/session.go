package client

import (
	"golang.org/x/crypto/ssh"
	"errors"
)

type Session struct {
	session *ssh.Session
	window  windowDimensionChangeMsg
	Client  *Client
}

type windowDimensionChangeMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}

func (c *Client) NewSession() (session *Session, err error) {
	session.Client = c
	session.session, err = c.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	c.Sessions = append(c.Sessions, session)
	return
}
func (s *Session) Resize(h, w int) error {
	s.window = windowDimensionChangeMsg{
		Columns: uint32(w),
		Rows:    uint32(h),
		Width:   uint32(w * 8),
		Height:  uint32(h * 8),
	}
	ok, err := s.session.SendRequest("window-change", true, ssh.Marshal(&s.window))
	if err == nil && !ok {
		err = errors.New("ssh: window-change failed")
	}
	return err
}
func (s *Session) Close() {
	s.session.Close()
	remove(s.Client.Sessions, s)
}

func remove(s []*Session, r *Session) []*Session {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
