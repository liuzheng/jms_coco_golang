package api

import "coco/util/errors"

//上报开启的Session
func (s *Server) ReportSession(sessionId string, serverId, userId, seq int) (ResponsePass, errors.Error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.ReportSession, data, &rd)
	return rd, err
}

func (s *Server) ReportSessionClose(sessionId string) (ResponsePass, errors.Error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.ReportSessionClose, data, &rd)
	return rd, err
}
