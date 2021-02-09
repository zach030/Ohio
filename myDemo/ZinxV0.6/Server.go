package main

import (
	"Ohio/ziface"
	"Ohio/znet"
	"fmt"
)

//ping test
type PingRouter struct {
	znet.BaseRouter
}

// test handle router
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router: handle")
	//先读取客户端数据，再回写
	fmt.Println("recv client msg id = ", request.GetMsgID(),
		" data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1,[]byte("ping"))
	if err != nil {
		fmt.Println("call back  ping error:", err)
	}
}

type HelloOhioRouter struct {
	znet.BaseRouter
}
// test handle router
func (p *HelloOhioRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call router: Hello ohio router handle")
	//先读取客户端数据，再回写
	fmt.Println("recv client msg id = ", request.GetMsgID(),
		" data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1,[]byte("Hello ohio"))
	if err != nil {
		fmt.Println("call back  ping error:", err)
	}
}
/*基于Ohio框架开发的服务器应用程序*/
func main() {
	// 1 创建server句柄
	s := znet.NewServer("ohio V0.6")
	// 2 添加自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloOhioRouter{})
	// 3 启动server
	s.Serve()
}
