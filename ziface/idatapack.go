package ziface

//todo 自定义协议
//解决tcp封包拆包
type IDataPack interface {
	//获取包头长度
	GetHeadLen() uint32
	//封包
	Pack(msg IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}
