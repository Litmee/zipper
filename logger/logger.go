package logger

import (
	"context"
	"log"
)

var Logger = make(chan string, 200)

func OutLog(s string) {
	Logger <- s
}

func InitLogServer(ctx context.Context) {
	for {
		str := <-Logger
		log.Println(str)
	}
}
