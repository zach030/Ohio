package znet

import (
	"Zinx/ziface"
	"fmt"
	"strconv"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个msgid所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度执行对应的router方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//找到msgid
	msgId := request.GetMsgID()
	handler,ok := m.Apis[msgId]
	if !ok{
		fmt.Println("api msgId = ",request.GetMsgID(),", is not found, need register router")
		return
	}
	//map中调度对应业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加路由
func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//判断当前msg绑定api是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeat api, msg id=" + strconv.Itoa(int(msgId)))
	}
	//绑定
	m.Apis[msgId] = router
	fmt.Println("add api msg id = ", msgId, " success")
}
