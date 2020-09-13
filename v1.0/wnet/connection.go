package wnet

import (
	"Win/v1.0/iface"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

//IConnection实现
type Connection struct {
	//当前Connection属于那个server
	Server iface.IServer

	//当前连接的套接字
	Conn *net.TCPConn

	//链接ID
	ConnId uint32

	//当前连接状态
	isClosed bool

	//告知当前连接是否已经推出的channel
	ExitChan chan bool

	//reader与writer通信channel
	msgChan chan []byte

	//消息处理管理
	MsgHandler iface.IMsgHandler

	//连接属性集合
	Property map[string]interface{}

	//连接属性锁
	pLock sync.RWMutex
}

func NewConnection(server iface.IServer, conn *net.TCPConn, connId uint32, msgHandler iface.IMsgHandler) *Connection {
	c := &Connection{
		Server:     server,
		Conn:       conn,
		ConnId:     connId,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		Property:   make(map[string]interface{}),
	}

	//将conn加入到GetConnMgr中
	server.GetConnMgr().Add(c)
	return c
}

//读业务方法
func (c *Connection) StartReader() {
	log.Println("[ reader goroutine is running... ]")
	defer log.Println("[ reader is exit,remote addr : ", c.RemoteAddr(), ", id : ", c.ConnId, " ]")

	defer c.Stop()

	for {
		//创建一个拆包对象
		dp := NewDatePack()

		//读取客户端 msg head 二进制流8字节
		headData := make([]byte, dp.GetHeaderLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {

			if err == io.EOF {
				log.Println("remote connection has been closed")
			} else {
				log.Println("read msg head error,", err)
			}

			break
		}

		//拆包

		msg, err := dp.Unpack(headData)
		if err != nil {
			log.Println("unpack msg head error,", err)
			break
		}

		//根据datalen读取data，放入msg.data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				if err == io.EOF {
					log.Println("remote connection has been closed")
				} else {
					log.Println("read msg head error,", err)
				}
			}
		}
		msg.SetData(data)

		//得到当前conn数据的request请求数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		//执行消息
		//c.MsgHandler.DoMsgHandler(req)
		c.MsgHandler.SendMsgToTaskQ(req)

	}

}

/*
	写消息goroutine,将消息发送给客户端
*/
func (c *Connection) StartWriter() {
	log.Println("[ writer goroutine is running... ]")
	defer log.Println("[ writer is exit,remote addr : ", c.RemoteAddr(), ", id : ", c.ConnId, " ]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("send message err,", err)
				return
			}

		case <-c.ExitChan:
			//reader已经推出
			return
		}
	}
}

//启动连接
func (c *Connection) Start() {
	log.Printf(" receive connection from %v - id:%d ", c.RemoteAddr(), c.ConnId)

	//启动当前连接度数据的业务
	go c.StartReader()

	//启动写数据业务
	go c.StartWriter()
	//运行Hook函数
	c.Server.CallOnConnStart(c)

}

//停止连接
func (c *Connection) Stop() {
	log.Printf("connId:%d is stopped!", c.ConnId)

	if c.isClosed {
		return
	}
	c.isClosed = true
	//运行Hook函数
	c.Server.CallOnConnStop(c)

	c.Conn.Close()
	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)

	//将当前连接从ConnMgr中去除
	c.Server.GetConnMgr().Remove(c)

}

//获取当前连接绑定sockert conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

//获取远程客户端的TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed!")
	}

	//将data进行封包
	dp := NewDatePack()

	binMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		log.Println("pack message err,", err)
		return err
	}

	c.msgChan <- binMsg

	return nil
}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.pLock.Lock()
	defer c.pLock.Unlock()

	c.Property[key] = value
}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.pLock.RLock()
	defer c.pLock.RUnlock()
	if v, ok := c.Property[key]; !ok {
		return nil, errors.New("key = " + key + " not exist!")
	} else {
		return v, nil
	}

}

//删除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.pLock.Lock()
	defer c.pLock.Unlock()

	delete(c.Property, key)
}
