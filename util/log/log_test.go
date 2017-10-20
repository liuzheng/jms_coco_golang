package log

import (
	"coco/util/errors"
	"testing"
	"flag"
)

func init() {
	flag.Parse()
	Initial()
}

func Test_HandleErr(t *testing.T) {
	err := errors.New("sdf", 2)
	Logs("", "INFO", "INFO")
	t.Log("Info Module")
	HandleErr("main", err, "test")
	t.Log("Debug Module")
	Logs("", "DEBUG", "INFO")
	HandleErr("main", err, "test")
}

func Test_HandleErr_one(t *testing.T) {
	err := errors.New("sdf", 2)
	Logs("", "INFO", "INFO")
	t.Log("Info Module")
	HandleErr("main", err)
	t.Log("Debug Module")
	Logs("", "DEBUG", "INFO")
	HandleErr("main", err)
}
