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
	err := request.GetConnection().SendMsg(1, []byte("ping"))
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
	err := request.GetConnection().SendMsg(1, []byte("Hello ohio"))
	if err != nil {
		fmt.Println("call back  ping error:", err)
	}
}

func DoConnBegin(connection ziface.IConnection) {
	fmt.Println("==========> Do connection Begin")
	if err := connection.SendMsg(202, []byte("Do connection Begin")); err != nil {
		fmt.Println("Call hook func failed:", err)
		return
	}
	//给当前连接设置属性
	fmt.Println("Set conn ...")
	connection.SetProperty("Name", "zach")
	connection.SetProperty("blog", "zach030")
	connection.SetProperty("Addr", "Nanjing")
}

func DoConnStop(connection ziface.IConnection) {
	fmt.Println("==========> Do connection Stop, conn id = ", connection.GetConnID())

	if name, err := connection.GetProperty("Name"); err == nil {
		fmt.Println("Name is:", name)
	}
}

/*基于Ohio框架开发的服务器应用程序*/
func main() {
	// 1 创建server句柄
	s := znet.NewServer("ohio V0.10")
	// 2 添加自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloOhioRouter{})
	// 3 注册hook回调func
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnStop)
	// 4 启动server
	s.Serve()
}
