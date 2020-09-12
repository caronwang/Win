package iface

/*
	路由抽象接口
	路由里的数据都是IRequest
*/

type IRouter interface {
	//处理conn业务前的钩子方法
	PreHandle(request IRequest)
	//处理业务方法
	Handle(request IRequest)
	//处理conn业务后的钩子方法
	PostHandle(request IRequest)
}
