package api

import "coco/util/errors"
//向JMS注册信息
func (s *Server) Register() (ResponsePass, errors.Error) {
	data := s.CreateQueryData()
	data["ip"] = s.Ip
	data["ws_port"] = s.WsPort
	data["ssh_port"] = s.SshPort
	var rd ResponsePass
	err := s.Query(s.Action.Register, data, &rd)
	return rd, err
}
