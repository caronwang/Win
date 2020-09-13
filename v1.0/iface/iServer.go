package iface

//定义一个服务器接口

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能
	AddRouter(msgId uint32, router IRouter)

	GetConnMgr() IConnManager

	//注册OnConnStart钩子方法
	SetOnConnStart(func(connection IConnection))

	//注册OnConnStop钩子方法
	SetOnConnStop(func(connection IConnection))

	//调用OnConnStart钩子方法
	CallOnConnStart(conn IConnection)

	//调用OnConnStart钩子方法
	CallOnConnStop(conn IConnection)
}
