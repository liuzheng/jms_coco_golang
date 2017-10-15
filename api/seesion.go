package api

//上报开启的Session
func (s *Server) ReportSession(sessionId, serverId, userId, seq int) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	var rd ResponsePass
	err := s.Query(s.Action.ReportSession, data, rd)
	return rd, err
}
