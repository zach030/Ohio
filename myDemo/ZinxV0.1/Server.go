package main

import (
	"Zinx/znet"
)

/*基于zinx框架开发的服务器应用程序*/
func main(){
	// 1 创建server句柄
	s := znet.NewServer("zinx V0.1")
	// 2 启动server
	s.Serve()
}
