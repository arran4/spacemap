package spacebtree

import (
	"github.com/google/go-cmp/cmp"
	"image"
	"spacemap/shared"
	"testing"
)

func NSM(shapes ...shared.Shape) func() *SpaceMap {
	return func() *SpaceMap {
		return NewSpaceMap().AddAll(shapes...)
	}
}

func TestSpaceBTreeAdd(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100)
	//rect2 := shared.NewRectangle(40, 40, 60, 60)
	//rect3 := shared.NewRectangle(10, 10, 60, 60)
	//rect4 := shared.NewRectangle(60, 60, 100, 100)
	for _, test := range []struct {
		Name             string
		SpaceMap         func() *SpaceMap
		ExpectedSpaceMap *SpaceMap
	}{
		{
			Name:     "rect1",
			SpaceMap: NSM(rect1),
			ExpectedSpaceMap: &SpaceMap{
				VTree: &Node{
					Value: 10,
					Here: []*Here{
						{
							Shape:  rect1,
							ZIndex: 0,
							Type:   Begin,
						},
					},
					Children: [2]*Node{
						nil,
						{
							Value: 100,
							Here: []*Here{
								{
									Shape:  rect1,
									ZIndex: 0,
									Type:   End,
								},
							},
							Children: [2]*Node{},
						},
					},
				},
				HTree: &Node{
					Value: 10,
					Here: []*Here{
						{
							Shape:  rect1,
							ZIndex: 0,
							Type:   Begin,
						},
					},
					Children: [2]*Node{
						nil,
						{
							Value: 100,
							Here: []*Here{
								{
									Shape:  rect1,
									ZIndex: 0,
									Type:   End,
								},
							},
							Children: [2]*Node{},
						},
					},
				},
			},
		},
		//{
		//	Name:     "rect1, rect2",
		//	SpaceMap: NSM(rect1, rect2),
		//},
		//{
		//	Name:     "rect1, rect3",
		//	SpaceMap: NSM(rect1, rect3),
		//},
		//{
		//	Name:     "rect1, rect4",
		//	SpaceMap: NSM(rect1, rect4),
		//},
	} {
		t.Run(test.Name, func(t *testing.T) {
			sm := test.SpaceMap()
			if s := cmp.Diff(sm, test.ExpectedSpaceMap); len(s) > 0 {
				t.Errorf("Failed stacks differ: %s", s)
			}
		})
	}
}

func TestSpaceBTreeEndToEnd(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100)
	rect2 := shared.NewRectangle(40, 40, 60, 60)
	rect3 := shared.NewRectangle(10, 10, 60, 60)
	rect4 := shared.NewRectangle(60, 60, 100, 100)
	for _, test := range []struct {
		Name     string
		Stack    []shared.Shape
		Position image.Point
		SpaceMap func() *SpaceMap
	}{
		{
			Name:     "Hit",
			Stack:    []shared.Shape{rect1},
			Position: image.Point{20, 20},
			SpaceMap: NSM(rect1),
		},
		{
			Name:     "Hit Low Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Min,
			SpaceMap: NSM(rect1),
		},
		{
			Name:     "Hit High Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Max,
			SpaceMap: NSM(rect1),
		},
		{
			Name:     "Miss -- Near Hit High Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Max.Add(image.Pt(1, 1)),
			SpaceMap: NSM(rect1),
		},
		{
			Name:     "Miss",
			Stack:    []shared.Shape{},
			Position: image.Point{-20, -20},
			SpaceMap: NSM(rect1),
		},
		{
			Name:     "Hit first with 2 overlapping",
			Stack:    []shared.Shape{rect1},
			Position: image.Point{20, 20},
			SpaceMap: NSM(rect1, rect2),
		},
		{
			Name:     "Hit both with 2 overlapping",
			Stack:    []shared.Shape{rect1, rect2},
			Position: image.Point{50, 50},
			SpaceMap: NSM(rect1, rect2),
		},
		{
			Name:     "Hit both with 2 overlapping same start",
			Stack:    []shared.Shape{rect1, rect3},
			Position: image.Point{20, 20},
			SpaceMap: NSM(rect1, rect3),
		},
		{
			Name:     "Hit both with 2 overlapping same end",
			Stack:    []shared.Shape{rect1, rect4},
			Position: image.Point{90, 90},
			SpaceMap: NSM(rect1, rect4),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			sm := test.SpaceMap()
			stack := sm.GetStackAt(test.Position.X, test.Position.Y)
			if s := cmp.Diff(stack, test.Stack); len(s) > 0 {
				t.Errorf("Failed stacks differ: %s", s)
			}
		})
	}
}
