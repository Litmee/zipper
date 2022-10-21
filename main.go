package main

import (
	"zipper/connect"
	"zipper/message"
	"zipper/server"
)

type MyRouter struct {
	connect.ZRouter
}

func (my *MyRouter) Handler(req connect.ZipperRequest) message.ZipperMessage {
	s := string(req.GetMsg().GetBody())
	s += " Hello World"
	zMessage := message.NewZMessage(1, []byte(s))
	return zMessage
}

func main() {
	zServer := server.NewZServer()
	zServer.AddRouter(1, &MyRouter{})
	zServer.Run()
}
