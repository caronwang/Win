package wnet

import "Win/v1.0/iface"

type Request struct {
	//已经和客户端建立好的连接
	conn iface.IConnection
	//客户端请求的数据
	//data []byte
	msg iface.IMessage

}


func (r *Request) GetConnection() iface.IConnection{
	return r.conn
}


func (r *Request) GetData() []byte{
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}