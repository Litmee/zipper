package route

type ZipperRouter interface {
	PreHandler()
	Handler()
	BackHandler()
}

type ZRouter struct{}

func (zr *ZRouter) PreHandler() {}

func (zr *ZRouter) Handler() {}

func (zr *ZRouter) BackHandler() {}
