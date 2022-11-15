package shared

import (
	"fmt"
	"image"
)

type Alignment int

const (
	NotAligned Alignment = iota
	Vertical
	Horizontal
)

type Shape interface {
	PointIn(x, y int) bool
	Bounds() image.Rectangle
	String() string
}

type Rectangle struct {
	image.Rectangle
	Name string
}

func (r Rectangle) String() string {
	var n string
	if len(r.Name) > 0 {
		n = r.Name + ":"
	}
	return fmt.Sprintf("%sRect(%s->%s)", n, r.Min, r.Max)
}

func (r Rectangle) PointIn(x, y int) bool {
	return (image.Point{x, y}).In(r.Rectangle)
}

func (r Rectangle) Bounds() image.Rectangle {
	return r.Rectangle
}

var _ Shape = (*Rectangle)(nil)

type Op any

type Name string

func NewRectangle(left, top, right, bottom int, ops ...Op) *Rectangle {
	r := &Rectangle{
		Rectangle: image.Rectangle{
			Min: image.Point{
				X: left,
				Y: top,
			},
			Max: image.Point{
				X: right,
				Y: bottom,
			},
		},
	}
	for _, op := range ops {
		switch op := op.(type) {
		case Name:
			r.Name = string(op)
		}
	}
	return r
}
