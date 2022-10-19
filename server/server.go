package server

type Server interface {
	// Run 启动方法
	Run()
}

type ZipperServer struct {
	// 服务名称
	Name string
	// ip 地址
	Ip string
	// 端口号
	Port uint16
}

func (zs *ZipperServer) Run() {

}
