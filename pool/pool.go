package pool

import "zipper/common"

// ZipperPool work queue task pool abstraction layer
type ZipperPool interface {
}

type ZPool struct {
	queue []*ZipperQueue
	size  uint16
}

// InitPool initialize the worker pool
func (zp *ZPool) InitPool() {
	zp.size = common.GlobalConfig.QueueSize
}
