package api

//上传录像文件
func (s *Server) Upload(sessionId int, isEnd bool) (bool, RespError) {
	return true, RespError{ErrorCode: 200, ErrorMsg: ""}
}
