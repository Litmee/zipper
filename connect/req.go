package connect

import (
	"zipper/message"
)

type ZipperRequest interface {
	// GetMsgId get message id
	GetMsgId() uint16
	// GetConnect get link parameters
	GetConnect() ZipperConnect
	// GetMsg get message
	GetMsg() message.ZipperMessage
}

type zRequest struct {
	conn ZipperConnect
	msg  message.ZipperMessage
}

func NewZRequest(c ZipperConnect, m message.ZipperMessage) ZipperRequest {
	return &zRequest{c, m}
}

func (zr *zRequest) GetMsgId() uint16 {
	return *zr.msg.GetMsgId()
}

func (zr *zRequest) GetConnect() ZipperConnect {
	return zr.conn
}

func (zr *zRequest) GetMsg() message.ZipperMessage {
	return zr.msg
}
