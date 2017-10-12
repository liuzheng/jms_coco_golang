package log

import (
	"github.com/liuzheng712/golog"
)

func Initial() {
	golog.Initial()
}

var Debug = golog.Debug
var Info = golog.Info
var Notice = golog.Notice
var Warn = golog.Warn
var Error = golog.Error
var Critical = golog.Critical
var Panic = golog.Panic
var Fatal = golog.Fatal

func HandleErr(name string, err error) bool {
	if err != nil {
		golog.Error(name, "%v", err)
		return true
	}
	return false
}
