package api

import "encoding/json"

//上报开启的Session
func (s *Server) ReportSession(sessionId, serverId, userId, seq int) (ResponsePass, error) {
	data := s.CreateQueryData()
	data["session_id"] = sessionId
	res, _ := s.Query(s.Action.ReportSession, data)
	var rd ResponsePass
	err := json.Unmarshal(res, &rd)
	return rd, err
}
