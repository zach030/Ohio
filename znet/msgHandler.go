package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"fmt"
	"strconv"
)

//消息处理模块的实现
type MsgHandle struct {
	//存放每个msgid所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//worker池的数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

//调度执行对应的router方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//找到msgid
	msgId := request.GetMsgID()
	handler, ok := m.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), ", is not found, need register router")
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

//启动worker工作池
func (m *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 分配 worker的channel管道
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前worker
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

//启动一个worker流程
func (m *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerId, " is Starting ... ")
	//阻塞等待对应channel的消息
	for {
		select {
		//如果有消息过来，出列一个request并执行
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//平均分配,找worker
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add Conn id = ", request.GetConnection().GetConnID(), ", request Msg ID = ", request.GetMsgID(),
		" to Worker:", workerId)
	//发送给对应worker的taskqueue
	m.TaskQueue[workerId] <- request
}
