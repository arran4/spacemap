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

func NewTSpaceMap(vTree *Node, hTree *Node) *SpaceMap {
	return &SpaceMap{
		VTree: vTree,
		HTree: hTree,
	}
}

func NewTNode(value int, leftNode, rightNode *Node, heres ...*Here) *Node {
	return &Node{
		Value: value,
		Here:  heres,
		Children: [2]*Node{
			leftNode,
			rightNode,
		},
	}
}

func NewTHere(shape shared.Shape, z int, t Type) *Here {
	return &Here{
		Shape:  shape,
		ZIndex: z,
		Type:   t,
	}
}

func TestSpaceBTreeAdd(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(40, 40, 60, 60, shared.Name("rect2"))
	rect3 := shared.NewRectangle(10, 10, 60, 60, shared.Name("rect3"))
	rect4 := shared.NewRectangle(60, 60, 100, 100, shared.Name("rect4"))
	for _, test := range []struct {
		Name             string
		SpaceMap         func() *SpaceMap
		ExpectedSpaceMap *SpaceMap
	}{
		{
			Name:     "rect1",
			SpaceMap: NSM(rect1),
			ExpectedSpaceMap: NewTSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100, nil, nil, NewTHere(rect1, 0, End)),
					NewTHere(rect1, 0, Begin),
				),
				NewTNode(10,
					nil,
					NewTNode(100, nil, nil, NewTHere(rect1, 0, End)),
					NewTHere(rect1, 0, Begin),
				),
			),
		},
		{
			Name:     "rect1, rect2",
			SpaceMap: NSM(rect1, rect2),
			ExpectedSpaceMap: NewTSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 1, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 1, Begin),
						),
						nil,
						NewTHere(rect1, 0, End)),
					NewTHere(rect1, 0, Begin),
				),
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 1, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 1, Begin),
						),
						nil,
						NewTHere(rect1, 0, End)),
					NewTHere(rect1, 0, Begin),
				),
			),
		},
		{
			Name:     "rect2, rect3",
			SpaceMap: NSM(rect2, rect3),
			ExpectedSpaceMap: NewTSpaceMap(
				NewTNode(40,
					NewTNode(10,
						nil,
						nil,
						NewTHere(rect3, 0, Begin),
					),
					NewTNode(60,
						nil,
						nil,
						NewTHere(rect2, 0, End),
						NewTHere(rect3, 1, End),
					),
					NewTHere(rect2, 0, Begin),
					NewTHere(rect3, 1, Middle),
				),
				NewTNode(40,
					NewTNode(10,
						nil,
						nil,
						NewTHere(rect3, 0, Begin),
					),
					NewTNode(60,
						nil,
						nil,
						NewTHere(rect2, 0, End),
						NewTHere(rect3, 1, End),
					),
					NewTHere(rect2, 0, Begin),
					NewTHere(rect3, 1, Middle),
				),
			),
		},
		{
			Name:     "rect1, rect4",
			SpaceMap: NSM(rect1, rect4),
			ExpectedSpaceMap: NewTSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(60,
							nil,
							nil,
							NewTHere(rect1, 0, Middle),
							NewTHere(rect4, 1, Begin),
						),
						nil,
						NewTHere(rect1, 0, End),
						NewTHere(rect4, 1, End),
					),
					NewTHere(rect1, 0, Begin),
				),
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(60,
							nil,
							nil,
							NewTHere(rect1, 0, Middle),
							NewTHere(rect4, 1, Begin),
						),
						nil,
						NewTHere(rect1, 0, End),
						NewTHere(rect4, 1, End),
					),
					NewTHere(rect1, 0, Begin),
				),
			),
		},
		{
			Name:     "rect1, rect2, rect3",
			SpaceMap: NSM(rect1, rect2, rect3),
			ExpectedSpaceMap: NewTSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 1, End),
								NewTHere(rect3, 2, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 1, Begin),
							NewTHere(rect3, 2, Middle),
						),
						nil,
						NewTHere(rect1, 0, End),
					),
					NewTHere(rect1, 0, Begin),
					NewTHere(rect3, 1, Begin),
				),
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 1, End),
								NewTHere(rect3, 2, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 1, Begin),
							NewTHere(rect3, 2, Middle),
						),
						nil,
						NewTHere(rect1, 0, End),
					),
					NewTHere(rect1, 0, Begin),
					NewTHere(rect3, 1, Begin),
				),
			),
		},
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
