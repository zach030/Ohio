package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"ohio/utils"
	"ohio/ziface"
)

// 封包，拆包 模块,处理TCP粘包问题
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度
func (d *DataPack) GetHeadLen() uint32 {
	//data len 4 byte, data id 4 byte
	return 8
}

//封包 ; data-len----data-id------data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	//data len --> buf
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		fmt.Println("write data_len in buf failed,", err)
		return nil, err
	}
	//data id --> buf
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		fmt.Println("write data_id in buf failed,", err)
		return nil, err
	}
	//data content --> buf
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		fmt.Println("write data in buf failed,", err)
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包:将包的head信息读出来，根据data长度进行读
func (d *DataPack) Unpack(binData []byte) (ziface.IMessage, error) {
	//创建io-reader
	dataBuff := bytes.NewReader(binData)
	//解压head信息，得到datalen和id
	msg := &Message{}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		fmt.Println("read data-len failed,", err)
		return nil, err
	}
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		fmt.Println("read data-id failed,", err)
		return nil, err
	}
	//判断长度是否已超出最大允许长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too long msg data error")
	}
	return msg, nil
}
