package utils

import (
	"Zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

// 存储一切有关zinx的全局参数，由zinx.json配置

type GlobalObj struct {
	TcpServer      ziface.IServer
	Host           string //服务器监听的ip
	TcpPort        int    //服务器监听的端口
	Name           string
	Version        string
	MaxConn        int    //允许的最大连接数
	MaxPackageSize uint32 // 数据包最大值
}

var GlobalObject *GlobalObj

// 加载json
func (s *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObj{})
	if err != nil {
		panic(err)
	}
}

func init() {
	//如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	//尝试从json中加载自定义参数
	GlobalObject.Reload()
}
