package api

import "coco/util/errors"

//上传录像文件
func (s *Server) Upload(sessionId int, isEnd bool) (bool, errors.Error) {
	return true, errors.New("", 200)
}
