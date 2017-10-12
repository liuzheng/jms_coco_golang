package log

import (
	"github.com/liuzheng712/golog"
)

func init() {
	golog.Initial()
}

var Initial = golog.Initial
var Debug = golog.Debug
var Info = golog.Info
var Notice = golog.Notice
var Warn = golog.Warn
var Error = golog.Error
var Critical = golog.Critical
var Panic = golog.Panic
var Fatal = golog.Fatal
