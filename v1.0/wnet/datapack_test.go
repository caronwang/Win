package wnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试datapack拆包
func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "0.0.0.0:7777")
	if err != nil {
		fmt.Println("server listen err,", err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err,", err)
			}

			go func(conn2 net.Conn) {
				//处理客户端请求
				dp := NewDatePack()

				for {
					headdata := make([]byte, dp.GetHeaderLen())
					_, err := io.ReadFull(conn, headdata)
					if err != nil {
						fmt.Println(err)
						return
					}

					msghead, err := dp.Unpack(headdata)
					if err != nil {
						fmt.Println("unpack err,", err)
						return
					}

					if msghead.GetMsgLen() > 0 {
						msg := msghead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("unpack data err,", err)
							return
						}

						fmt.Println(
							"-->recv msg id = ", msg.GetMsgId(), ",dataLen = ", msg.GetMsgLen(),
							",data = ", string(msg.GetData()))

					}
				}
			}(conn)




		}
	}()


	//模拟客户端

	conn,err := net.Dial("tcp","127.0.0.1:7777")
	if err != nil {
		fmt.Println("connect server err,", err)
		return
	}

	//创建一个封包对象
	dp:= NewDatePack()

	//模拟粘包
	msg1 := &Message{
		Id: 1,
		DataLen: 5,
		Data: []byte("hello"),
	}

	msg2 := &Message{
		Id: 2,
		DataLen: 5,
		Data: []byte("world"),
	}
	data1,err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack data err,", err)
		return
	}
	data2,_ := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack data err,", err)
		return
	}
	data1 = append(data1,data2...)
	_,err = conn.Write(data1)
	if err != nil {
		fmt.Println("send data err,", err)
		return
	}
	select {

	}
}
