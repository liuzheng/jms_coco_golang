package api

// 获取可用服务器列表
func (s *Server) GetList() ([]Machine, error) {
	data := s.CreateQueryData()
	var rd []Machine
	err := s.Query(s.Action.GetMachineList, data, rd)
	return rd, err
}

// 获取服务器登陆凭证
func (s *Server) GetLoginCredit(serverId, userId int) (LoginCredit, error) {
	data := s.CreateQueryData()
	var rd LoginCredit
	err := s.Query(s.Action.GetLoginCredit, data, rd)
	return rd, err
}
