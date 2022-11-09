package connect

import (
	"github.com/Litmee/zipper/message"
)

type ZipperRequest interface {
	// GetMsgId get message id
	GetMsgId() uint16
	// getConnect get link parameters
	getConnect() ZipperConnect
	// GetMsgBody get message
	GetMsgBody() []byte
}

type zRequest struct {
	conn ZipperConnect
	msg  message.ZipperMessage
}

func NewZRequest(c ZipperConnect, m message.ZipperMessage) ZipperRequest {
	return &zRequest{
		c,
		m,
	}
}

// GetMsgId get message id
func (zr *zRequest) GetMsgId() uint16 {
	return *zr.msg.GetMsgId()
}

// GetConnect get link parameters
func (zr *zRequest) getConnect() ZipperConnect {
	return zr.conn
}

// GetMsgBody get message
func (zr *zRequest) GetMsgBody() []byte {
	return zr.msg.GetBody()
}
