package connect

import (
	"zipper/message"
)

type ZipperRouter interface {
	Handler(ZipperRequest) message.ZipperMessage
}

type ZRouter struct{}

func (zr *ZRouter) Handler(req ZipperRequest) message.ZipperMessage {
	return nil
}
