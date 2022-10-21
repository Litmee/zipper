package connect

import (
	"zipper/message"
)

type ZipperRouter interface {
	// Handler business processing hook function
	Handler(ZipperRequest) message.ZipperMessage
}

type ZRouter struct{}

// Handler business processing hook function
func (zr *ZRouter) Handler(req ZipperRequest) message.ZipperMessage {
	return nil
}
