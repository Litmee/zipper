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

// logChan receive Log Pipeline Parameters
var logChan chan string

// OutLog output log to console
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

// logServer log server
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
