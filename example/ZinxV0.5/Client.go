package main

import (
	"ohio/znet"
	"fmt"
	"io"
	"net"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	// 1 连接远程服务器，得到conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error,", err)
		return
	}
	for {
		//发送封包的msg消息
		dp := znet.NewDataPack()
		msgBinary, err := dp.Pack(znet.NewMsgPackage(0, []byte("ohio V0.5 client test message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		conn.Write(msgBinary)

		//服务器返回message：id，
		headData := make([]byte, dp.GetHeadLen())
		//先读头
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println("read head error:", err)
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack msgHead error:", err)
			break
		}
		//再读data部分
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
				return
			}

			fmt.Println("----->Recv server message id = ", msg.Id, ",len =", msg.DataLen, ",data = ", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}

}
