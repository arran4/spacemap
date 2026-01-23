package shared

import "sort"

type Point struct {
	ZIndex int
	Shape
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

type PointArray []*Point

func (pa PointArray) Remove(shape Shape) ([]*Point, int) {
	shrink := 0
	for i := range pa {
		for pa[i].Shape == shape && len(pa)-shrink > i {
			shrink++
			pa[i] = pa[len(pa)-shrink]
		}
	}
	// Avoid memory leak by niling out the pointers in the removed section of the underlying array
	for i := len(pa) - shrink; i < len(pa); i++ {
		pa[i] = nil
	}
	return pa[:len(pa)-shrink], shrink
}

func (pa PointArray) Insert(point *Point) []*Point {
	a := pa
	p := 0
	for ; p < len(a) && point.ZIndex >= a[p].ZIndex; p++ {
	}
	a = append(a, nil)
	copy(a[p+1:], a[p:])
	a[p] = point
	return a
}
