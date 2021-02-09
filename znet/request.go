package znet

import "Ohio/ziface"

type Request struct {
	//已经与客户端建立好的连接
	conn ziface.IConnection
	//客户端请求的数据
	data ziface.IMessage
}

//得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

//得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.data.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.data.GetMsgId()
}
