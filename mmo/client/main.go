package main

import (
	pb "Win/mmo/proto"
	"Win/v1.0/wnet"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"log"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {

			dp := wnet.NewDatePack()

			headData := make([]byte, dp.GetHeaderLen())
			_, err := io.ReadFull(conn, headData)
			if err != nil {
				if err == io.EOF {
					log.Println("remote connection has been closed!")
				} else {
					log.Println("read head err,", err)
				}

				break
			}

			msg, err := dp.Unpack(headData)
			if err != nil {
				log.Println("unpack  message err,", err)
				break
			}

			if msg.GetMsgLen() > 0 {
				dataBuf := make([]byte, msg.GetMsgLen())
				_, err := io.ReadFull(conn, dataBuf)
				if err != nil {
					log.Println("read data err,", err)
					break
				}
				msg.SetData(dataBuf)
			}

			switch msg.GetMsgId() {
			case 1:
				var data pb.SyncPid
				err := proto.Unmarshal(msg.GetData(), &data)
				if err != nil {
					fmt.Println("Unmarshal err,", err)
				}

				fmt.Println("[玩家上线] --> msgid = ", msg.GetMsgId(), ",玩家id = ", data.Pid)
			case 200:
				var data pb.BroadCast
				err := proto.Unmarshal(msg.GetData(), &data)
				if err != nil {
					fmt.Println("Unmarshal err,", err)
				}
				if data.Tp == 2 {
					fmt.Println("[更新玩家位置] --> msgid = ", msg.GetMsgId(), ",postion = ", data.Data)
				}
			case 201:
				var data pb.SyncPid
				err := proto.Unmarshal(msg.GetData(), &data)
				if err != nil {
					fmt.Println("Unmarshal err,", err)
				}

				fmt.Println("[玩家下线] --> msgid = ", msg.GetMsgId(), ",玩家id = ", data.Pid)
			}
		}

	}()
	time.Sleep(time.Second * 3)
	conn.Close()
	select {}

	//var msgId uint32
	//msgId = 1
	//for{
	//
	//	dp := wnet.NewDatePack()
	//
	//	msg := wnet.NewMessage(msgId%2,[]byte(time.Now().String()))
	//
	//	bindata,err  :=dp.Pack(msg)
	//	if err!=nil{
	//		if err == io.EOF{
	//			log.Println("remote connection has been closed!")
	//		}else{
	//			log.Println("pack msg err,",err)
	//		}
	//
	//		break
	//	}
	//
	//	_,err = conn.Write(bindata)
	//
	//
	//	if err!=nil{
	//		if err == io.EOF{
	//			log.Println("remote connection has been closed!")
	//		}else{
	//			log.Println("write err,",err)
	//		}
	//		break
	//	}
	//
	//	fmt.Println("<-- send msg.id = ",msg.GetMsgId(),",data = ",string(msg.GetData()))
	//	time.Sleep(time.Second)
	//	msgId++
	//}

}
