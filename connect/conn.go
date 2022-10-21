package connect

import (
	"context"
	"io"
	"net"
	"zipper/logger"
	"zipper/message"
	"zipper/pack"
)

// ZipperConnect TCP connection layer abstraction
type ZipperConnect interface {
	// Start run link
	Start(ctx context.Context)
	// ZRead 读业务
	zRead(ctx context.Context, f context.CancelFunc)
	// ZWrite 写业务
	zWrite(ctx context.Context, f context.CancelFunc)
	// Send write-back message
	Send(m message.ZipperMessage)
}

type zConnect struct {
	// TCP connection
	conn *net.TCPConn
	// queue worker pool parameters
	pool ZipperPool
	// write-back message receive pipeline
	msgChan chan message.ZipperMessage
}

func NewZConnect(conn *net.TCPConn, pool ZipperPool) ZipperConnect {
	return &zConnect{conn, pool, make(chan message.ZipperMessage)}
}

func (zc *zConnect) Start(ctx context.Context) {
	// iterates out new context control parameters
	newCtx, cancelFunc := context.WithCancel(ctx)
	// open the read method
	go zc.zRead(newCtx, cancelFunc)
	// open the write method
	go zc.zWrite(newCtx, cancelFunc)
}

func (zc *zConnect) zRead(ctx context.Context, f context.CancelFunc) {
	defer f()
	// create Packet Unpacking Parameters
	zPack := pack.NewZPack()
	for {
		// create a binary byte container that receives the header of the message body
		headData := make([]byte, zPack.GetHeadLen())
		if _, err := io.ReadFull(zc.conn, headData); err != nil {
			logger.OutLog("Read msg head err: "+err.Error(), logger.ERROR)
			break
		}
		// unpacking
		msg, err := zPack.UnPack(headData)
		if err != nil {
			logger.OutLog("UnPack msg head err: "+err.Error(), logger.ERROR)
			break
		}
		// get message body
		var body []byte
		n := *msg.GetMsgLen()
		if n > 0 {
			body = make([]byte, n)
			if _, err = io.ReadFull(zc.conn, body); err != nil {
				logger.OutLog("Read msg body err: "+err.Error(), logger.ERROR)
				break
			}
		}
		// splice body
		msg.SetBody(body)
		// send to queue
		zc.pool.AddQueue(zc, msg)
	}

}

func (zc *zConnect) zWrite(ctx context.Context, f context.CancelFunc) {
	defer f()
	zPack := pack.NewZPack()
	for {
		msg := <-zc.msgChan
		b, err := zPack.Pack(msg)
		if err != nil {
			logger.OutLog("The write-back data packet is abnormal", logger.WARN)
			continue
		}
		_, err = zc.conn.Write(b)
		if err != nil {
			logger.OutLog("TCP write-back data err", logger.ERROR)
			break
		}
	}
}

func (zc *zConnect) Send(m message.ZipperMessage) {
	zc.msgChan <- m
}
