package iface


type IDatePack interface {
	//获取包头长度
	GetHeaderLen() uint32
	//封包
	Pack(msg IMessage)([]byte,error)
	//拆包
	Unpack([]byte) (IMessage,error)

}