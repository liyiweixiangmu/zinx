package utils

import (
	"encoding/json"
	"github.com/liyiweixiangmu/zinx/ziface"
	"io/ioutil"
)

type _GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	/*
		Zinx
	*/
	Version          string
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前业务工作worker池的Goroutine数量
	MaxWorkerTaskLen uint32 // zinx 框架允许用户最多开辟多少个worker
}

//定义一个全聚德对外GlobalObj
var GlobalObject *_GlobalObj

/**
 * @Description:提供一个init方法，初始化当前的GlobalObject
 */
func init() {
	// 若配置文件没有加载，默认的值
	GlobalObject = &_GlobalObj{
		Name:             "ZinxServer",
		Version:          "V0.4",
		TcpPort:          8999,
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   //
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的任务的数量的最大值
	}

	//应该尝试从conf/zinx.json
	GlobalObject.Reload()
}

/**
 * @Description:从zinx.json去加载用于自定义的参数
 * @receiver g
 */
func (g *_GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析道struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
