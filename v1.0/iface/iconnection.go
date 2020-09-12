package iface

import "net"

//定义连接模块的抽象层
type IConnection interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//获取当前连接绑定sockert conn
	GetTCPConnection() *net.TCPConn
	//获取当前连接模块的连接ID
	GetConnId() uint32
	//获取远程客户端的TCP状态 IP Port
	RemoteAddr() net.Addr
	//发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32,data []byte) error
}

//定义哥处理连接业务的方法
type HandleFunc func(*net.TCPConn,[]byte,int) error