package ziface

//将请求的消息封装到message中
type IMessage interface {
	GetMsgId() uint32 //获取消息id
	GetMsgLen() uint32
	GetData() []byte
	SetMsgId(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
