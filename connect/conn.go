package connect

import "context"

// ZipperConnect TCP connection layer abstraction
type ZipperConnect interface {
	// start up
	start(ctx context.Context)
}
