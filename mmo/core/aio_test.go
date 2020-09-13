package core

import (
	"fmt"
	"testing"
)

func TestAIOManager(t *testing.T) {
	//初始化
	aio := NewAIOManager(0, 500, 5, 0, 500, 5)

	//fmt.Println(aio)
	c := aio.GetSuroundingGirdByGird(0)
	fmt.Println(len(c))
	fmt.Println(c)
}
