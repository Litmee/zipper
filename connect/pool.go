package connect

import (
	"context"
	"strconv"
	"time"
	"zipper/common"
	"zipper/logger"
	"zipper/message"
)

// ZipperPool work queue task pool abstraction layer
type ZipperPool interface {
	// InitPool initialize the worker pool
	InitPool(ctx context.Context)
	// AddRouter add message id routing mapping
	AddRouter(id uint16, rt ZipperRouter)
	// RunQueue open queue job
	RunQueue(ctx context.Context, id int, queue chan ZipperRequest)
	// ConsumeMessage Consuming messages
	ConsumeMessage(req ZipperRequest)
	// AddQueue add messages to the queue by pseudo load balancing
	AddQueue(c ZipperConnect, m message.ZipperMessage)
}

type ZPool struct {
	QueuePool []chan ZipperRequest
	Rm        map[uint16]ZipperRouter
	Size      uint8
}

// InitPool initialize the worker pool
func (zp *ZPool) InitPool(ctx context.Context) {
	if common.GlobalConfig.PoolSize == 0 || common.GlobalConfig.QueueSize == 0 {
		logger.OutLog("There is an exception in the PoolSize or QueueSize configuration items in the json configuration file", logger.ERROR)
		panic("The PoolSize or QueueSize configuration item in the system json configuration file has a value of 0 or an abnormal value")
	}
	zp.Size = common.GlobalConfig.PoolSize
	zp.QueuePool = make([]chan ZipperRequest, zp.Size)
	for i := 0; i < int(zp.Size); i++ {
		zp.QueuePool[i] = make(chan ZipperRequest, common.GlobalConfig.QueueSize)
		go zp.RunQueue(ctx, i, zp.QueuePool[i])
	}
}

// AddRouter add message id routing mapping
func (zp *ZPool) AddRouter(id uint16, rt ZipperRouter) {
	if _, ok := zp.Rm[id]; ok {
		logger.OutLog("Message id map insert failed, id = "+strconv.Itoa(int(id)), logger.ERROR)
		panic("There is a conflict in the routing message id mapping")
		return
	}
	zp.Rm[id] = rt
}

// RunQueue open queue job
func (zp *ZPool) RunQueue(ctx context.Context, id int, queue chan ZipperRequest) {
	logger.OutLog("Queue id = "+strconv.Itoa(id)+" is running", logger.INFO)
	for {
		select {
		case req := <-queue:
			zp.ConsumeMessage(req)
		}
	}
}

// ConsumeMessage Consuming messages
func (zp *ZPool) ConsumeMessage(req ZipperRequest) {
	// get the corresponding route according to the message id
	router, ok := zp.Rm[req.GetMsgId()]
	if !ok {
		logger.OutLog("", logger.ERROR)
		return
	}
	msg := router.Handler(req)
	if msg != nil {
		req.GetConnect().Send(msg)
	}
}

// AddQueue add messages to the queue by pseudo load balancing
func (zp *ZPool) AddQueue(c ZipperConnect, m message.ZipperMessage) {
	// generate request
	zReq := NewZRequest(c, m)
	// take the remainder to achieve pseudo load balancing, and stuff the request into the queue
	i := time.Now().Second() % int(zp.Size)
	zp.QueuePool[i] <- zReq
}
