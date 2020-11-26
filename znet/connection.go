package znet

import (
	"Zinx/utils"
	"Zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

//当前连接模块
type Connection struct {
	//所属server
	TcpServer ziface.IServer
	//当前连接的socket tcp 套接字
	Conn *net.TCPConn
	//连接的id
	ConnID uint32
	//当前连接的状态
	isClosed bool
	//告知当前连接已停止的channel
	ExitChan chan bool
	//用于读写协程缓冲的消息通信
	msgChan chan []byte
	//消息管理msgid与处理业务api
	MsgHandler ziface.IMsgHandler
}

//初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	//加入到manager中
	c.TcpServer.GetConnManager().Add(c)
	return c
}

//连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID=", c.ConnID, "[Reader is exit], remote Addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端数据到buf中
		//创建拆包，解包
		dp := NewDataPack()
		//读客户端msg head 二进制流 8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error:", err)
			break
		}
		//拆包，得到msgid，datalen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error:", err)
			break
		}
		var data []byte
		//再次读取，放在msg data中
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)
		req := Request{
			conn: c,
			data: msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

//写消息协程，专门发送给客户端消息
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit]")
	defer c.Stop()
	for {
		//阻塞等待channel消息
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			//代表reader已退出，则writer需退出
			return
		}
	}
}

//提供send msg方法，将要发送给客户端的数据先封包
func (c *Connection) SendMsg(msgid uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed when send msg")
	}
	//将data进行封包
	dp := NewDataPack()
	msgBinary, err := dp.Pack(NewMsgPackage(msgid, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgid)
		return err
	}
	//将数据发送给客户端
	c.msgChan <- msgBinary
	return nil
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("connection start,conID=", c.ConnID)
	//启动从当前连接的读数据业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()
}

//停止连接，结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("connection stop(), conID=", c.ConnID)
	//如果当前连接已关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()
	//告知writer退出
	c.ExitChan <- true
	//关闭管道,回收资源
	close(c.msgChan)
	close(c.ExitChan)

	c.TcpServer.GetConnManager().Remove(c)
}

//获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接的id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的tcp状态
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
