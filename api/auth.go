package api

import (
	"coco/util/errors"
	"coco/util/log"
)

// 用户名获取用户的 pubkey
func (s *Server) GetUserPubKey(username string) (UserPubKey, errors.Error) {
	data := s.CreateQueryData()
	data["username"] = username
	var rd UserPubKey
	err := s.Query(s.Action.GetUserPubKey, data, &rd)
	return rd, err
}

// 获取登陆用户TOKEN
func (s *Server) GetLoginToken(username, ticket string) (UserToken, errors.Error) {
	data := s.CreateQueryData()
	data["username"] = username
	data["ticket"] = ticket
	var rd UserToken
	err := s.Query(s.Action.GetUserToken, data, &rd)
	s.Token = rd
	return rd, err
}

// 检查用户能否开启监控SHELL
func (s *Server) CheckMonitorToken(sessionId int) (ResponsePass, errors.Error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.CheckMonitorToken, data, &rd)
	return rd, err
}

func (s *Server) Login(username string) bool {
	up, err := s.GetUserPubKey(username)
	if log.HandleErr("Login", err, "login issue") {
		return false
	}
	_, err = s.GetLoginToken(username, up.Ticket)
	if log.HandleErr("Login", err, "login issue") {
		return false
	}
	return true

}
