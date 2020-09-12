package wnet

import (
	"Win/v1.0/iface"
	"Win/v1.0/utils"
	"log"
	"strconv"
)

/*
	消息处理模块实现
*/
type MsgHandler struct {
	//存放每个msgID对应的处理方法
	apis map[uint32] iface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	//业务工作worker池的woker数量
	WorkerPoolSize uint32
}


func NewMsgHandler() *MsgHandler{
	m := &MsgHandler{
		apis: make(map[uint32]iface.IRouter),
		WorkerPoolSize:utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan iface.IRequest,utils.GlobalObject.WorkerPoolSize),
	}
	m.StartWorkerPool()
	return m
}


//调度、执行对应的router消息处理方法
func (m *MsgHandler) DoMsgHandler(request iface.IRequest){
	//从request中找到msgID
	handler ,ok := m.apis[request.GetMsgID()]
	if !ok{
		log.Println("msgID", request.GetMsgID(),"not found,you should register it first!")
		return
	}
	//根据msgID调度对应的router
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (m *MsgHandler) AddRouter(msgID uint32,router iface.IRouter){
	//当前的msg绑定的API是否已经存在
	if _,ok := m.apis[msgID];ok{
		panic("repeat add api,msgId = "+strconv.Itoa(int(msgID)) )
	}
	//添加msg与API的处理关系
	m.apis[msgID] = router
}


//启动一个worker工作池(开启工作池的动作只能发生一次)
func (m *MsgHandler) StartWorkerPool(){
	for i:=0 ;i<int(m.WorkerPoolSize);i++{
		//给当前worker对应的channel消息队列开辟空间
		m.TaskQueue[i] = make(chan iface.IRequest,utils.GlobalObject.MaxWorkerTaskSize)
		//启动当前的worker，阻塞等待消息
		go m.StartWorker(i,m.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (m *MsgHandler) StartWorker(workerID int,taskQ chan iface.IRequest){
	log.Println("workerID = ",workerID,"is started...")

	//等待消息队列的消息
	for {
		select {
		//如果有消息过来
		case req := <-taskQ:
			m.DoMsgHandler(req)
		}
	}
}

//处理发送给MsgHandler的消息
func (m *MsgHandler) SendMsgToTaskQ(req iface.IRequest){
	wokerId := req.GetConnection().GetConnId() % m.WorkerPoolSize
	log.Println("assgin connID = ",req.GetConnection().GetConnId(),", msgId = ",req.GetMsgID(),"=> workerid = ",wokerId)
	m.TaskQueue[wokerId]<-req
}