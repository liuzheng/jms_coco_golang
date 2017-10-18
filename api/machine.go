package api

// 获取可用服务器列表
func (s *Server) GetList(keyword string, groupId int) ([]Machine, RespError) {
	data := s.CreateQueryData()
	data["keyword"] = keyword
	data["group_id"] = groupId
	var rd []Machine
	err := s.Query(s.Action.GetMachineList, data, &rd)
	return rd, err
}

// 获取服务器登陆凭证
func (s *Server) GetLoginCredit(serverId, userId int) (LoginCredit, RespError) {
	data := s.CreateQueryData()
	var rd LoginCredit
	err := s.Query(s.Action.GetLoginCredit, data, &rd)
	return rd, err
}

//获取服务器组
func (s *Server) GetGroupList() ([]MachineGroup, RespError) {
	data := s.CreateQueryData()
	var rd []MachineGroup
	err := s.Query(s.Action.GetMachineGroupList, data, &rd)
	return rd, err
}
