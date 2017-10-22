package log

import (
	"github.com/liuzheng712/golog"
	"strings"
	//"reflect"
	//"coco/util/errors"
)

func Initial() {
	golog.Initial()
}

type Password string

func (p Password) Redacted() interface{} {
	return strings.Repeat("*", len(p))
}

var Logs = golog.Logs

var Debug = golog.Debug
var Info = golog.Info
var Notice = golog.Notice
var Warn = golog.Warn
var Error = golog.Error
var Critical = golog.Critical
var Panic = golog.Panic
var Fatal = golog.Fatal

func HandleErr(name string, err ... interface{}) bool {
	if err != nil && err[0] != nil {
		if len(err) > 1 {
			// TODO: discuss if it need judge the err's type to print error's code <liuzheng712@gmail.com, xrain@simcu.com>
			//if reflect.TypeOf(err[1]).String() == "*error.UtilErr" {
			//	golog.Error(name, "%v, Code: %v", err[1], err[0].(errors.Error).GetCode())
			//} else {
			golog.Error(name, "%v", err[1])
			//}
		} else {
			golog.Error(name, "Error Unknow")
		}
		golog.Debug(name, "%v", err[0])
		return true
	}
	return false
}
