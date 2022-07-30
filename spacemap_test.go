package spacemap

import (
	"github.com/google/go-cmp/cmp"
	"image"
	"testing"
)

func TestSpaceMap(t *testing.T) {
	m := NewSpaceMap()
	rect1 := NewRectangle(10, 10, 100, 100)
	m.Add(rect1)
	for _, test := range []struct {
		Name     string
		Stack    []Shape
		Position image.Point
	}{
		{
			Name:     "Hit",
			Stack:    []Shape{rect1},
			Position: image.Point{20, 20},
		},
		{
			Name:     "Miss",
			Stack:    []Shape{},
			Position: image.Point{-20, -20},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			stack := m.GetStackAt(test.Position.X, test.Position.Y)
			if s := cmp.Diff(stack, test.Stack); len(s) > 0 {
				t.Errorf("Failed stacks differ: %s", s)
			}
		})
	}
}
