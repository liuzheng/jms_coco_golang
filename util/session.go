package util

import "fmt"

//获取录像文件路径
func GetReportFile(sessionId, seq int) (string, error) {
	return fmt.Sprintf("/tmp/jms/record/$d-$d.jmsr", sessionId, seq), nil
}

//删除录像文件
func DeleteReportFile(sessionId, seq int) (bool, error) {
	return true, nil
}
