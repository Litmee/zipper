package logger

import (
	"context"
	"github.com/Litmee/zipper/common"
	"log"
)

const (
	INFO = iota
	WARN
	ERROR
)

var Logger = make(chan string, common.GlobalConfig.LogSize)

func OutLog(s string, t int) {
	switch t {
	case INFO:
		s = "INFO " + s
	case WARN:
		s = "WARN " + s
	case ERROR:
		s = "ERROR " + s
	}
	Logger <- s
}

func InitLogServer(ctx context.Context) {
	for {
		str := <-Logger
		log.Println(str)
	}
}
