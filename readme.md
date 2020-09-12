# 介绍
该项目是一个轻量级GO的TCP服务框架

## 消息封装
定义一个解决TCP粘包问题的封包拆包模块
- 针对Message进行TLV格式的封装
- 针对Message进行TLV格式的拆包

## 多路由模式


## 读写协程分离
实现思路:  
- 添加一个reader和writer通信的channel
- 添加一个writer goroutine
- reader由之前直接发送给客户端改成发送给channel
- 启动reader和writer一起工作

## 消息队列和多任务处理
- 创建一个消息队列
- 创建多任务worker工作池
- 将之前的发送消息，全部改成把消息发送给消息队列和worker工作池


## 连接管理
对于连接数做限制，超过一定数量需要拒绝请求。
- 创建一个连接管理模块（定义、属性、方法）


