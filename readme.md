# 介绍
该项目是一个轻量级GO的TCP服务框架



## 消息封装

定义一个解决TCP粘包问题的封包拆包模块
- 针对Message进行TLV格式的封装
- 针对Message进行TLV格式的拆包



## 多路由模式



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
|       |                                                        |                                                              |

