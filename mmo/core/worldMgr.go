package core

import (
	"sync"
)

/*
	服务器管理模块
*/
type WorldMgr struct {
	AoiMgr  *AIOManager
	Players map[int32]*Player
	pLock   sync.RWMutex
}

var WorldMgrObj *WorldMgr

func init() {
	WorldMgrObj = &WorldMgr{
		AoiMgr:  NewAIOManager(0, 1000, 5, 0, 1000, 5),
		Players: make(map[int32]*Player),
	}
}

func (w *WorldMgr) AddPlayer(player *Player) {
	w.pLock.Lock()
	defer w.pLock.Unlock()
	w.Players[player.Pid] = player

	w.AoiMgr.AddToGirdByPos(int(player.Pid), player.X, player.Z)
}

func (w *WorldMgr) RemovePlayerByPid(pid int32) {
	player := w.Players[pid]
	w.AoiMgr.RemoveToGirdByPos(int(player.Pid), player.X, player.Z)

	w.pLock.Lock()
	defer w.pLock.Unlock()
	delete(w.Players, pid)
}

func (w *WorldMgr) GetPlayerByPid(pid int32) *Player {
	w.pLock.RLock()
	defer w.pLock.RUnlock()
	return w.Players[pid]
}
