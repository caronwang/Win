package iface

/*

客户端请求的连接信息，和请求数据包装到一个request中

*
 */

type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMsgID() uint32
}
