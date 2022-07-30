package shared

import "image"

type Alignment int

const (
	NotAligned Alignment = iota
	Vertical
	Horizontal
)

type Shape interface {
	PointIn(x, y int) bool
	Bounds() image.Rectangle
}

type Rectangle image.Rectangle

func (r Rectangle) PointIn(x, y int) bool {
	return (image.Point{x, y}).In(image.Rectangle(r))
}

func (r Rectangle) Bounds() image.Rectangle {
	return image.Rectangle(r)
}

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
