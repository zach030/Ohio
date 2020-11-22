package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试拆包、封包的单元测试
func TestDataPack_Pack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept error:", err)
		}
		go func(conn net.Conn) {
			//处理客户端请求
			dp := NewDataPack()
			for {
				//1 拆包,读head
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head error:", err)
					break
				}
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("server unpack head error:", err)
					return
				}
				if msgHead.GetMsgLen() > 0 {
					//进行第二次读取
					//2 读data内容
					msg := msgHead.(*Message)
					msg.Data = make([]byte, msg.GetMsgLen())
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data failed,error:", err)
						return
					}
					fmt.Println("----> Recv MsgId= ", msg.Id, ", dataLen= ", msg.DataLen, ",data= ", string(msg.Data))
				}
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("client dial failed,", err)
		return
	}
	dp := NewDataPack()
	//封装两个message分别发送
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack message1 failed,", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack message1 failed,", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	select {}
}
