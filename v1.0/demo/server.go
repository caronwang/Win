package main

import (
	"Win/v1.0/iface"
	"Win/v1.0/wnet"
	_ "Win/v1.0/utils"
	"fmt"
	"log"
)

/*
	基于框架开发的应用程序
*/


// ping test 测试
type PingRouter struct {
	wnet.BaseRouter
}

func (p *PingRouter) PreHandle(req iface.IRequest){
	log.Println("prehandle...")
}

func (p *PingRouter) Handle(req iface.IRequest){
	log.Println("handle...")
	log.Println("recv from client: msgId = ",req.GetMsgID(),",data = ",string(req.GetData()))

	err := req.GetConnection().SendMsg(req.GetMsgID(),[]byte(fmt.Sprintln(" got msg. id =",req.GetMsgID(),",data:",string(req.GetData()))))
	if err!=nil{
		log.Println("send message err,",err)
		return
	}
}


// ping test 测试
type HelloRouter struct {
	wnet.BaseRouter
}

func (p *HelloRouter) Handle(req iface.IRequest){
	log.Println("handle...")
	log.Println("recv from client: msgId = ",req.GetMsgID(),",data = ",string(req.GetData()))

	err := req.GetConnection().SendMsg(req.GetMsgID(),[]byte(fmt.Sprintln("hello" )))
	if err!=nil{
		log.Println("send message err,",err)
		return
	}
}


func (p *PingRouter) PostHandle(req iface.IRequest){
	log.Println("Posthandle...")
}

func main(){
	//创建一个server句饼
	s:= wnet.NewServer("win")

	//注册自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloRouter{})


	//启动server
	s.Serve()
}
