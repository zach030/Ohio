package ziface

//定义一个服务器接口
type IServer interface {
	//启动
	Start()
	//停止
	Stop()
	//运行
	Serve()
	//路由：给当前服务注册一个路由方法，供客户端的连接使用
	AddRouter(id uint32,router IRouter)
	GetConnManager() IConnManager
}

