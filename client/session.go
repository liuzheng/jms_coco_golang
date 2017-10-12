package client

import (
	"golang.org/x/crypto/ssh"
	"errors"
)

type Session struct {
	session *ssh.Session
	window  windowDimensionChangeMsg
}

type windowDimensionChangeMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}

func (c *Client) NewSession() (session *Session, err error) {
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
}
