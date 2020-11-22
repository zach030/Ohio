package main

import (
	"Zinx/ziface"
	"Zinx/znet"
	"fmt"
)

//ping test
type PingRouter struct {
	znet.BaseRouter
}

// test pre-handle router
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call router: pre-handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error:", err)
	}
}

// test handle router
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router: handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping...\n"))
	if err != nil {
		fmt.Println("call back  ping error:", err)
	}
}

// test post-handle router
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call router: after-handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error:", err)
	}
}

/*基于zinx框架开发的服务器应用程序*/
func main() {
	// 1 创建server句柄
	s := znet.NewServer("zinx V0.3")
	// 2 添加自定义router
	s.AddRouter(&PingRouter{})
	// 3 启动server
	s.Serve()
}
