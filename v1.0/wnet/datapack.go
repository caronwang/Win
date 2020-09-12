package wnet

import (
	"Win/v1.0/iface"
	"Win/v1.0/utils"
	"bytes"
	"encoding/binary"
	"errors"
)

type DatePack struct {}


//拆封包实例初始化方法
func NewDatePack() *DatePack{
	return &DatePack{}
}

//获取包头长度
func (d *DatePack) GetHeaderLen() uint32{
	//dataLen(4)+ ID(4)
	return 8
}

//封包
func (d *DatePack) Pack(msg iface.IMessage)([]byte,error){
	//创建一个存放byte字节流的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	//datalen写入databuf
	if err := binary.Write(dataBuf,binary.BigEndian,msg.GetMsgLen());err!=nil{
		return nil,err
	}

	//dataid写入databuf
	if err := binary.Write(dataBuf,binary.BigEndian,msg.GetMsgId());err!=nil{
		return nil,err
	}
	//data数据写入databuf
	if err := binary.Write(dataBuf,binary.BigEndian,msg.GetData());err!=nil{
		return nil,err
	}
	return dataBuf.Bytes(),nil
}


//拆包
func (d *DatePack) Unpack(data []byte) ( iface.IMessage,error){
	//创建一个ioReader
	dataBuf := bytes.NewReader(data)

	msg := &Message{}

	//读到dataLen信息
	if err := binary.Read(dataBuf,binary.BigEndian,&msg.DataLen);err!=nil{
		return nil,err
	}

	//读到dataId信息
	if err := binary.Read(dataBuf,binary.BigEndian,&msg.Id);err!=nil{
		return nil,err
	}

	//包长度判断是否超出允许范围
	if utils.GlobalObject.MaxPackageSize >0 && msg.DataLen >utils.GlobalObject.MaxPackageSize{
		return nil,errors.New("msg data too large")
	}

	return msg,nil
}