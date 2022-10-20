package main

import (
	"zipper/server"
)

func main() {
	zServer := server.NewZServer()
	zServer.Run()
}
