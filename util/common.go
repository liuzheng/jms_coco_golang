package util

import (
	"flag"
	"coco/util/log"
)

func Initial() {
	flag.Parse()
	log.Initial()
}
