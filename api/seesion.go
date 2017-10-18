package api


//上报开启的Session
func (s *Server) ReportSession(sessionId string, serverId, userId, seq int) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.ReportSession, data, &rd)
	return rd, err
}

func (s *Server) ReportSessionClose(sessionId string) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.ReportSessionClose, data, &rd)
	return rd, err
}
