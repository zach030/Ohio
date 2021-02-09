package utils

import (
	"Ohio/ziface"
	"encoding/json"
	"io/ioutil"
)

// 存储一切有关Ohio的全局参数，由Ohio.json配置

type GlobalObj struct {
	TcpServer        ziface.IServer
	Host             string //服务器监听的ip
	TcpPort          int    //服务器监听的端口
	Name             string
	Version          string
	MaxConn          int    //允许的最大连接数
	MaxPackageSize   uint32 // 数据包最大值
	WorkerPoolSize   uint32 //当前工作池goroutine大小
	MaxWorkerTaskLen uint32 //允许开辟的最多worker数
}

var GlobalObject *GlobalObj

// 加载json
func (s *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("myDemo/OhioV0.5/conf/ohio.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, s)
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
		Name:           "OhioServerApp",
		Version:        "V0.7",
		MaxConn:        1000,
		MaxPackageSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	//尝试从json中加载自定义参数
	GlobalObject.Reload()
}
