package spacemap

import (
	"github.com/google/go-cmp/cmp"
	"image"
	"testing"
)

func TestSpaceMap(t *testing.T) {
	rect1 := NewRectangle(10, 10, 100, 100)
	for _, test := range []struct {
		Name     string
		Stack    []Shape
		Position image.Point
		SpaceMap *SpaceMap
	}{
		{
			Name:     "Hit",
			Stack:    []Shape{rect1},
			Position: image.Point{20, 20},
			SpaceMap: NewSpaceMap().Add(rect1),
		},
		{
			Name:     "Miss",
			Stack:    []Shape{},
			Position: image.Point{-20, -20},
			SpaceMap: NewSpaceMap().Add(rect1),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			if len(test.SpaceMap.VSplits)*len(test.SpaceMap.HSplits) != len(test.SpaceMap.Stacks) {
				t.Errorf("Incorrect number of cells expected %d ( %d * %d ) but got %d",
					len(test.SpaceMap.VSplits)*len(test.SpaceMap.HSplits), len(test.SpaceMap.VSplits),
					len(test.SpaceMap.HSplits), len(test.SpaceMap.Stacks))
			}
			stack := test.SpaceMap.GetStackAt(test.Position.X, test.Position.Y)
			if s := cmp.Diff(stack, test.Stack); len(s) > 0 {
				t.Errorf("Failed stacks differ: %s", s)
			}
		})
	}
}
