package connect

import (
	"github.com/Litmee/zipper/message"
)

// ZipperRouter route layer abstract interface
type ZipperRouter interface {
	// Handler business processing hook function
	Handler(ZipperRequest) message.ZipperMessage
}

type ZRouter struct{}

// Handler business processing hook function
func (zr *ZRouter) Handler(req ZipperRequest) message.ZipperMessage {
	return nil
}
