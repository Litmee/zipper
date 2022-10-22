# zipper

### 项目功能描述
```
并发 TCP 异步消息处理服务框架
```
### 简单的使用例子

#### 1. 在项目中集成 zipper
```
// 后缀的版本号根据情况定
go get github.com/Litmee/zipper@v0.1.1 
```
#### 2. 在项目中使用 zipper

```go
package main

import (
	"github.com/Litmee/zipper/connect"
	"github.com/Litmee/zipper/message"
	"github.com/Litmee/zipper/server"
)

// MyRouter 消息的业务处理路由模型要组合 connect.ZRouter
// 仅此而已即可, 不要试图给你组合后的结构体添加任何属性, 那都是并发不安全的
type MyRouter struct {
	connect.ZRouter
}

// Handler 具体的业务逻辑实现需要重写 connect.ZRouter 的 Handler 方法
// connect.ZipperRequest 是消息传过来封装后的接口模型结构
// connect.ZipperRequest 对用户暴露了 GetMsgId() 获取消息id 和  GetMsgBody() 获取业务消息内容的方法
// 如果需要给发送方返回消息, 则可以通过 NewZMessage(id uint16, body []byte) 该方法构造 message.ZipperMessage 抽象模型, 然后返回即可
// 如果不需要返回消息给发送方, 请务必 return nil
func (my *MyRouter) Handler(req connect.ZipperRequest) message.ZipperMessage {
	s := string(req.GetMsgBody())
	s += " Hello World"
	zMessage := message.NewZMessage(req.GetMsgId(), []byte(s))
	return zMessage
}

// 准备好消息的业务处理环节, 接下来就是服务本身的启动与路由的注册
// 通过 server.NewZServer() 方法构造一个抽象 server.ZipperServer
// 通过  server.ZipperServer 的 AddRouter(id uint16, rt connect.ZipperRouter) 方法注册路由
// id 参数为消息id, rt 为组合 connect.ZRouter 实现的具体业务处理路由模型
// 例子: 当发送方传过来消息 id = 1 的消息时将调用 MyRouter 的 Handler 方法进行业务处理
// 注册完毕消息的路由配置, 调用 server.ZipperServer 的 Run() 方法启动服务
// 注意: 就目前的初衷来看, 是希望集成 zipper 的服务单独部署的
// 所以 Run() 方法执行后正常情况下是持续阻塞的, 其后的业务代码在 zipper 服务整体崩溃前是不会被执行的
func main() {
	zServer := server.NewZServer()
	zServer.AddRouter(1, &MyRouter{})
	zServer.Run()
}
```
#### 3. 默认配置
```
端口: 8066
消息队列容量: 200          ===> 即一个队列里最多容纳的消息数量, 容量满了后会产生阻塞
队列工作池容量: 6          ===> 即消息队列的数量 
最大数据包容量: 1024 字节   ===> 即 TCP 传输二进制消息数据允许的最大字节长度
最大 TCP 链接数: 50        ===> 即服务可接收的最大 TCP 链接数量, 满了后其他 TCP 链接无法连接
日志管道缓冲区容量: 200     ===> 日志服务是通过有缓冲管道异步打印的, 此数值为管道的缓冲区大小
```
#### 4. 自定义配置
```
用户可以通过新增该相对路径下 ./conf/zipper.json 文件进行对 zipper 的自定义配置

zipper.json 有效字段格式如下

Port        ===> 端口(uint16 类型, 注意范围)
QueueSize   ===> 消息队列容量(uint16 类型, 注意范围)
PoolSize    ===> 队列工作池容量(uint8 类型, 注意范围)
MaxPackSize ===> 最大数据包容量(uint16 类型, 注意范围)
MaxConnect  ===> 最大 TCP 链接数(uint16 类型, 注意范围)
LogSize     ===> 日志管道缓冲区容量(uint16 类型, 注意范围)
```

```json
{
  "Port": 8066,
  "QueueSize": 200,
  "PoolSize": 6,
  "MaxPackSize": 1024,
  "MaxConnect": 50,
  "LogSize": 200
}
```