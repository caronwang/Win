# 介绍
该项目是一个轻量级GO的TCP服务框架。



目录说明

服务框架

```
└── v1.0										//服务框架目录
    ├── conf								//配置文件
    │   └── conf.json
    ├── demo								//测试程序
    │   ├── client.go					//客户端测试程序
    │   └── server.go					//服务端程序
    ├── iface									//抽象层
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
    │   └── global.go					//全局配置管理
    └── wnet									//实现层
        ├── connMgr.go				//连接管理
        ├── connection.go			//客户端连接对象
        ├── datapack.go				//消息解析器
        ├── datapack_test.go	//消息解析器-单元测试
        ├── message.go				//消息对象
        ├── msgHandler.go			//消息处理器
        ├── request.go				//客户端请求对象
        ├── router.go					//路由对象
        └── server.go					//服务对象
```

MMO应用

```shell
├── mmo
│   ├── apis									//api功能实现
│   │   └── move.go
│   ├── client
│   │   └── main.go						//客户端测试程序
│   ├── conf
│   │   └── conf.json					//配置文件
│   ├── core
│   │   ├── aio_test.go				//AIO管理-单元测试
│   │   ├── aoi.go						//AIO管理
│   │   ├── gird.go						//网格对象
│   │   ├── player.go					//玩家对象
│   │   └── worldMgr.go				//世界管理器
│   ├── gen.sh
│   ├── main.go								//MMO服务程序
│   └── proto
│       ├── msg.pb.go
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

- 添加一个reader和writer通信的channel
- 添加一个writer goroutine
- reader由之前直接发送给客户端改成发送给channel
- 启动reader和writer一起工作



## 消息队列和多任务处理

- 创建一个消息队列
- 创建多任务worker工作池
- 将之前的发送消息，全部改成把消息发送给消息队列和worker工作池



## 连接管理

- 对于连接数做限制，超过一定数量需要拒绝请求。
- 连接创建后/连接销毁前提供Hook接口



## 连接属性配置




## 消息协议

| MsgID | 事件                                                   | 信息                                                         |
| ----- | ------------------------------------------------------ | ------------------------------------------------------------ |
| 1     | SynPid：<br />同步玩家本次登录的ID<br />发起者：server | Pid：玩家信息                                                |
| 2     | Talk：<br />同步玩家聊天信息<br />发起者：client       | Content：聊天信息                                            |
| 3     | Move:<br />玩家移动坐标数据<br />发起者：client        | X:X 坐标<br />Y:Y 坐标<br />Z:Z坐标<br />V:角度              |
| 200   | BroadCast<br />广播消息<br />发起者:server             | Pid:玩家<br />Topic:消息类型（1:世界聊天，2:坐标，3:动作）<br />Content:消息 |
| 201   | SynPid<br />广播消息 掉线/消失视野<br />发起者：Server | Pid：玩家ID                                                  |
| 202   | SynPos<br />将玩家信息同步给周围人                     | Player: 玩家信息（Pid：玩家ID，Position：位置信息）          |
|       |

