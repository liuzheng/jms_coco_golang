package log

import (
	"github.com/liuzheng712/golog"
)

func init() {
	golog.Initial()
}

func Initial() {
	golog.Initial()
}
func Debug(selector string, format string, v ...interface{}) {
	golog.Debug(selector, format, v)
}
func Info(selector string, format string, v ...interface{}) {
	golog.Info(selector, format, v)
}
func Notice(selector string, format string, v ...interface{}) {
	golog.Notice(selector, format, v)
}
func Warn(selector string, format string, v ...interface{}) {
	golog.Warn(selector, format, v)
}
func Error(selector string, format string, v ...interface{}) {
	golog.Error(selector, format, v)
}
func Critical(selector string, format string, v ...interface{}) {
	golog.Critical(selector, format, v)
}
func Panic(selector string, format string, v ...interface{}) {
	golog.Panic(selector, format, v)
}
func Fatal(selector string, format string, v ...interface{}) {
	golog.Fatal(selector, format, v)
}
