package api

import (
	"encoding/json"
)

func (s *Server) Auth() (u UserAuth, err error) {
	u.Username = "root"
	u.Action = s.Action
	u.Server = s
	return
}

//根据用户名获取pubkey
func (u *UserAuth) GetUserPubKey() (UserPubKey, error) {
	data := u.Server.CreateQueryData()
	data["username"] = u.Username
	res, _ := u.Server.Query(u.Action.GetUserPubKey, data)
	err := json.Unmarshal(res, &u.UserPubKey)
	return u.UserPubKey, err
}

//根据pubkey和username获取登陆TOKEN
func (u *UserAuth) GetLoginToken() (UserToken, error) {
	data := u.Server.CreateQueryData()
	data["username"] = u.Username
	data["ticket"] = u.UserPubKey.Ticket
	res, _ := u.Server.Query(u.Action.GetUserToken, data)
	err := json.Unmarshal(res, &u.UserToken)
	return u.UserToken, err
}

//检查用户能否开启监控SHELL
func (u *UserAuth) CheckMonitorToken(sessionId int) (ResponsePass, error) {
	data := u.Server.CreateQueryData()
	data["session_id"] = sessionId
	res, _ := u.Server.Query(u.Action.CheckMonitorToken, data)
	var rd ResponsePass
	err := json.Unmarshal(res, &rd)
	return rd, err
}
