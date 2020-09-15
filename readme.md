# 介绍
该项目是一个轻量级GO的TCP服务框架。



目录说明

服务框架

```
└── v1.0	//服务框架目录
    ├── conf	//配置文件
    │   └── conf.json
    ├── demo	//测试程序
    │   ├── client.go	//客户端测试程序
    │   └── server.go	//服务端程序
    ├── iface	//抽象层
    │   ├── IRequest.go				
    │   ├── iDataPack.go
    │   ├── iMsgHandler.go
    │   ├── iServer.go
    │   ├── iconnMgr.go
    │   ├── iconnection.go
    │   ├── imessage.go
    │   └── irouter.go
    ├── readme.md
    ├── utils
    │   └── global.go	//全局配置管理
    └── wnet	//实现层
        ├── connMgr.go	//连接管理
        ├── connection.go	//客户端连接对象
        ├── datapack.go	//消息解析器
        ├── datapack_test.go	//消息解析器-单元测试
        ├── message.go	//消息对象
        ├── msgHandler.go	//消息处理器
        ├── request.go	//客户端请求对象
        ├── router.go	//路由对象
        └── server.go	//服务对象
```

MMO应用

```shell
├── mmo
│   ├── apis	//api功能实现
│   │   └── move.go
│   ├── client
│   │   └── main.go	//客户端测试程序
│   ├── conf
│   │   └── conf.json	//配置文件
│   ├── core
│   │   ├── aio_test.go	//AIO管理-单元测试
│   │   ├── aoi.go	//AIO管理
│   │   ├── gird.go	//网格对象
│   │   ├── player.go	//玩家对象
│   │   └── worldMgr.go	//世界管理器
│   ├── gen.sh
│   ├── main.go	//MMO服务程序
│   └── proto
│       ├── msg.pb.go	//msg对象
│       └── msg.proto
```



## 服务器



```go
//服务对象
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip地址
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的Port
	Port int
	//server消息管理模块，用来绑定msgid和业务处理API关系
	MsgHandler iface.IMsgHandler
	//连接管理器
	ConnManager iface.IConnManager
	//连接创建后的Hook函数
	OnConnStart func(conn iface.IConnection)
	//连接销毁前的Hook函数
	OnConnStop func(conn iface.IConnection)
}


//服务接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能
	AddRouter(msgId uint32, router IRouter)
  //连接管理
	GetConnMgr() IConnManager
	//注册OnConnStart钩子方法
	SetOnConnStart(func(connection IConnection))
	//注册OnConnStop钩子方法
	SetOnConnStop(func(connection IConnection))
	//调用OnConnStart钩子方法
	CallOnConnStart(conn IConnection)
	//调用OnConnStart钩子方法
	CallOnConnStop(conn IConnection)
}

```





## 消息封装

定义一个解决TCP粘包问题的封包拆包模块

改造思路：

- 针对Message进行TLV格式的封装
- 针对Message进行TLV格式的拆包

消息包的结构分为`包头`和`包体`两部分。其中包头又由`消息长度`(4字节)、`消息类型`(4字节)组成，包体为发送的具体消息，采用protobuf格式。



```go
message.go

//消息对象
type Message struct {
	Id uint32	//消息ID
	DataLen uint32 	//消息长度
	Data []byte	//消息内容
}

IDataPack.go

//消息解析器接口
type IDatePack interface {
	//获取包头长度
	GetHeaderLen() uint32
	//封包
	Pack(msg IMessage)([]byte,error)
	//拆包
	Unpack([]byte) (IMessage,error)
}

wnet/request.go
//客户端请求的连接信息，和请求数据包装到一个request对象中
type Request struct {
	//已经和客户端建立好的连接
	conn iface.IConnection
	//客户端请求的数据
	//data []byte
	msg iface.IMessage

}

iface/IRequest.go

//请求接口
type IRequest interface {
  //得到request对象绑定的connection对象
	GetConnection() IConnection
  //设置request的数据
	GetData() []byte
  //设置request的消息类型
	GetMsgID() uint32
}
```







## 多路由模式

针对不同的消息，执行不同的函数方法，将接口提供给框架使用者调用。

```go
//实现router时，先嵌入baserouter类，根据需求对该类进行重写
type BaseRouter struct {}

//处理conn业务前的钩子方法
func (b *BaseRouter) PreHandle(request iface.IRequest){
}

//处理业务方法
func (b *BaseRouter) Handle(request iface.IRequest){
}

//处理conn业务后的钩子方法
func (b *BaseRouter) PostHandle(request iface.IRequest){
}
```

服务端调用

```go
func main(){
  //创建服务器句柄
	s := NewServer("game")
	
  ...
	s.AddRouter(3, &api.MoveApi{}) //移动
	
	...
	//启动服务
	s.Serve()
}


//玩家移动,继承BaseRouter
type MoveApi struct {
	BaseRouter
}


//重写Handle方法
func (*MoveApi) Handle(request IRequest) {
	 ....
}
```




## 读写协程分离

改造思路：

- 添加一个reader和writer通信的channel
- 添加一个writer goroutine
- reader由之前直接发送给客户端改成发送给channel
- 启动reader和writer一起工作



## 消息队列和多任务处理

改造思路：

- 创建多任务worker工作池，每个worker绑定一个消息channel接受connection reader解析后的数据封装成的request对象
- 将之前的发送消息给客户端，全部改成把消息发送给消息队列和worker工作池。发送给worker的策略采用connId对工作池大小取余的方式选择。

```go
wnet/msgHandler.go
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


iface/iMsgHandler.go
/*
	消息处理抽象层
*/
type IMsgHandler interface {
	//调度、执行对应的router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32,router IRouter)
	//启动工作池
	StartWorkerPool()
	//处理发送给IMsgHandler的消息
	SendMsgToTaskQ(IRequest)
}


```







## 连接管理

连接管理的作用：

- 对于连接数做限制，超过一定数量需要拒绝请求。
- 连接创建后/连接销毁前提供Hook接口



```go
wnet/connMgr.go

//连接管理对象
type ConnManager struct {
	//记录connId和连接对象的映射关系
	cons map[uint32]iface.IConnection

	lock sync.RWMutex
}


iface/iconnMgr.go
/*
	连接管理抽象模块
*/
type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据connID获取链接
	Get(connID uint32) (IConnection, error)
	//得到当前连接数
	Len() int
	//清除所有链接
	ClearConn()
}


```





## 连接属性配置

对于框架使用者，需要为其提供连接属性的功能，将连接和一些业务信息作绑定，比如玩家ID等

```go
wnet/connection.go
type Connection struct {
  ...
	//连接属性集合
	Property map[string]interface{}
	//连接属性锁
	pLock sync.RWMutex
}

iface/iconnection.go
//定义连接模块的抽象层
type IConnection interface {
	...
	//设置连接属性
	SetProperty(key string, value interface{})
	//获取连接属性
	GetProperty(key string) (interface{}, error)
	//删除连接属性
	RemoveProperty(key string)
}

```

MMO服务中，当玩家登陆时，服务端接受到连接，给客户端发送一个随机分配的PlayerID

```go
mmo/main.go

//当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)
	...
  //将pid和连接作绑定
	conn.SetProperty("pid", player.Pid)
	//fmt.Println(core.WorldMgrObj.GetPlayerByPid(player.Pid))
	fmt.Println("=====> Player pidId = ", player.Pid, " arrived ====")
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

```




## MMO消息协议

| MsgID | 事件                                                   | 信息                                                         |
| ----- | ------------------------------------------------------ | ------------------------------------------------------------ |
| 1     | SynPid：<br />同步玩家本次登录的ID<br />发起者：server | Pid：玩家信息                                                |
| 2     | Talk：<br />同步玩家聊天信息<br />发起者：client       | Content：聊天信息                                            |
| 3     | Move:<br />玩家移动坐标数据<br />发起者：client        | X:X 坐标<br />Y:Y 坐标<br />Z:Z坐标<br />V:角度              |
| 200   | BroadCast<br />广播消息<br />发起者:server             | Pid:玩家<br />Topic:消息类型（1:世界聊天，2:坐标，3:动作）<br />Content:消息 |
| 201   | SynPid<br />广播消息 掉线/消失视野<br />发起者：Server | Pid：玩家ID                                                  |
| 202   | SynPos<br />将玩家信息同步给周围人<br />发起者：Server | Player: 玩家信息（Pid：玩家ID，Position：位置信息）          |
|       |                                                        |                                                              |

部分消息对应的protobuf结构定义

```protobuf
syntax="proto3";                //Proto协议
package pb;                     //当前包名
option csharp_namespace="Pb";   //给C#提供的选项

//同步客户端玩家ID
message SyncPid{
  int32 Pid=1;
}

//玩家位置
message Position{
  float X=1;
  float Y=2;
  float Z=3;
  float V=4;
}

//玩家广播数据
message BroadCast{
  int32 Pid=1;
  int32 Tp=2;
  oneof Data {
    string Content=3;
    Position P=4;
    int32 ActionData=5;
  }
}
```

