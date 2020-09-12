package wnet

import (
	"Win/v1.0/iface"
	"Win/v1.0/utils"
	"fmt"
	"log"
	"net"
)

//定义一个iServer的接口实现
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip地址
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的Port
	Port int
	//当前server添加一个router,server注册的连接对应的处理业务
	//Router iface.IRouter

	//server消息管理模块，用来绑定msgid和业务处理API关系
	MsgHandler iface.IMsgHandler

}




func (s *Server) Start(){

	log.Printf("server [%v] Listening %s:%d is starting\n",utils.GlobalObject.Name,utils.GlobalObject.Host,utils.GlobalObject.TcpPort)

	go func() {
		//获取一个TCP的addr
		addr,err:=net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err!=nil{
			log.Fatal(err)
		}

		//监听服务器地址
		lsner,err :=net.ListenTCP(s.IPVersion,addr)
		if err!=nil{
			log.Fatal(err)
		}
		log.Println("server start success")
		//阻塞的等待客户端连接,处理客户端业务
		var cid uint32
		cid = 0
		for{
			conn,err :=lsner.AcceptTCP()
			if err!=nil{
				log.Println("accept con err,",err)
				continue
			}
			cid++
			//将该处理新连接的业务方法和conn进行绑定
			c := NewConnection(conn,cid,s.MsgHandler)
			go c.Start()

		}
	}()
}

func (s *Server) Stop(){

}

func (s *Server) Serve(){
	s.Start()

	// TODO:启动服务器后的工作

	select {

	}
}

func (s *Server) AddRouter(msgId uint32,router iface.IRouter) {
	s.MsgHandler.AddRouter(msgId,router)
	log.Println("add Router success!")
}

/*
 	初始化Server模块
*/
func NewServer(name string) iface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
	}
	return s
}