package main

import (
	api "Win/mmo/apis"
	"Win/mmo/core"
	. "Win/v1.0/iface"
	. "Win/v1.0/wnet"
	"fmt"
)

//当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)
	//同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()
	core.WorldMgrObj.AddPlayer(player)
	conn.SetProperty("pid", player.Pid)
	//fmt.Println(core.WorldMgrObj.GetPlayerByPid(player.Pid))
	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
}

//当客户端关闭连接的时候的hook函数
func OnConnecionClose(conn IConnection) {
	//获取玩家
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("=====> Player pidId = ", pid, " leaving ====")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	if player == nil {
		fmt.Println("player not found")
		return
	}
	player.Leave()
	fmt.Println("=====> Player pidId = ", player.Pid, " left ====")
}

func main() {
	//创建服务器句柄
	s := NewServer("game")

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnecionClose)

	//注册路由
	s.AddRouter(3, &api.MoveApi{}) //移动

	//启动服务
	s.Serve()
}
