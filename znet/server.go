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
	//当前server的消息管理模块，绑定消息id与处理业务api
	MsgHandler  ziface.IMsgHandler
	ConnManager ziface.IConnManager
	//connection 建立后的hook 函数
	OnConnStart func(connection ziface.IConnection)
	//connection 销毁前的hook 函数
	OnConnStop func(connection ziface.IConnection)
}

//The writer who wrote the codes below is a fat and lazy pig. -That's true.But,he has a beautiful girlfriend that everyone envies.
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name:%s, listener at IP:%s,Port:%d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s,MaxConn %d, MaxPackageSize %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 0 开启消息队列，工作池
		s.MsgHandler.StartWorkerPool()
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
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// 给客户端响应
				fmt.Println("========= Too many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			//客户端已经建立连接,做一个回写业务
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态，已经开辟的连接信息进行停止或者回收
	fmt.Println("[STOP] Zinx Server is stop")
	s.ConnManager.Clear()
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()
	//TODO 做一些启动服务后的业务
	//处于阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Successfully!")
}

//初始化server的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

//注册start hook func
func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart!=nil{
		fmt.Println("----> Call On Conn Start func......")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop!=nil{
		fmt.Println("----> Call On Conn Stop func......")
		s.OnConnStop(connection)
	}
}


