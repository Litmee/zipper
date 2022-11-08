package connect

import (
	"context"
	"github.com/Litmee/zipper/common"
	"github.com/Litmee/zipper/logger"
	"github.com/Litmee/zipper/message"
	"github.com/Litmee/zipper/pack"
	"io"
	"net"
	"sync"
)

// ZipperConnect TCP connection layer abstraction
type ZipperConnect interface {
	// Start run link
	Start(ctx context.Context)
	// zRead read business part processing function
	zRead(ctx context.Context)
	// ZWrite write the business part processing function
	zWrite(ctx context.Context)
	// send write-back message
	send(m message.ZipperMessage)
	// stop read and write
	stop()
}

type zConnect struct {
	// TCP connection
	conn *net.TCPConn
	// queue worker pool parameters
	pool ZipperPool
	// write-back message receive pipeline
	msgChan chan message.ZipperMessage
	// whether the link has been closed
	isClosed bool
	// receive and send exit signal
	exitChan chan bool
	// lock
	rwl sync.RWMutex
}

func NewZConnect(conn *net.TCPConn, pool ZipperPool) ZipperConnect {
	return &zConnect{
		conn,
		pool,
		make(chan message.ZipperMessage, common.GlobalConfig.QueueSize),
		false,
		make(chan bool, 1),
		sync.RWMutex{},
	}
}

// Start run link
func (zc *zConnect) Start(ctx context.Context) {
	// open the read method
	go zc.zRead(ctx)
	// open the write method
	go zc.zWrite(ctx)
}

// zRead read business part processing function
func (zc *zConnect) zRead(ctx context.Context) {
	defer common.Cut()
	defer zc.stop()
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
			continue
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
func (zc *zConnect) zWrite(ctx context.Context) {
	zPack := pack.NewZPack()
	for {
		select {
		case msg := <-zc.msgChan:
			b, err := zPack.Pack(msg)
			if err != nil {
				logger.OutLog("The write-back data packet is abnormal", logger.WARN)
				break
			}
			_, err = zc.conn.Write(b)
			if err != nil {
				logger.OutLog("TCP write-back data err", logger.ERROR)
				return
			}
		case <-zc.exitChan:
			return
		}
	}
}

// send write-back message
func (zc *zConnect) send(m message.ZipperMessage) {
	zc.rwl.RLock()
	defer zc.rwl.RUnlock()
	if zc.isClosed {
		return
	}
	zc.msgChan <- m
}

// stop read and write
func (zc *zConnect) stop() {
	zc.rwl.Lock()
	defer zc.rwl.Unlock()
	if zc.isClosed {
		return
	}
	// send exit signal
	zc.exitChan <- true
	zc.isClosed = true
	close(zc.msgChan)
	close(zc.exitChan)
}
