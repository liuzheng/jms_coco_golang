package api

import (
	"encoding/json"
)

// 用户权限认证
func (s *Server) Auth() (u UserAuth, err error) {
	u.Username = "root"
	u.Action = s.Action
	u.server = s
	return
}

// Query 方法
func (u *UserAuth) Query(action string, data map[string]interface{}) ([]byte, error) {
	return u.server.Query(u.Action.GetUserPubKey, data)
}

// 用户名获取用户的 pubkey
func (u *UserAuth) GetUserPubKey() (UserPubKey, error) {
	data := u.server.CreateQueryData()
	data["username"] = u.Username
	res, _ := u.Query(u.Action.GetUserPubKey, data)
	err := json.Unmarshal(res, &u.UserPubKey)
	return u.UserPubKey, err
}

// 获取登陆用户TOKEN
func (u *UserAuth) GetLoginToken() (UserToken, error) {
	data := u.server.CreateQueryData()
	data["username"] = u.Username
	data["ticket"] = u.UserPubKey.Ticket
	res, _ := u.Query(u.Action.GetUserToken, data)
	err := json.Unmarshal(res, &u.UserToken)
	return u.UserToken, err
}

// 检查用户能否开启监控SHELL
func (u *UserAuth) CheckMonitorToken(sessionId int) (ResponsePass, error) {
	data := u.server.CreateQueryData()
	data["session_id"] = sessionId
	res, _ := u.Query(u.Action.CheckMonitorToken, data)
	var rd ResponsePass
	err := json.Unmarshal(res, &rd)
	return rd, err
}
