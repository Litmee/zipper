package server

import (
	"context"
	"github.com/Litmee/zipper/common"
	"github.com/Litmee/zipper/connect"
	"github.com/Litmee/zipper/logger"
	"net"
	"strconv"
)

// ZipperServer service layer interface abstraction
type ZipperServer interface {
	// Run start method
	Run()
	// AddRouter add message id routing mapping
	AddRouter(id uint16, rt connect.ZipperRouter)
}

type zServer struct {
	// message queue worker pool
	pool connect.ZipperPool
}

func NewZServer() ZipperServer {
	return &zServer{
		pool: connect.NewZPool(),
	}
}

// Run start method
func (zs *zServer) Run() {
	// 1. initialize the top-level context, for coroutine safety
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 2. start the log service
	go logger.InitLogServer(ctx)

	// 3. start the message queue worker pool according to the configuration
	zs.pool.InitPool(ctx)

	// 4. monitor
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
		// if the current number of links is already greater than MaxConnect, close
		if !common.Add() {
			conn.Close()
			continue
		}
		// create a Linked Layer Model
		zConnect := connect.NewZConnect(conn, zs.pool)
		go zConnect.Start(ctx)
	}
}

// AddRouter add message id routing mapping
func (zs *zServer) AddRouter(id uint16, rt connect.ZipperRouter) {
	zs.pool.AddRouter(id, rt)
}
