package api

import "encoding/json"

//获取可用服务器列表
func (s *Server) GetList() ([]Machine, error) {
	data := s.CreateQueryData()
	res, _ := s.Query(s.Action.GetMachineList, data)
	var rd []Machine
	err := json.Unmarshal(res, &rd)
	return rd, err
}

//获取服务器登陆凭证
func (s *Server) GetLoginCredit(serverId, userId int) (LoginCredit, error) {
	data := s.CreateQueryData()
	res, _ := s.Query(s.Action.GetLoginCredit, data)
	var rd LoginCredit
	err := json.Unmarshal(res, &rd)
	return rd, err
}
