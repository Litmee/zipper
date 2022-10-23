package connect

import (
	"context"
	"github.com/Litmee/zipper/common"
	"github.com/Litmee/zipper/logger"
	"github.com/Litmee/zipper/message"
	"strconv"
	"time"
)

// ZipperPool work queue task pool abstraction layer
type ZipperPool interface {
	// InitPool initialize the worker pool
	InitPool(ctx context.Context)
	// AddRouter add message id routing mapping
	AddRouter(id uint16, rt ZipperRouter)
	// runQueue open queue job
	runQueue(ctx context.Context, id int, queue chan ZipperRequest)
	// consumeMessage Consuming messages
	consumeMessage(req ZipperRequest)
	// addQueue add messages to the queue by pseudo load balancing
	addQueue(c ZipperConnect, m message.ZipperMessage)
}

type zPool struct {
	queuePool []chan ZipperRequest
	rm        map[uint16]ZipperRouter
	size      uint8
}

func NewZPool() ZipperPool {
	return &zPool{rm: make(map[uint16]ZipperRouter)}
}

// InitPool initialize the worker pool
func (zp *zPool) InitPool(ctx context.Context) {
	if common.GlobalConfig.PoolSize == 0 || common.GlobalConfig.QueueSize == 0 {
		logger.OutLog("There is an exception in the PoolSize or QueueSize configuration items in the json configuration file", logger.ERROR)
		panic("The PoolSize or QueueSize configuration item in the system json configuration file has a value of 0 or an abnormal value")
	}
	zp.size = common.GlobalConfig.PoolSize
	zp.queuePool = make([]chan ZipperRequest, zp.size)
	for i := 0; i < int(zp.size); i++ {
		zp.queuePool[i] = make(chan ZipperRequest, common.GlobalConfig.QueueSize)
		go zp.runQueue(ctx, i, zp.queuePool[i])
	}
}

// AddRouter add message id routing mapping
func (zp *zPool) AddRouter(id uint16, rt ZipperRouter) {
	if _, ok := zp.rm[id]; ok {
		logger.OutLog("Message id map insert failed, id = "+strconv.Itoa(int(id)), logger.ERROR)
		panic("There is a conflict in the routing message id mapping")
		return
	}
	zp.rm[id] = rt
}

// runQueue open queue job
func (zp *zPool) runQueue(ctx context.Context, id int, queue chan ZipperRequest) {
	logger.OutLog("Queue id = "+strconv.Itoa(id)+" is running", logger.INFO)
	for {
		req := <-queue
		zp.consumeMessage(req)
	}
}

// consumeMessage Consuming messages
func (zp *zPool) consumeMessage(req ZipperRequest) {
	// get the corresponding route according to the message id
	router, ok := zp.rm[req.GetMsgId()]
	if !ok || router == nil {
		logger.OutLog("The route corresponding to the message id = "+strconv.Itoa(int(req.GetMsgId()))+" is not found, the id may be wrong or the route itself is nil", logger.ERROR)
		return
	}
	msg := router.Handler(req)
	if msg != nil {
		req.getConnect().send(msg)
	}
}

// addQueue add messages to the queue by pseudo load balancing
func (zp *zPool) addQueue(c ZipperConnect, m message.ZipperMessage) {
	// generate request
	zReq := NewZRequest(c, m)
	// take the remainder to achieve pseudo load balancing, and stuff the request into the queue
	i := time.Now().Nanosecond() % int(zp.size)
	zp.queuePool[i] <- zReq
}
