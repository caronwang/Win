package core

import (
	"fmt"
	"sync"
)

/*
	AOI地图中的格子类型
*/

type Gird struct {
	//格子ID
	GID int
	//格子左边界坐标
	MinX int
	//格子右边边界
	MaxX int
	//格子上边边界
	MinY int
	//格子下边边界
	MaxY int
	//格子内玩家或物体成员ID集合
	playerIDs map[int]bool
	//保护当前集合的锁
	pIDLock sync.RWMutex
}

func NewGird(gid, minX, maxX, minY, maxY int) *Gird {
	return &Gird{
		GID:       gid,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

//给格子添加一个玩家
func (g *Gird) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

//删除一个玩家
func (g *Gird) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

//得到格子中所有玩家ID
func (g *Gird) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

//打印格子的基本信息
func (g *Gird) String() string {
	return fmt.Sprintf("gid:%v,minX:%v,maxX:%v,minY:%v,maxY:%v,playerIDS:%v,",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
