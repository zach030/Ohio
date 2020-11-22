package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"net"
)

//当前连接模块
type Connection struct {
	//当前连接的socket tcp 套接字
	Conn *net.TCPConn
	//连接的id
	ConnID uint32
	//当前连接的状态
	isClosed bool
	//告知当前连接已停止的channel
	ExitChan chan bool
	//当前链接处理的方法router
	Router ziface.IRouter
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit, remote Addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("receive buf err ", err)
			continue
		}
		req := Request{
			conn: c,
			data: buf,
		}
		//执行注册的路由方法
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
		//从路由中，找到注册绑定的conn对应的router调用
	}
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("connection start,conID=", c.ConnID)
	//启动从当前连接的读数据业务
	//TODO 启动从当前连接写数据的业务
	go c.StartReader()
}

//停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("connection stop(), conID=", c.ConnID)
	//如果当前连接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()
	//关闭管道,回收资源
	close(c.ExitChan)
}

//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接的id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的tcp状态
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
