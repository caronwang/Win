package wnet

import (
	"Win/v1.0/iface"
	"errors"
	"log"
	"sync"
)

type ConnManager struct {
	cons map[uint32]iface.IConnection
	lock sync.RWMutex
}

//创建当前链接
func NewConnManager() *ConnManager {
	return &ConnManager{
		cons: make(map[uint32]iface.IConnection),
	}
}

//添加链接
func (c *ConnManager) Add(conn iface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cons[conn.GetConnId()] = conn
	log.Println("connId =", conn.GetConnId(), "  add to ConnManager successfully: conn num =", c.Len())
}

//删除链接
func (c *ConnManager) Remove(conn iface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	//删除
	delete(c.cons, conn.GetConnId())
	log.Println("connId =", conn.GetConnId(), "  remove from ConnManager successfully: conn num =", c.Len())
}

//根据connID获取链接
func (c *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	conn, ok := c.cons[connID]
	if !ok {
		return nil, errors.New("connID not found")
	}
	return conn, nil
}

//得到当前连接数
func (c *ConnManager) Len() int {
	return len(c.cons)
}

//清除所有链接
func (c *ConnManager) ClearConn() {
	c.lock.Lock()
	defer c.lock.Unlock()

	//删除conn，停止工作
	for connId, conn := range c.cons {
		conn.Stop()
		delete(c.cons, connId)
	}
}
