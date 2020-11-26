package ziface

//连接管理模块
type IConnManager interface {
	Add(connection IConnection)
	Remove(connection IConnection)
	Get(id uint32)(IConnection,error)
	Len()int
	Clear()
}
