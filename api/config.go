package api

import "encoding/json"

//向JMS注册信息
func (s *Server) Register() (ResponsePass, error) {
	data := s.CreateQueryData()
	data["ip"] = s.Ip
	data["ws_port"] = s.WsPort
	data["ssh_port"] = s.SshPort
	res, _ := s.Query(s.Action.Register, data)
	var rd ResponsePass
	err := json.Unmarshal(res, &rd)
	return rd, err
}
