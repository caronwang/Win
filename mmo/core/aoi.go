package core

import (
	"fmt"
)

type AIOManager struct {
	//区域左边界坐标
	MinX int
	//区域右边边界
	MaxX int
	//X方向格子的数量
	CntsX int
	//区域上边边界
	MinY int
	//区域下边边界
	MaxY int
	//Y方向格子的数量
	CntsY int
	//当前区域有哪些格子map  -key=格子ID，value=格子对象
	girds map[int]*Gird
}

func NewAIOManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AIOManager {
	aoiMgr := &AIOManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		girds: make(map[int]*Gird),
	}

	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsY + x
			aoiMgr.girds[gid] = NewGird(
				gid, aoiMgr.MinX+x*aoiMgr.girdWidth(), aoiMgr.MinX+(x+1)*aoiMgr.girdWidth(),
				aoiMgr.MinY+x*aoiMgr.girdLength(), aoiMgr.MinY+(y+1)*aoiMgr.girdWidth(),
			)
		}
	}
	fmt.Println(aoiMgr)
	return aoiMgr
}

func (a *AIOManager) girdWidth() int {

	return (a.MaxX - a.MinX) / a.CntsX
}

func (a *AIOManager) girdLength() int {

	return (a.MaxY - a.MinY) / a.CntsY
}

func (a *AIOManager) GetSuroundingGirdByGird(gid int) (girds []*Gird) {
	if _, ok := a.girds[gid]; !ok {
		return
	}

	girds = append(girds, a.girds[gid])

	idx := gid % a.CntsX
	if idx > 0 {
		girds = append(girds, a.girds[gid-1])
	}

	if idx < a.CntsX-1 {
		girds = append(girds, a.girds[gid+1])
	}

	gidsX := make([]int, 0, len(girds))
	for _, v := range girds {
		gidsX = append(gidsX, v.GID)
	}

	for _, v := range gidsX {
		idy := v / a.CntsY
		if idy > 0 {
			girds = append(girds, a.girds[v-a.CntsX])
		}
		if idy < a.CntsY-1 {
			girds = append(girds, a.girds[v+a.CntsX])
		}
	}

	return
}

func (a *AIOManager) GetGirdByPos(x, y float32) int {
	idx := (int(x) - a.MinX) / a.girdWidth()
	idy := (int(y) - a.MinY) / a.girdLength()
	return idy*a.CntsX + idx
}

func (a *AIOManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	gid := a.GetGirdByPos(x, y)
	girds := a.GetSuroundingGirdByGird(gid)

	for _, v := range girds {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
	}
	return
}

//添加一个player到一个gird中
func (a *AIOManager) AddPidToGird(pid, gid int) {
	a.girds[gid].Add(pid)
}

//删除一个player到一个gird中
func (a *AIOManager) RemovePidToGird(pid, gid int) {
	a.girds[gid].Remove(pid)
}

//通过gid获取全部的playerID
func (a *AIOManager) GetPidByGird(gid int) (pids []int) {
	pids = a.girds[gid].GetPlayerIDs()
	return
}

//通过坐标将player添加到一个格子中
func (a *AIOManager) AddToGirdByPos(pid int, x, y float32) {
	gid := a.GetGirdByPos(x, y)
	gird := a.girds[gid]
	if gird == nil {
		panic("gird not found")
	}

	gird.Add(pid)
	return
}

//通过坐标将player在一个格子中删除
func (a *AIOManager) RemoveToGirdByPos(pid int, x, y float32) {
	gid := a.GetGirdByPos(x, y)
	gird := a.girds[gid]
	gird.Remove(pid)
	return
}

func (a *AIOManager) String() string {
	s := fmt.Sprintf("[AIOManager]\nMinX:%v,MaxX:%v,cntsX:%v,MinY:%v,MaxY:%v,cntsY:%v\n",
		a.MinX, a.MaxX, a.CntsX, a.MinY, a.MaxY, a.CntsY,
	)
	s += "girds:\n"
	for _, gird := range a.girds {
		s += "\t" + fmt.Sprintln(gird)
	}

	return s
}
