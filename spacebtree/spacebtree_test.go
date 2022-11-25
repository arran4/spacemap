package spacebtree

import (
	"github.com/google/go-cmp/cmp"
	"image"
	"reflect"
	"spacemap/shared"
	"testing"
)

func NSMBalanced(shapes ...shared.Shape) func() *SpaceMap {
	return func() *SpaceMap {
		return NewSpaceMap().AddAll(shapes...)
	}
}

func NSMUnbalanced(shapes ...shared.Shape) func() *SpaceMap {
	return func() *SpaceMap {
		return NewSpaceMap().Unbalance().AddAll(shapes...)
	}
}

func NewTUnbalancedSpaceMap(vTree *Node, hTree *Node) *SpaceMap {
	return &SpaceMap{
		VTree: vTree,
		HTree: hTree,
	}
}

func NewTBalancedSpaceMap(vTree *Node, hTree *Node) *SpaceMap {
	s := &SpaceMap{
		VTree:    vTree,
		HTree:    hTree,
		Balanced: true,
	}
	s.VTree.RecalculateDepth(0)
	s.HTree.RecalculateDepth(0)
	return s
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
			Name:     "unbalanced rect1",
			SpaceMap: NSMUnbalanced(rect1),
			ExpectedSpaceMap: NewTUnbalancedSpaceMap(
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
			Name:     "unbalanced rect1, rect2",
			SpaceMap: NSMUnbalanced(rect1, rect2),
			ExpectedSpaceMap: NewTUnbalancedSpaceMap(
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
			Name:     "unbalanced rect2, rect3",
			SpaceMap: NSMUnbalanced(rect2, rect3),
			ExpectedSpaceMap: NewTUnbalancedSpaceMap(
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
			Name:     "unbalanced rect1, rect4",
			SpaceMap: NSMUnbalanced(rect1, rect4),
			ExpectedSpaceMap: NewTUnbalancedSpaceMap(
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
			Name:     "unbalanced rect1, rect2, rect3",
			SpaceMap: NSMUnbalanced(rect1, rect2, rect3),
			ExpectedSpaceMap: NewTUnbalancedSpaceMap(
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
		{
			Name:     "balanced rect1",
			SpaceMap: NSMBalanced(rect1),
			ExpectedSpaceMap: NewTBalancedSpaceMap(
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
			Name:     "balanced rect1, rect2",
			SpaceMap: NSMBalanced(rect1, rect2),
			ExpectedSpaceMap: NewTBalancedSpaceMap(
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
			Name:     "balanced rect2, rect3",
			SpaceMap: NSMBalanced(rect2, rect3),
			ExpectedSpaceMap: NewTBalancedSpaceMap(
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
			Name:     "balanced rect1, rect4",
			SpaceMap: NSMBalanced(rect1, rect4),
			ExpectedSpaceMap: NewTBalancedSpaceMap(
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
			Name:     "balanced rect1, rect2, rect3",
			SpaceMap: NSMBalanced(rect1, rect2, rect3),
			ExpectedSpaceMap: NewTBalancedSpaceMap(
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
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(40, 40, 60, 60, shared.Name("rect2"))
	rect3 := shared.NewRectangle(10, 10, 60, 60, shared.Name("rect3"))
	rect4 := shared.NewRectangle(60, 60, 100, 100, shared.Name("rect4"))
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
			SpaceMap: NSMBalanced(rect1),
		},
		{
			Name:     "Hit Low Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Min,
			SpaceMap: NSMBalanced(rect1),
		},
		{
			Name:     "Hit High Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Max,
			SpaceMap: NSMBalanced(rect1),
		},
		{
			Name:     "Miss -- Near Hit High Border",
			Stack:    []shared.Shape{rect1},
			Position: rect1.Max.Add(image.Pt(1, 1)),
			SpaceMap: NSMBalanced(rect1),
		},
		{
			Name:     "Miss",
			Stack:    []shared.Shape{},
			Position: image.Point{-20, -20},
			SpaceMap: NSMBalanced(rect1),
		},
		{
			Name:     "Hit first with 2 overlapping",
			Stack:    []shared.Shape{rect1},
			Position: image.Point{20, 20},
			SpaceMap: NSMBalanced(rect1, rect2),
		},
		{
			Name:     "Hit both with 2 overlapping",
			Stack:    []shared.Shape{rect1, rect2},
			Position: image.Point{50, 50},
			SpaceMap: NSMBalanced(rect1, rect2),
		},
		{
			Name:     "Hit both with 2 overlapping same start",
			Stack:    []shared.Shape{rect1, rect3},
			Position: image.Point{20, 20},
			SpaceMap: NSMBalanced(rect1, rect3),
		},
		{
			Name:     "Hit both with 2 overlapping same end",
			Stack:    []shared.Shape{rect1, rect4},
			Position: image.Point{90, 90},
			SpaceMap: NSMBalanced(rect1, rect4),
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

func TestNode_VerticalRotate(t *testing.T) {
	tests := []struct {
		name      string
		N         *Node
		NF        func(r *Node, direction Direction)
		direction Direction
		want      *Node
	}{
		{
			name: "Right",
			N: &Node{
				Value: 12,
				Children: [2]*Node{
					{
						Value: 10,
						Children: [2]*Node{
							{
								Value: 8,
								Children: [2]*Node{
									{Value: 4},
									{Value: 9},
								},
							},
							{Value: 11},
						},
					},
					{Value: 14},
				},
			},
			direction: 0,
			want: &Node{
				Value: 12,
				Children: [2]*Node{
					{
						Value: 8,
						Children: [2]*Node{
							{Value: 4},
							{
								Value: 10,
								Children: [2]*Node{
									{Value: 9},
									{Value: 11},
								},
							},
						},
					},
					{Value: 14},
				},
			},
			NF: func(r *Node, direction Direction) {
				r.Children[0] = r.Children[0].VerticalRotate(direction)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.NF(tt.N, tt.direction)
			if !reflect.DeepEqual(tt.N, tt.want) {
				t.Errorf("VerticalRotate() = %v, want %v", tt.N, tt.want)
			}
		})
	}
}
