package btree

import (
	"image"
	"spacemap/shared"
)

type Type int

const (
	Min Type = iota
	Max Type = iota
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

func (n *Node) Add(p int, s shared.Shape, hType Type, zIndex *int) *Node {
	if n == nil {
		return NewNode(p, s, hType)
	}
	if n.Value == p {
		zi := len(n.Here)
		if zIndex != nil {
			zi = *zIndex
		}
		HereInsert(n.Here, zi)
		h := &Here{
			Shape:  s,
			ZIndex: zi,
			Type:   hType,
		}
		n.Here = append(n.Here, h)
	}
	if n.Value > p {
		n.Children[0] = n.Children[0].Add(p, s, hType, zIndex)
	} else {
		n.Children[1] = n.Children[0].Add(p, s, hType, zIndex)
	}
	return n
}

func HereInsert(heres []*Here, zi int) (r int) {
	prev := 0
	r = 0
	for i, h := range heres {
		if h.ZIndex < zi {
			r = i + 1
			continue
		}
		if prev == 0 {
			r = i
			prev = zi
		}
		if h.ZIndex == prev {
			h.ZIndex++
		}
		prev = h.ZIndex
	}
	return
}

func NewNode(p int, s shared.Shape, hType Type) *Node {
	nn := &Node{
		Value: p,
		Here: []*Here{
			{
				Shape:  s,
				ZIndex: 0,
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
	m.VTree = m.VTree.Add(b.Min.Y, shape, Min, nil)
	m.VTree = m.VTree.Add(b.Max.Y, shape, Max, nil)
	m.HTree = m.HTree.Add(b.Min.X, shape, Min, nil)
	m.HTree = m.HTree.Add(b.Max.X, shape, Max, nil)
	return m
}

func (m *SpaceMap) GetXYPositions(p image.Point) (int, int) {
	panic("not implmented todo")
}

func (m *SpaceMap) GetStackAt(x int, y int) []shared.Shape {
	panic("not implmented todo")
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{}
}
