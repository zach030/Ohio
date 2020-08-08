package ziface

import "net"

//定义连接的抽象层
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接的工作
	Stop()
	//获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接的id
	GetConnID() uint32
	//获取远程客户端的tcp状态
	RemoteAddr() net.Addr
	//发送数据，发送给远程客户端
	Send(data []byte) error
}

//定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
