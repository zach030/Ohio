package ziface

//消息管理抽象层

type IMsgHandler interface {
	//调度执行对应的router方法
	DoMsgHandler(request IRequest)
	//为消息添加路由
	AddRouter(msgId uint32, router IRouter)
}
