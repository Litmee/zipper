package logger

import (
	"context"
	"log"
)

const (
	INFO = iota
	WARN
	ERROR
)

var logChan chan string

func OutLog(s string, t int) {
	switch t {
	case INFO:
		s = "INFO " + s
	case WARN:
		s = "WARN " + s
	case ERROR:
		s = "ERROR " + s
	}
	logChan <- s
}

func logServer(ctx context.Context) {
	for {
		str := <-logChan
		log.Println(str)
	}
}

func InitLogServer(ctx context.Context, s uint16) {
	logChan = make(chan string, s)
	go logServer(ctx)
}
