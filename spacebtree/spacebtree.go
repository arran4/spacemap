package spacebtree

import (
	"fmt"
	"spacemap/shared"
)

type Type int

const (
	Middle Type = iota
	Begin  Type = iota
	End    Type = iota
)

func (t Type) String() string {
	switch t {
	case Middle:
		return "Middle"
	case Begin:
		return "Begin"
	case End:
		return "End"
	}
	return "Unknown"
}

type Here struct {
	Shape  shared.Shape
	ZIndex int
	Type   Type
}

func (h *Here) String() string {
	return fmt.Sprintf("Here(%s, Z:%d, Type:%s)", h.Shape, h.ZIndex, h.Type)
}

func (h *Here) Copy() *Here {
	return &Here{
		Shape:  h.Shape,
		ZIndex: h.ZIndex,
		Type:   h.Type,
	}
}

var _ fmt.Stringer = (*Here)(nil)

type Node struct {
	Value    int
	Depth    int
	Here     []*Here
	Children [2]*Node
}

func (n *Node) AddBetween(from, to int, s shared.Shape, zIndex *int, leftMost, rightMost bool, parent *Node, balance int) *Node {

	if n == nil {
		var r *Node
		if leftMost || rightMost {
			var lm *Node
			if leftMost {
				r = NewNode(from, s, Begin, zIndex, parent, balance)
				lm = r
			}
			if rightMost {
				if lm != nil && balance >= 0 {
					balance++
				}
				r = NewNode(to, s, End, zIndex, parent, balance)
				if lm != nil {
					lm.Children[1] = r
					lm.Depth = r.Depth
					r = lm
				}
			}
		}
		return r
	}

	if n.Value == from && leftMost {
		n.InsertHere(zIndex, s, Begin)
		leftMost = false
	} else if n.Value == to && rightMost {
		n.InsertHere(zIndex, s, End)
		rightMost = false
	} else if from < n.Value && n.Value < to {
		n.InsertHere(zIndex, s, Middle)
	}
	if balance >= 0 {
		balance++
	}
	if n.Value > from {
		n.Children[0] = n.Children[0].AddBetween(from, to, s, zIndex, leftMost, rightMost && n.Value >= to, n, balance)
	}
	if n.Value < to {
		n.Children[1] = n.Children[1].AddBetween(from, to, s, zIndex, leftMost && n.Value <= from, rightMost, n, balance)
	}
	r := n
	for balance >= 0 {
		var ldepth int
		var rdepth int
		if r.Children[0] != nil {
			ldepth = r.Children[0].Depth
		}
		if r.Children[1] != nil {
			rdepth = r.Children[1].Depth
		}
		n.Depth = ldepth
		if n.Depth < rdepth {
			n.Depth = rdepth
		}
		if !(ldepth-rdepth > 1 || ldepth-rdepth < 0) {
			break
		}
		rn := r.Rotate(ldepth - rdepth)
		if rn == nil {
			break
		}
		r = rn
	}
	return r
}

func (n *Node) InsertHere(zIndex *int, s shared.Shape, t Type) {
	var zi int
	if n != nil {
		zi = len(n.Here)
	}
	if zIndex != nil {
		zi = *zIndex
	}
	r := 0
	for ; r < len(n.Here) && n.Here[r].ZIndex < zi; r++ {
	}
	n.Here = append(n.Here, nil)
	if r < len(n.Here) {
		copy(n.Here[r+1:], n.Here[r:len(n.Here)-1])
		n.Here[r] = &Here{
			Shape:  s,
			ZIndex: zi,
			Type:   t,
		}
	}
}

func (n *Node) Get(v int) (result []shared.Shape) {
	if n == nil {
		return
	}
	if n.Value > v {
		result = n.Children[0].Get(v)
	} else if n.Value < v {
		result = n.Children[1].Get(v)
	} else {
		for _, e := range n.Here {
			result = append(result, e.Shape)
		}
		return
	}
	if result == nil {
		for _, e := range n.Here {
			if n.Value >= v && e.Type != Begin {
				result = append(result, e.Shape)
			} else if n.Value <= v && e.Type != End {
				result = append(result, e.Shape)
			}
		}
	}
	return
}

func (n *Node) Rotate(direction int) *Node {
	c, sb := 0, 1
	if direction > 0 {
		c, sb = sb, c
	}
	if n.Children[c] == nil {
		return nil
	}
	depth := n.Depth
	var r *Node
	r, n.Children[c], n.Children[c].Children[sb] = n.Children[c], n.Children[c].Children[sb], n
	r.RecalculateDepth(depth)
	return r
}

func (n *Node) RecalculateDepth(depth int) int {
	if n == nil {
		return depth
	}
	ld := n.Children[0].RecalculateDepth(depth + 1)
	rd := n.Children[1].RecalculateDepth(depth + 1)
	n.Depth = ld
	if n.Depth < rd {
		n.Depth = rd
	}
	return n.Depth
}

func NewNode(p int, s shared.Shape, hType Type, zIndex *int, parent *Node, depth int) *Node {
	var here []*Here
	if parent != nil {
		for _, ph := range parent.Here {
			if ph.Shape == s {
				continue
			}
			h := ph.Copy()
			h.Type = Middle
			switch ph.Type {
			case Middle:
				here = append(here, h)
			case End:
				if p < parent.Value {
					here = append(here, h)
				}
			case Begin:
				if p > parent.Value {
					here = append(here, h)
				}
			}
		}
	}
	nn := &Node{
		Value: p,
		Here:  here,
	}
	if depth >= 0 {
		nn.Depth = depth + 1
	}
	nn.InsertHere(zIndex, s, hType)
	return nn
}

type SpaceMap struct {
	VTree    *Node
	HTree    *Node
	Balanced bool
}

func (m *SpaceMap) AddAll(shapes ...shared.Shape) *SpaceMap {
	for _, shape := range shapes {
		m.Add(shape)
	}
	return m
}

func (m *SpaceMap) Add(shape shared.Shape) *SpaceMap {
	b := shape.Bounds()
	var balance = -1
	if m.Balanced {
		balance = 0
	}
	m.VTree = m.VTree.AddBetween(b.Min.Y, b.Max.Y, shape, nil, true, true, nil, balance)
	m.HTree = m.HTree.AddBetween(b.Min.X, b.Max.X, shape, nil, true, true, nil, balance)
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

func (m *SpaceMap) Unbalance() *SpaceMap {
	m.Balanced = false
	return m
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{
		Balanced: true,
	}
}
