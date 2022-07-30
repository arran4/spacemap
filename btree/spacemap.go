package btree

import (
	"image"
	"spacemap/shared"
)

type Node struct {
	Position  int
	BecauseOf []shared.Shape
	Alignment shared.Alignment
}

type SpaceMap struct {
	VTree []*Node
	HTree []*Node
}

func (m *SpaceMap) AddAll(shapes ...shared.Shape) *SpaceMap {
	for _, shape := range shapes {
		m.Add(shape)
	}
	return m
}

func (m *SpaceMap) Add(shape shared.Shape) *SpaceMap {
	b := shape.Bounds()
}

func (m *SpaceMap) GetXYPositions(p image.Point) (int, int) {

}

func (m *SpaceMap) GetStackAt(x int, y int) []shared.Shape {

}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{
		VTree: []*Node{},
		HTree: []*Node{},
	}
}
