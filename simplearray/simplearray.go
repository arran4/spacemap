package simplearray

import (
	"github.com/arran4/spacemap/shared"
	"sort"
)

type Struct struct {
	Shapes []*shared.Point
}

func New() *Struct {
	return &Struct{
		Shapes: []*shared.Point{},
	}
}

func (sm *Struct) AddAll(s ...shared.Shape) *Struct {
	for _, e := range s {
		sm.Add(e, 0)
	}
	return sm
}

func (sm *Struct) Add(shape shared.Shape, zIndex int) {
	sm.Shapes = append(sm.Shapes, &shared.Point{
		Shape:  shape,
		ZIndex: zIndex,
	})
}

func (sm *Struct) Remove(s shared.Shape) {
	shrink := 0
	for i := range sm.Shapes {
		for sm.Shapes[i] != nil && sm.Shapes[i].Shape == s && len(sm.Shapes)-shrink > i {
			shrink++
			sm.Shapes[i] = sm.Shapes[len(sm.Shapes)-shrink]
		}
	}
	newLen := len(sm.Shapes) - shrink
	for i := newLen; i < len(sm.Shapes); i++ {
		sm.Shapes[i] = nil
	}
	sm.Shapes = sm.Shapes[:newLen]
}

func (sm *Struct) GetStackAt(x int, y int) (result []shared.Shape) {
	var r []*shared.Point
	for i := range sm.Shapes {
		if sm.Shapes[i].PointIn(x, y) {
			r = append(r, sm.Shapes[i])
		}
	}
	sort.Sort(shared.ZSort(r))
	for _, e := range r {
		result = append(result, e.Shape)
	}
	return
}

func (sm *Struct) GetAt(x int, y int) shared.Shape {
	s := sm.GetStackAt(x, y)
	if len(s) > 0 {
		return s[0]
	}
	return nil
}
