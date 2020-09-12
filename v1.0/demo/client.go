package main

import (
	"Win/v1.0/wnet"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main(){

	conn,err :=net.Dial("tcp","127.0.0.1:8888")
	if err!=nil{
		log.Fatal(err)
	}

	go func() {
		for {

			dp := wnet.NewDatePack()


			headData := make([]byte,dp.GetHeaderLen())
			_,err :=io.ReadFull(conn,headData)
			if err!=nil{
				log.Println("read head err,",err)
				break
			}

			msg,err := dp.Unpack(headData)
			if err!=nil{
				log.Println("unpack  message err,",err)
				break
			}

			if msg.GetMsgLen() >0{
				dataBuf := make([]byte,msg.GetMsgLen())
				_,err :=io.ReadFull(conn,dataBuf)
				if err!=nil{
					log.Println("read data err,",err)
					break
				}
				msg.SetData(dataBuf)
			}

			fmt.Println("-->",string(msg.GetData()))
		}

	}()

	var msgId uint32
	msgId = 1
	for{

		dp := wnet.NewDatePack()

		msg := wnet.NewMessage(msgId%2,[]byte(time.Now().String()))

		bindata,err  :=dp.Pack(msg)
		if err!=nil{
			log.Println("pack msg err,",err)
			break
		}

		_,err = conn.Write(bindata)
		if err == io.EOF{
			break
		}

		if err!=nil{
			log.Println("write err,",err)
			continue
		}

		fmt.Println("<-- send msg.id = ",msg.GetMsgId(),",data = ",string(msg.GetData()))
		time.Sleep(time.Second)
		msgId++
	}


}
