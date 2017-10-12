package util

import (
	"flag"
	"github.com/liuzheng712/golog"
)

func Init() {
	flag.Parse()
	golog.Initial()
}
