package space2trees

import (
	"bytes"
	"fmt"
	"image"
	"reflect"
	"testing"

	"github.com/arran4/spacemap/shared"
	"github.com/google/go-cmp/cmp"
)

func NSMBalanced(shapes ...shared.Shape) func() *Struct {
	return func() *Struct {
		return New().AddAll(shapes...)
	}
}

func NSMUnbalanced(shapes ...shared.Shape) func() *Struct {
	return func() *Struct {
		return New().Unbalance().AddAll(shapes...)
	}
}

func NewTUnbalancedSpaceMap(vTree *Node, hTree *Node) *Struct {
	return &Struct{
		VTree: vTree,
		HTree: hTree,
	}
}

func NewTBalancedSpaceMap(vTree *Node, hTree *Node) *Struct {
	s := &Struct{
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
		Name        string
		Constructor func() *Struct
		Expected    *Struct
	}{
		{
			Name:        "unbalanced rect1",
			Constructor: NSMUnbalanced(rect1),
			Expected: NewTUnbalancedSpaceMap(
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
			Name:        "unbalanced rect1, rect2",
			Constructor: NSMUnbalanced(rect1, rect2),
			Expected: NewTUnbalancedSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 0, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 0, Begin),
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
								NewTHere(rect2, 0, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 0, Begin),
						),
						nil,
						NewTHere(rect1, 0, End)),
					NewTHere(rect1, 0, Begin),
				),
			),
		},
		{
			Name:        "unbalanced rect2, rect3",
			Constructor: NSMUnbalanced(rect2, rect3),
			Expected: NewTUnbalancedSpaceMap(
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
						NewTHere(rect3, 0, End),
					),
					NewTHere(rect2, 0, Begin),
					NewTHere(rect3, 0, Middle),
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
						NewTHere(rect3, 0, End),
					),
					NewTHere(rect2, 0, Begin),
					NewTHere(rect3, 0, Middle),
				),
			),
		},
		{
			Name:        "unbalanced rect1, rect4",
			Constructor: NSMUnbalanced(rect1, rect4),
			Expected: NewTUnbalancedSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(60,
							nil,
							nil,
							NewTHere(rect1, 0, Middle),
							NewTHere(rect4, 0, Begin),
						),
						nil,
						NewTHere(rect1, 0, End),
						NewTHere(rect4, 0, End),
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
							NewTHere(rect4, 0, Begin),
						),
						nil,
						NewTHere(rect1, 0, End),
						NewTHere(rect4, 0, End),
					),
					NewTHere(rect1, 0, Begin),
				),
			),
		},
		{
			Name:        "unbalanced rect1, rect2, rect3",
			Constructor: NSMUnbalanced(rect1, rect2, rect3),
			Expected: NewTUnbalancedSpaceMap(
				NewTNode(10,
					nil,
					NewTNode(100,
						NewTNode(40,
							nil,
							NewTNode(60,
								nil,
								nil,
								NewTHere(rect1, 0, Middle),
								NewTHere(rect2, 0, End),
								NewTHere(rect3, 0, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 0, Begin),
							NewTHere(rect3, 0, Middle),
						),
						nil,
						NewTHere(rect1, 0, End),
					),
					NewTHere(rect1, 0, Begin),
					NewTHere(rect3, 0, Begin),
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
								NewTHere(rect2, 0, End),
								NewTHere(rect3, 0, End),
							),
							NewTHere(rect1, 0, Middle),
							NewTHere(rect2, 0, Begin),
							NewTHere(rect3, 0, Middle),
						),
						nil,
						NewTHere(rect1, 0, End),
					),
					NewTHere(rect1, 0, Begin),
					NewTHere(rect3, 0, Begin),
				),
			),
		},
		{
			Name:        "balanced rect1",
			Constructor: NSMBalanced(rect1),
			Expected: NewTBalancedSpaceMap(
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
			Name:        "balanced rect1, rect2",
			Constructor: NSMBalanced(rect1, rect2),
			Expected:    nil,
		},
		{
			Name:        "balanced rect2, rect3",
			Constructor: NSMBalanced(rect2, rect3),
			Expected:    nil,
		},
		{
			Name:        "balanced rect1, rect4",
			Constructor: NSMBalanced(rect1, rect4),
			Expected:    nil,
		},
		{
			Name:        "balanced rect1, rect2, rect3",
			Constructor: NSMBalanced(rect1, rect2, rect3),
			Expected:    nil,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			sm := test.Constructor()
			if sm.Balanced {
				balancedDepthTest(sm.VTree, 0, t, []int{})
				balancedDepthTest(sm.HTree, 0, t, []int{})
			}
			if test.Expected != nil {
				if s := cmp.Diff(sm, test.Expected); len(s) > 0 {
					t.Errorf("Failed stacks differ: %s", s)
				}
			}
			if t.Failed() {
				value := func(n *Node) string {
					return fmt.Sprintf("%d", n.Value)
				}
				depth := func(n *Node) string {
					return fmt.Sprintf("%d", n.MaxDepth)
				}
				t.Logf("Result VTree:\n%s", plotTree(sm.VTree, 3, value))
				t.Logf("Result HTree:\n%s", plotTree(sm.HTree, 3, value))
				t.Logf("Result VTree.depth:\n%s", plotTree(sm.VTree, 3, depth))
				t.Logf("Result HTree.depth:\n%s", plotTree(sm.HTree, 3, depth))
				if test.Expected != nil {
					t.Logf("Expected VTree:\n%s", plotTree(test.Expected.VTree, 3, value))
					t.Logf("Expected HTree:\n%s", plotTree(test.Expected.HTree, 3, value))
				}
			}
		})
	}
}

func balancedDepthTest(n *Node, i int, t *testing.T, p []int) int {
	if n == nil {
		return i - 1
	}
	p = append(p, n.Value)
	lr := balancedDepthTest(n.Children[0], i+1, t, p)
	rr := balancedDepthTest(n.Children[1], i+1, t, p)
	b := Balance(rr - lr)
	if b.Extreme() {
		t.Errorf("%#v is unbalanced: %d vs %d", p, lr, rr)
	}
	m := lr
	if lr < rr {
		m = rr
	}
	if n.MaxDepth != m {
		t.Errorf("%#v is depth is wrong: .MaxDepth= %d vs max depth= %d", p, n.MaxDepth, m)
	}
	return m
}

func plotTree(n *Node, w int, value func(*Node) string) string {
	nextLine := []*Node{n}
	var p [][]string
	for c := 1; c > 0; {
		c = 0
		line := nextLine
		nextLine = make([]*Node, len(line)*2)
		strl := len(line)*2 - 1
		for pli := range p {
			ns := make([]string, strl)
			for i, v := range p[pli] {
				ns[i*2+1] = v
			}
			p[pli] = ns
		}
		sl := make([]string, strl)
		for i, en := range line {
			if en == nil {
				continue
			}
			sl[i*2] = value(en)
			if en.Children[0] != nil {
				c++
				nextLine[i*2] = en.Children[0]
			}
			if en.Children[1] != nil {
				c++
				nextLine[i*2+1] = en.Children[1]
			}
		}
		p = append(p, sl)
	}
	b := &bytes.Buffer{}
	for _, pl := range p {
		fmt.Fprintf(b, fmt.Sprintf("%%%ds ", w), pl)
		b.WriteString("\n")
	}
	return b.String()
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
		SpaceMap func() *Struct
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

func TestNode_AvlBalance(t *testing.T) {
	tests := []struct {
		name string
		node *Node
		want *Node
	}{
		{
			name: "Don't disorder",
			node: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
			want: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
		},
		{
			name: "Right right",
			node: &Node{
				Value:    10,
				MaxDepth: 2,
				Children: [2]*Node{
					nil,
					{
						Value:    40,
						MaxDepth: 2,
						Children: [2]*Node{
							nil,
							{
								Value:    60,
								MaxDepth: 2,
								Children: [2]*Node{
									nil,
									nil,
								},
							},
						},
					},
				},
			},
			want: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
		},
		{
			name: "Right left",
			node: &Node{
				Value:    10,
				MaxDepth: 2,
				Children: [2]*Node{
					nil,
					{
						Value:    60,
						MaxDepth: 2,
						Children: [2]*Node{
							{
								Value:    40,
								MaxDepth: 2,
								Children: [2]*Node{
									nil,
									nil,
								},
							},
							nil,
						},
					},
				},
			},
			want: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
		},
		{
			name: "Left left",
			node: &Node{
				Value:    60,
				MaxDepth: 2,
				Children: [2]*Node{
					{
						Value:    40,
						MaxDepth: 2,
						Children: [2]*Node{
							{
								Value:    10,
								MaxDepth: 2,
								Children: [2]*Node{
									nil,
									nil,
								},
							},
							nil,
						},
					},
					nil,
				},
			},
			want: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
		},
		{
			name: "Left right",
			node: &Node{
				Value:    60,
				MaxDepth: 2,
				Children: [2]*Node{
					{
						Value:    40,
						MaxDepth: 2,
						Children: [2]*Node{
							{
								Value:    10,
								MaxDepth: 2,
								Children: [2]*Node{
									nil,
									nil,
								},
							},
							nil,
						},
					},
					nil,
				},
			},
			want: &Node{
				Value: 40, MaxDepth: 1,
				Children: [2]*Node{
					{Value: 10, MaxDepth: 1},
					{Value: 60, MaxDepth: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.AvlBalance(0)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("AvlBalance() = \n%s", diff)
			}
			if t.Failed() {
				value := func(n *Node) string {
					return fmt.Sprintf("%d", n.Value)
				}
				t.Logf("Result VTree:\n%s", plotTree(got, 3, value))
			}
		})
	}
}

func TestStruct_Remove(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(40, 40, 60, 60, shared.Name("rect2"))
	rect3 := shared.NewRectangle(10, 10, 60, 60, shared.Name("rect3"))
	rect4 := shared.NewRectangle(60, 60, 100, 100, shared.Name("rect4"))
	tests := []struct {
		name         string
		Constructor  func() *Struct
		Expected     func() *Struct
		shape        shared.Shape
		ScanExpected *int
		ScanFor      shared.Shape
	}{
		{
			name:        "Remove the only unbalanced",
			Constructor: NSMUnbalanced(rect1),
			shape:       rect1,
			Expected:    NSMUnbalanced(),
		},
		{
			name:        "Remove the only balanced",
			Constructor: NSMBalanced(rect1),
			shape:       rect1,
			Expected:    NSMBalanced(),
		},
		{
			name:        "Remove the first of 2 unbalanced",
			Constructor: NSMUnbalanced(rect1, rect2),
			shape:       rect1,
			Expected: func() *Struct {
				return NewTUnbalancedSpaceMap(
					NewTNode(40,
						nil,
						NewTNode(60,
							nil,
							nil,
							NewTHere(rect2, 0, End),
						),
						NewTHere(rect2, 0, Begin),
					),
					NewTNode(40,
						nil,
						NewTNode(60,
							nil,
							nil,
							NewTHere(rect2, 0, End),
						),
						NewTHere(rect2, 0, Begin),
					),
				)
			},
		},
		{
			name:        "Remove the first of 2 balanced",
			Constructor: NSMBalanced(rect1, rect2),
			shape:       rect1,
			Expected: func() *Struct {
				return NewTBalancedSpaceMap(
					NewTNode(60,
						NewTNode(40,
							nil,
							nil,
							NewTHere(rect2, 0, Begin),
						),
						nil,
						NewTHere(rect2, 0, End),
					),
					NewTNode(60,
						NewTNode(40,
							nil,
							nil,
							NewTHere(rect2, 0, Begin),
						),
						nil,
						NewTHere(rect2, 0, End),
					),
				)
			},
		},
		{
			name:         "Remove the first of 4 balanced",
			Constructor:  NSMBalanced(rect1, rect2, rect3, rect4),
			shape:        rect1,
			ScanFor:      rect1,
			ScanExpected: PI(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := tt.Constructor()
			sm.Remove(tt.shape)
			if sm.Balanced {
				balancedDepthTest(sm.VTree, 0, t, []int{})
				balancedDepthTest(sm.HTree, 0, t, []int{})
			}
			var expected *Struct = nil
			if tt.Expected != nil {
				expected = tt.Expected()
				if s := cmp.Diff(sm, expected); len(s) > 0 {
					t.Errorf("Failed stacks differ: %s", s)
				}
			}
			if tt.ScanExpected != nil {
				if found := scanFor(sm, tt.ScanFor); found != *tt.ScanExpected {
					t.Errorf("Failed scan for found %d when expecting %d of: %s", found, *tt.ScanExpected, tt.ScanFor)
				}
			}
			if t.Failed() {
				value := func(n *Node) string {
					return fmt.Sprintf("%d", n.Value)
				}
				depth := func(n *Node) string {
					return fmt.Sprintf("%d", n.MaxDepth)
				}
				t.Logf("Result VTree:\n%s", plotTree(sm.VTree, 3, value))
				t.Logf("Result HTree:\n%s", plotTree(sm.HTree, 3, value))
				t.Logf("Result VTree.depth:\n%s", plotTree(sm.VTree, 3, depth))
				t.Logf("Result HTree.depth:\n%s", plotTree(sm.HTree, 3, depth))
				if tt.Expected != nil {
					t.Logf("Expected VTree:\n%s", plotTree(expected.VTree, 3, value))
					t.Logf("Expected HTree:\n%s", plotTree(expected.HTree, 3, value))
				}
			}
		})
	}
}

func scanFor(sm *Struct, scanFor shared.Shape) int {
	r := 0
	r += scanForInNode(sm.VTree, scanFor)
	r += scanForInNode(sm.HTree, scanFor)
	return r
}

func scanForInNode(tree *Node, shape shared.Shape) int {
	if tree == nil {
		return 0
	}
	r := 0
	for _, e := range tree.Here {
		if e.Shape == shape {
			r++
		}
	}
	r += scanForInNode(tree.Children[0], shape)
	r += scanForInNode(tree.Children[1], shape)
	return r
}

func PI(i int) *int {
	return &i
}

func TestZIndex(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect2"))
	tests := []struct {
		name        string
		Constructor func() *Struct
		want        shared.Shape
	}{
		{
			name: "r1&2 want r2",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect1, 1)
				s.Add(rect2, 0)
				return s
			},
			want: rect2,
		},
		{
			name: "r1&2 want r1",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect1, 1)
				s.Add(rect2, 2)
				return s
			},
			want: rect1,
		},
		{
			name: "r2&1 want r2",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect2, 0)
				s.Add(rect1, 1)
				return s
			},
			want: rect2,
		},
		{
			name: "r2&1 want r1",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect2, 2)
				s.Add(rect1, 1)
				return s
			},
			want: rect1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			sm := test.Constructor()
			r := sm.GetAt(50, 50)
			if r != test.want {
				tt.Errorf("Failed got %s expected %s", r, test.want)
			}
		})
	}
}

func TestNode_InsertHere(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect2"))
	tests := []struct {
		name string
		*Node
		ZIndex   *int
		Shape    shared.Shape
		Type     Type
		Expected []*Here
	}{
		{
			name: "Insert into empty okay",
			Node: &Node{
				Here: []*Here{},
			},
			ZIndex: PI(0),
			Shape:  rect1,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect1, ZIndex: 0},
			},
		},
		{
			name: "Insert into duplicate zid",
			Node: &Node{
				Here: []*Here{
					&Here{Shape: rect1, ZIndex: 0},
				},
			},
			ZIndex: PI(0),
			Shape:  rect2,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect1, ZIndex: 0},
				&Here{Shape: rect2, ZIndex: 0},
			},
		},
		{
			name: "Insert into duplicate zid reversed",
			Node: &Node{
				Here: []*Here{
					&Here{Shape: rect2, ZIndex: 0},
				},
			},
			ZIndex: PI(0),
			Shape:  rect1,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect2, ZIndex: 0},
				&Here{Shape: rect1, ZIndex: 0},
			},
		},
		{
			name: "Insert into one lower",
			Node: &Node{
				Here: []*Here{
					&Here{Shape: rect2, ZIndex: -1},
				},
			},
			ZIndex: PI(0),
			Shape:  rect1,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect2, ZIndex: -1},
				&Here{Shape: rect1, ZIndex: 0},
			},
		},
		{
			name: "Insert into one greater",
			Node: &Node{
				Here: []*Here{
					&Here{Shape: rect2, ZIndex: 1},
				},
			},
			ZIndex: PI(0),
			Shape:  rect1,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect1, ZIndex: 0},
				&Here{Shape: rect2, ZIndex: 1},
			},
		},
		{
			name: "Insert into position 5 of 10",
			Node: &Node{
				Here: []*Here{
					&Here{Shape: rect2, ZIndex: 1},
					&Here{Shape: rect2, ZIndex: 2},
					&Here{Shape: rect2, ZIndex: 3},
					&Here{Shape: rect2, ZIndex: 4},
					&Here{Shape: rect2, ZIndex: 6},
					&Here{Shape: rect2, ZIndex: 7},
					&Here{Shape: rect2, ZIndex: 8},
					&Here{Shape: rect2, ZIndex: 9},
					&Here{Shape: rect2, ZIndex: 10},
				},
			},
			ZIndex: PI(5),
			Shape:  rect1,
			Type:   0,
			Expected: []*Here{
				&Here{Shape: rect2, ZIndex: 1},
				&Here{Shape: rect2, ZIndex: 2},
				&Here{Shape: rect2, ZIndex: 3},
				&Here{Shape: rect2, ZIndex: 4},
				&Here{Shape: rect1, ZIndex: 5},
				&Here{Shape: rect2, ZIndex: 6},
				&Here{Shape: rect2, ZIndex: 7},
				&Here{Shape: rect2, ZIndex: 8},
				&Here{Shape: rect2, ZIndex: 9},
				&Here{Shape: rect2, ZIndex: 10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.InsertHere(tt.ZIndex, tt.Shape, tt.Type)
			if diff := cmp.Diff(tt.Here, tt.Expected); diff != "" {
				t.Errorf("Error mismatch with here and expected;\n%s", diff)
			}
		})
	}
}
