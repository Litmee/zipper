package server

import (
	"context"
	"net"
	"strconv"
	"zipper/common"
	"zipper/connect"
	"zipper/logger"
)

// ZipperServer service layer interface abstraction
type ZipperServer interface {
	// Run start method
	Run()
	// AddRouter add message id routing mapping
	AddRouter(id uint16, tr connect.ZipperRouter)
}

type zServer struct {
	// message queue worker pool
	pool connect.ZipperPool
}

func NewZServer() ZipperServer {
	return &zServer{
		pool: &connect.ZPool{Rm: make(map[uint16]connect.ZipperRouter)},
	}
}

func (zs *zServer) Run() {
	// 1. initialize the top-level context, for coroutine safety
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 2. start the log service
	go logger.InitLogServer(ctx)

	// 3. start the message queue worker pool according to the configuration
	zs.pool.InitPool(ctx)

	// monitor
	addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(int(common.GlobalConfig.Port)))
	if err != nil {
		logger.OutLog(err.Error(), logger.ERROR)
		panic(err)
	}
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.OutLog(err.Error(), logger.ERROR)
		panic(err)
	}
	for {
		conn, err := tcp.AcceptTCP()
		if err != nil {
			logger.OutLog(err.Error(), logger.ERROR)
			continue
		}
		// create a Linked Layer Model
		zConnect := connect.NewZConnect(conn, zs.pool)
		go zConnect.Start(ctx)
	}
}

func (zs *zServer) AddRouter(id uint16, rt connect.ZipperRouter) {
	zs.pool.AddRouter(id, rt)
}
