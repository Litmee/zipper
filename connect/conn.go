package connect

import (
	"context"
	"fmt"
	"io"
	"net"
	"zipper/common"
	"zipper/logger"
	"zipper/message"
	"zipper/pack"
)

// ZipperConnect TCP connection layer abstraction
type ZipperConnect interface {
	// Start run link
	Start(ctx context.Context)
	// zRead read business part processing function
	zRead(ctx context.Context, f context.CancelFunc)
	// ZWrite write the business part processing function
	zWrite(ctx context.Context, f context.CancelFunc)
	// send write-back message
	send(m message.ZipperMessage)
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

// Start run link
func (zc *zConnect) Start(ctx context.Context) {
	// iterates out new context control parameters
	newCtx, cancelFunc := context.WithCancel(ctx)
	// open the read method
	go zc.zRead(newCtx, cancelFunc)
	// open the write method
	go zc.zWrite(newCtx, cancelFunc)
}

// zRead read business part processing function
func (zc *zConnect) zRead(ctx context.Context, f context.CancelFunc) {
	defer f()
	defer common.Cut()
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
		zc.pool.addQueue(zc, msg)
	}

}

// zWrite write the business part processing function
func (zc *zConnect) zWrite(ctx context.Context, f context.CancelFunc) {
	defer f()
	defer fmt.Println("一切都结束了2")
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

// send write-back message
func (zc *zConnect) send(m message.ZipperMessage) {
	zc.msgChan <- m
}
