package ziface

//路由的抽象接口，路由的数据是irequest

type IRouter interface {
	//处理业务之前的方法hook
	PreHandle(request IRequest)
	//处理业务的主方法hook
	Handle(request IRequest)
	//处理业务后的方法hook
	PostHandle(request IRequest)
}
