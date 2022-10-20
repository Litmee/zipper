package pool

// ZipperQueue 消息队列抽象层
type ZipperQueue interface {
	Do()
}

type ZQueue struct {
}
