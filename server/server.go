package server

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"zipper/common"
	"zipper/logger"
	"zipper/pool"
	"zipper/route"
)

// ZipperServer service layer interface abstraction
type ZipperServer interface {
	// Run start method
	Run()
	// AddRouter add message id routing mapping
	AddRouter(id uint16, tr route.ZipperRouter)
}

type zServer struct {
	// message queue worker pool
	pool pool.ZipperPool
}

func NewZServer() ZipperServer {
	return &zServer{
		pool: &pool.ZPool{},
	}
}

func (zs *zServer) Run() {
	// 1. initialize the top-level context, for coroutine safety
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 2. start the log service
	go logger.InitLogServer(ctx)

	// 3. start the message queue worker pool according to the configuration
	// zPool := &pool.ZPool{}
	// zPool.InitPool()
	// 2. initialize the service object

	// 监听
	addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(int(common.GlobalConfig.Port)))
	if err != nil {
		logger.OutLog(err.Error())
		panic(err)
	}
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.OutLog(err.Error())
		panic(err)
	}
	for {
		conn, err := tcp.AcceptTCP()
		if err != nil {
			logger.OutLog(err.Error())
			continue
		}
		// 创建链接层模型
		fmt.Println(conn)
	}
}

func (zs *zServer) AddRouter(id uint16, rt route.ZipperRouter) {
	// zs.pool.AddRouter(id, rt)
}
