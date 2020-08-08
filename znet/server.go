package znet

import (
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
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listening at IP : %s, Port:%d, is starting\n\n", s.IP, s.Port)

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
			fmt.Println("listen ", s.IPVersion, " error ", err)
			return
		}
		fmt.Println("Start zinx Server successfully, ", s.Name, " Listening...")
		// 3 阻塞的等待客户端连接，处理客户端业务
		for {
			//如果有客户端连接，会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//客户端已经建立连接,做一个回写业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf err,", err)
						continue
					}
					//回写功能
					if _, err := conn.Write(buf[0:cnt]); err != nil {
						fmt.Println("write back buf error,", err)
						continue
					}
				}
			}()
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

//初始化server的方法
func NewServer(name string) Server {
	s := Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
