package ziface

//IRequest接口

type IRequest interface {
	//得到当前连接
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte
	//得到消息id
	GetMsgID() uint32
}
