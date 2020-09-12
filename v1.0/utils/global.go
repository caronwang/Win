package utils

import (
	"Win/v1.0/iface"
	"encoding/json"
	"io/ioutil"
	"log"
)

/*
 存储一切有关框架的全局参数
*/

type GlobalObj struct {
	TCPServer      iface.IServer //全局server对象
	Host           string        //当前服务器监听IP
	TcpPort        int           //当前服务器监听端口
	Name           string        //当前服务器名称
	Version        string        //当前版本号
	MaxConn        int           //服务器主机允许的最大连接数
	MaxPackageSize uint32        //数据包最大值

	WorkerPoolSize uint32        //任务池大小
	MaxWorkerTaskSize uint32	//任务池最大数量
}

//conf/conf.json中加载自定义参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("v1.0/conf/conf.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}
}

/*
	定义一个全局对象
*/
var GlobalObject *GlobalObj

// 初始化全局对象
func init() {
	log.Println("loading configuration...")
	GlobalObject = &GlobalObj{
		Name:           "TCPServer",
		Version:        "0.1",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1024,
		MaxPackageSize: 4096,
		WorkerPoolSize:  10,
		MaxWorkerTaskSize:1024,
	}

	GlobalObject.Reload()
	//log.Println(GlobalObject)
}
