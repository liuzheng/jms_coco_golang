package api

// 用户名获取用户的 pubkey
func (s *Server) GetUserPubKey(username string) (UserPubKey, error) {
	data := s.CreateQueryData()
	data["username"] = username
	var rd UserPubKey
	err := s.Query(s.Action.GetUserPubKey, data, &rd)
	return rd, err
}

// 获取登陆用户TOKEN
func (s *Server) GetLoginToken(username, ticket string) (UserToken, error) {
	data := s.CreateQueryData()
	data["username"] = username
	data["ticket"] = ticket
	var rd UserToken
	err := s.Query(s.Action.GetUserToken, data, &rd)
	s.Token = rd
	return rd, err
}

// 检查用户能否开启监控SHELL
func (s *Server) CheckMonitorToken(sessionId int) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.CheckMonitorToken, data, &rd)
	return rd, err
}
