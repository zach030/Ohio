package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"net"
)

//定义一个server的服务器模块
type Server struct {
	//名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//端口
	Port int
	//当前server添加router
	Router ziface.IRouter
}

//The writer who wrote the codes below is a fat and lazy pig. -That's true.But,he has a beautiful girlfriend that everyone envies.
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s, listener at IP:%s,Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s,MaxConn %d, MaxPackageSize %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	go func() {
		// 1 获取一个tcp addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error: ", err)
			return
		}
		fmt.Println("Start zinx Server successfully, ", s.Name, " Listening...")
		var cid uint32
		cid = 0
		// 3 阻塞的等待客户端连接，处理客户端业务
		for {
			//如果有客户端连接，会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}
			//客户端已经建立连接,做一个回写业务
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态，已经开辟的连接信息进行停止或者回收
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()
	//TODO 做一些启动服务后的业务
	//处于阻塞状态
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Successfully!")
}

//初始化server的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
