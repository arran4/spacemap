package simplearray

import (
	"sort"
	"spacemap/shared"
)

type Point struct {
	ZIndex int
	shared.Shape
}

type Struct struct {
	Shapes []*Point
}

func New() *Struct {
	return &Struct{
		Shapes: []*Point{},
	}
}

func (sm *Struct) AddAll(s ...shared.Shape) *Struct {
	for _, e := range s {
		sm.Add(e)
	}
	return sm
}

func (sm *Struct) Add(s shared.Shape) *Struct {
	sm.Shapes = append(sm.Shapes, &Point{
		Shape: s,
	})
	return sm
}

func (sm *Struct) Remove(s shared.Shape) *Struct {
	for i := range sm.Shapes {
		if sm.Shapes[i] == s {
			sm.Shapes[i] = sm.Shapes[len(sm.Shapes)-1]
			sm.Shapes = sm.Shapes[:len(sm.Shapes)-1]
		}
	}
	return sm
}

func (sm *Struct) GetStackAt(x int, y int) (result []shared.Shape) {
	var r []*Point
	for i := range sm.Shapes {
		if sm.Shapes[i].PointIn(x, y) {
			r = append(r, sm.Shapes[i])
		}
	}
	sort.Sort(ZSort(r))
	for _, e := range r {
		result = append(result, e.Shape)
	}
	return
}

type ZSort []*Point

func (Z ZSort) Len() int {
	return len(Z)
}

func (Z ZSort) Less(i, j int) bool {
	return Z[i].ZIndex < Z[j].ZIndex
}

func (Z ZSort) Swap(i, j int) {
	Z[i], Z[j] = Z[j], Z[i]
}

var _ sort.Interface = ZSort(nil)
