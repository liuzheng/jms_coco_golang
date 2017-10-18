package client

import (
	"golang.org/x/crypto/ssh"
	"coco/util/errors"
)

type Session struct {
	Session *ssh.Session
	window  windowDimensionChangeMsg
	Client  *Client
}

type windowDimensionChangeMsg struct {
	Columns uint32
	Rows    uint32
	Width   uint32
	Height  uint32
}

// 新建会话
func (c *Client) NewSession() (session *Session, erro errors.Error) {
	var err error
	session.Client = c
	session.Session, err = c.Client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
		erro = errors.New(err.Error(), 400)
	}
	c.Sessions = append(c.Sessions, session)
	return
}

// 调整窗口大小
func (s *Session) Resize(h, w int) (erro errors.Error) {
	s.window = windowDimensionChangeMsg{
		Columns: uint32(w),
		Rows:    uint32(h),
		Width:   uint32(w * 8),
		Height:  uint32(h * 8),
	}
	ok, err := s.Session.SendRequest("window-change", true, ssh.Marshal(&s.window))
	if err == nil && !ok {
		erro = errors.New("ssh: window-change failed", 400)
	}
	return
}
func (s *Session) Close() {
	s.Session.Close()
	remove(s.Client.Sessions, s)
}

// 从列表移除
func remove(s []*Session, r *Session) []*Session {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i + 1:]...)
		}
	}
	return s
}
