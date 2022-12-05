package common

import (
	"sync"
)

// curConnectNum current number of TCP connections
var curConnectNum uint16

// lock
var lock sync.Mutex

// Add link count add
func Add() bool {
	lock.Lock()
	defer lock.Unlock()
	curConnectNum++
	if curConnectNum > GlobalConfig.MaxConnect {
		return false
	}
	return true
}

// Cut link count reduce
func Cut() {
	lock.Lock()
	defer lock.Unlock()
	curConnectNum--
}
