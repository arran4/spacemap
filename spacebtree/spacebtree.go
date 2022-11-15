package spacebtree

import (
	"spacemap/shared"
)

type Type int

const (
	Middle Type = iota
	Begin  Type = iota
	End    Type = iota
)

type Here struct {
	Shape  shared.Shape
	ZIndex int
	Type   Type
}

type Node struct {
	Value    int
	Here     []*Here
	Children [2]*Node
}

func (n *Node) AddBetween(from, to int, s shared.Shape, zIndex *int, leftMost, rightMost bool) *Node {
	var zi int
	if n != nil {
		zi = len(n.Here)
	}
	if zIndex != nil {
		zi = *zIndex
	}

	if n == nil && (leftMost || rightMost) {
		r := n
		var lm *Node
		if leftMost {
			r = NewNode(from, s, Begin, zi)
			lm = r
			leftMost = false
		}
		if rightMost {
			r = NewNode(to, s, End, zi)
			if lm != nil {
				lm.Children[1] = r
				r = lm
			}
			rightMost = false
		}
		return r
	}

	if n.Value == from && rightMost {
		n.InsertHere(zi, s, Begin)
		leftMost = false
	} else if n.Value == to && leftMost {
		n.InsertHere(zi, s, End)
		rightMost = false
	} else if from < n.Value && n.Value < to {
		n.InsertHere(zi, s, Middle)
	}

	if n.Value > from {
		n.Children[0] = n.Children[0].AddBetween(from, to, s, zIndex, leftMost, rightMost)
	} else {
		n.Children[1] = n.Children[0].AddBetween(from, to, s, zIndex, leftMost, rightMost)
	}
	return n
}

func (n *Node) InsertHere(zi int, s shared.Shape, t Type) {
	r := 0
	for ; r < len(n.Here); r++ {
		if zi >= n.Here[r].ZIndex {
			break
		}
	}
	h := &Here{
		Shape:  s,
		ZIndex: zi,
		Type:   t,
	}
	n.Here = append(n.Here, h)
	if r <= len(n.Here) {
		copy(n.Here[r:], n.Here[r+1:])
		n.Here[r] = h
	}
}

func (n *Node) Get(v int) (result []shared.Shape) {
	if n == nil {
		return
	}
	if n.Value > v {
		result = n.Children[0].Get(v)
	} else if n.Value < v {
		result = n.Children[0].Get(v)
	}
	if result == nil {
		for _, e := range n.Here {
			if n.Value >= v {
				result = append(result, e.Shape)
			}
		}
	}
	return
}

func NewNode(p int, s shared.Shape, hType Type, zi int) *Node {
	nn := &Node{
		Value: p,
		Here: []*Here{
			{
				Shape:  s,
				ZIndex: zi,
				Type:   hType,
			},
		},
	}
	return nn
}

type SpaceMap struct {
	VTree *Node
	HTree *Node
}

func (m *SpaceMap) AddAll(shapes ...shared.Shape) *SpaceMap {
	for _, shape := range shapes {
		m.Add(shape)
	}
	return m
}

func (m *SpaceMap) Add(shape shared.Shape) *SpaceMap {
	b := shape.Bounds()
	m.VTree = m.VTree.AddBetween(b.Min.Y, b.Max.Y, shape, nil, true, true)
	m.HTree = m.HTree.AddBetween(b.Min.X, b.Max.X, shape, nil, true, true)
	return m
}

func (m *SpaceMap) GetStackAt(x int, y int) []shared.Shape {
	xs := m.HTree.Get(x)
	if len(xs) == 0 {
		return []shared.Shape{}
	}
	seen := map[shared.Shape]struct{}{}
	for _, e := range xs {
		seen[e] = struct{}{}
	}
	ys := m.VTree.Get(y)
	if len(ys) == 0 {
		return []shared.Shape{}
	}
	result := make([]shared.Shape, 0)
	for _, e := range ys {
		if _, ok := seen[e]; ok {
			result = append(result, e)
		}
	}
	return result
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{}
}
