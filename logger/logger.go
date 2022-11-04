package logger

import (
	"context"
	"log"
)

// const log level
const (
	INFO = iota
	WARN
	ERROR
)

// logChan receive Log Pipeline Parameters
var logChan chan string

// OutLog output log to console
func OutLog(s string, t int) {
	// judge log type
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

// InitLogServer init logServer
func InitLogServer(ctx context.Context, s uint16) {
	logChan = make(chan string, s)
	// go start log server
	go logServer(ctx)
}
