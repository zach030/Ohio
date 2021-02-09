package znet

import "Ohio/ziface"

//实现router时，先嵌入BaseRouter基类，再根据业务对基类方法进行重写
type BaseRouter struct {}


//这些方法为空，因为有的router不希望有pre和post两个业务，建立router时要全部继承BaseRouter
//处理业务之前的方法hook
func (br *BaseRouter)PreHandle(request ziface.IRequest){}

//处理业务的主方法hook
func (br *BaseRouter)Handle(request ziface.IRequest){}

//处理业务后的方法hook
func (br *BaseRouter)PostHandle(request ziface.IRequest){}

