package api

import (
	"encoding/json"
)

//根据用户名获取pubkey
func (s *Server) GetUserPubKey(username string) (UserPubKey, error) {
	data := s.CreateQueryData()
	data["username"] = username
	res, _ := s.Query(s.Action.GetUserPubKey, data)
	var rd UserPubKey
	err := json.Unmarshal(res, &rd)
	return rd, err
}

//根据pubkey和username获取登陆TOKEN
func (s *Server) GetLoginToken(username string, ticket string) (UserToken, error) {
	data := s.CreateQueryData()
	data["username"] = username
	data["ticket"] = ticket
	res, _ := s.Query(s.Action.GetUserToken, data)
	var rd UserToken
	err := json.Unmarshal(res, &rd)
	return rd, err
}

//检查用户能否开启监控SHELL
func (s *Server) CheckMonitorToken(sessionId int) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	res, _ := s.Query(s.Action.CheckMonitorToken, data)
	var rd ResponsePass
	err := json.Unmarshal(res, &rd)
	return rd, err
}
