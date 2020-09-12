package wnet

import "Win/v1.0/iface"

//实现router时，先嵌入baserouter类，根据需求对该类进行重写
type BaseRouter struct {

}



//处理conn业务前的钩子方法
func (b *BaseRouter) PreHandle(request iface.IRequest){

}
//处理业务方法
func (b *BaseRouter) Handle(request iface.IRequest){

}

//处理conn业务后的钩子方法
func (b *BaseRouter) PostHandle(request iface.IRequest){

}


