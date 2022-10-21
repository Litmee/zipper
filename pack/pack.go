package pack

import (
	"bytes"
	"encoding/binary"
	"zipper/message"
)

// ZipperPack TCP data unpacking abstraction layer
type ZipperPack interface {
	// GetHeadLen get the length of the packet header information, fixed to 4
	GetHeadLen() uint16
	// Pack packing method
	Pack(m message.ZipperMessage) ([]byte, error)
	// UnPack unpacking method
	UnPack(d []byte) (message.ZipperMessage, error)
}

type ZPack struct{}

func NewZPack() ZipperPack {
	return &ZPack{}
}

func (zp *ZPack) GetHeadLen() uint16 {
	return 4
}

// Pack packing method
func (zp *ZPack) Pack(m message.ZipperMessage) (d []byte, err error) {
	// create a buffer to hold byte data
	buf := bytes.NewBuffer([]byte{})
	// first write the length information of the message into buf
	if err = binary.Write(buf, binary.LittleEndian, m.GetMsgLen()); err != nil {
		return
	}
	// then write the data id to buf
	if err = binary.Write(buf, binary.LittleEndian, m.GetMsgId()); err != nil {
		return
	}
	// finally write the data body to buf
	if err = binary.Write(buf, binary.LittleEndian, m.GetBody()); err != nil {
		return
	}
	d = buf.Bytes()
	return
}

// UnPack unpacking method
func (zp *ZPack) UnPack(d []byte) (m message.ZipperMessage, err error) {
	// create an io.Reader that reads binary data
	reader := bytes.NewReader(d)
	msg := message.NewZMessage(0, nil)
	// first read the data length
	if err = binary.Read(reader, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return
	}
	// then read the message id
	if err = binary.Read(reader, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return
	}
	m = msg
	return
}
