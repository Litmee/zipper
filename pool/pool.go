package pool

import "zipper/route"

// ZipperPool work queue task pool abstraction layer
type ZipperPool interface {
	InitPool()
}

type ZPool struct {
	Queue []*ZipperQueue
	Rm    map[uint16]route.ZipperRouter
}

// InitPool initialize the worker pool

func (zp *ZPool) InitPool() {

}
