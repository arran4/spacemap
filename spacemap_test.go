package spacemap

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"image"
	"strings"
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
				LogStructure(t, test.SpaceMap)
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

func LogStructure(t *testing.T, spaceMap *SpaceMap) {
	b := strings.Builder{}
	b.WriteString("\n____\t")
	for _, v := range spaceMap.VSplits {
		b.WriteString(fmt.Sprintf("p: %d,\t", v.Position))
	}
	b.WriteString("\n")
	for _, h := range spaceMap.HSplits {
		b.WriteString(fmt.Sprintf("p: %d,\t", h.Position))
		for _, v := range spaceMap.VSplits {
			if sq, ok := spaceMap.Stacks[SplitCoordination{h, v}]; ok {
				b.WriteString(fmt.Sprintf("c: %d,\t", len(sq)))
			} else {
				b.WriteString(fmt.Sprintf("null,\t"))
			}
		}
		b.WriteString("\n")
	}
	t.Log(b.String())
}
