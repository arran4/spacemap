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
