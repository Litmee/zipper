package message

// ZipperMessage message abstraction layer
type ZipperMessage interface {
	// GetMsgLen get message length
	GetMsgLen() *uint16
	// GetMsgId get message id
	GetMsgId() *uint16
	// GetBody get message body
	GetBody() []byte
	// SetMsgLen set message length
	SetMsgLen(l uint16)
	// SetMsgId set message id
	SetMsgId(id uint16)
	// SetBody set message body
	SetBody(b []byte)
}

type zMessage struct {
	// length of the message body
	l uint16
	// message id
	id uint16
	// message body
	body []byte
}

func NewZMessage(id uint16, body []byte) ZipperMessage {
	return &zMessage{l: uint16(len(body)), id: id, body: body}
}

// GetMsgLen get message length
func (zm *zMessage) GetMsgLen() *uint16 {
	return &zm.l
}

// GetMsgId get message id
func (zm *zMessage) GetMsgId() *uint16 {
	return &zm.id
}

// GetBody get message body
func (zm *zMessage) GetBody() []byte {
	return zm.body
}

// SetMsgLen set message length
func (zm *zMessage) SetMsgLen(l uint16) {
	zm.l = l
}

// SetMsgId set message id
func (zm *zMessage) SetMsgId(id uint16) {
	zm.id = id
}

// SetBody set message body
func (zm *zMessage) SetBody(b []byte) {
	zm.body = b
}
