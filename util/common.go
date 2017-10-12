package util

import (
	"flag"
	"coco/util/log"
)

func Init() {
	flag.Parse()
	log.Initial()
}
