package spacemap

import "image"

type SpaceMap struct {
}

func (m SpaceMap) Add(shape Shape) {

}

func (m SpaceMap) GetStackAt(x int, y int) []Shape {
	return []Shape{}
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{}
}

type Shape interface{}

type Rectangle image.Rectangle

var _ Shape = (*Rectangle)(nil)

func NewRectangle(left, top, right, bottom int) *Rectangle {
	return &Rectangle{
		Min: image.Point{
			left,
			top,
		},
		Max: image.Point{
			right,
			bottom,
		},
	}
}
