package space2trees

import (
	"strconv"

	"github.com/arran4/spacemap/shared"
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
	return "Here(" + h.Shape.String() + ", Z:" + strconv.Itoa(h.ZIndex) + ", Type:" + h.Type.String() + ")"
}

func (h *Here) Copy() *Here {
	return &Here{
		Shape:  h.Shape,
		ZIndex: h.ZIndex,
		Type:   h.Type,
	}
}

// Ensure Here implements the Stringer interface without importing fmt.
type stringer interface{ String() string }

var _ stringer = (*Here)(nil)

type Node struct {
	Value    int
	MaxDepth int
	Here     []*Here
	Children [2]*Node
}

func (n *Node) AddBetween(from, to int, s shared.Shape, zIndex *int, leftMost, rightMost bool, parent *Node, depth int) *Node {

	if n == nil {
		var r *Node
		if leftMost || rightMost {
			if leftMost {
				r = NewNode(from, s, Begin, zIndex, parent, depth)
				if rightMost {
					var nDepth = depth
					r.Children[1] = r.Children[1].AddBetween(from, to, s, zIndex, false, rightMost, r, nDepth)
				}
			} else if rightMost {
				r = NewNode(to, s, End, zIndex, parent, depth)
			}
		}
		if depth >= 0 {
			r = r.AvlBalance(depth)
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
	var nDepth = depth
	if depth >= 0 {
		nDepth = depth + 1
	}
	if n.Value > from {
		n.Children[0] = n.Children[0].AddBetween(from, to, s, zIndex, leftMost, rightMost && n.Value >= to, n, nDepth)
	}
	if n.Value < to {
		n.Children[1] = n.Children[1].AddBetween(from, to, s, zIndex, leftMost && n.Value <= from, rightMost, n, nDepth)
	}
	r := n
	if depth >= 0 {
		r = r.AvlBalance(depth)
	}
	return r
}

func (n *Node) InsertHere(zIndex *int, s shared.Shape, t Type) {
	var zi int
	if zIndex != nil {
		zi = *zIndex
	}
	h := &Here{
		Shape:  s,
		ZIndex: zi,
		Type:   t,
	}
	p := 0
	for ; p < len(n.Here) && h.ZIndex >= n.Here[p].ZIndex; p++ {
	}
	n.Here = append(n.Here, nil)
	copy(n.Here[p+1:], n.Here[p:])
	n.Here[p] = h
}

func (n *Node) InsertHereReZIndex(zIndex *int, s shared.Shape, t Type) {
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

func (n *Node) AvlBalance(depth int) *Node {
	if n == nil {
		return nil
	}
	_ = n.RecalculateDepth(depth)
	nBal := n.Balance()
	var cBal = Balanced
	switch nBal {
	case Left:
		cBal = n.Children[0].Balance()
	case Right:
		cBal = n.Children[1].Balance()
	default:
		return n
	}
	if !cBal.SameSide(nBal) {
		d := 0
		if cBal.Left() {
			d = 1
		}
		n.Children[d] = n.Children[d].VerticalRotate(Direction(cBal))
	}
	r := n.VerticalRotate(Direction(nBal))
	_ = r.RecalculateDepth(depth)
	return r
}

type Direction int
type Balance int

func (b Balance) Left() bool {
	return b < 0
}

func (b Balance) Right() bool {
	return b > 0
}

func (b Balance) SameSide(b2 Balance) bool {
	if b < 0 && b2 < 0 {
		return true
	}
	if b > 0 && b2 > 0 {
		return true
	}
	return false
}

func (b Balance) Extreme() bool {
	switch b {
	case Left, Right:
		return true
	default:
		return false
	}
}

const (
	Left      Balance = -2
	LeftLean  Balance = -1
	Balanced  Balance = 0
	RightLean Balance = 1
	Right     Balance = 2
)

func (n *Node) Balance() Balance {
	var lDepth int
	var rDepth int
	if n.Children[0] != nil {
		lDepth = n.Children[0].MaxDepth
	}
	if n.Children[1] != nil {
		rDepth = n.Children[1].MaxDepth
	}
	nb := rDepth - lDepth
	switch nb {
	case 0:
		return Balanced
	case -1:
		return LeftLean
	case 1:
		return RightLean
	default:
		if nb < 0 {
			return Left
		} else {
			return Right
		}
	}
}

func (n *Node) VerticalRotate(direction Direction) *Node {
	c, cs := 0, 1
	if direction > 0 {
		c, cs = cs, c
	}
	var r *Node
	if n == nil || n.Children[c] == nil {
		return n
	}
	r, n.Children[c], n.Children[c].Children[cs] = n.Children[c], n.Children[c].Children[cs], n
	return r
}

func (n *Node) RecalculateDepth(depth int) int {
	if n == nil {
		return depth - 1
	}
	var ld int
	var rd int
	if depth >= 0 {
		ld = n.Children[0].RecalculateDepth(depth + 1)
		rd = n.Children[1].RecalculateDepth(depth + 1)
	} else {
		if n.Children[0] != nil {
			ld = n.Children[0].MaxDepth
		}
		if n.Children[1] != nil {
			rd = n.Children[1].MaxDepth
		}
	}
	n.MaxDepth = ld
	if n.MaxDepth < rd {
		n.MaxDepth = rd
	}
	return n.MaxDepth
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
		nn.MaxDepth = depth + 1
	}
	nn.InsertHere(zIndex, s, hType)
	return nn
}

type Struct struct {
	VTree    *Node
	HTree    *Node
	Balanced bool
}

func (m *Struct) AddAll(shapes ...shared.Shape) *Struct {
	for _, shape := range shapes {
		m.Add(shape, 0)
	}
	return m
}

func (m *Struct) Add(shape shared.Shape, zIndex int) {
	b := shape.Bounds()
	var balance = -1
	if m.Balanced {
		balance = 0
	}
	m.VTree = m.VTree.AddBetween(b.Min.Y, b.Max.Y, shape, &zIndex, true, true, nil, balance)
	m.HTree = m.HTree.AddBetween(b.Min.X, b.Max.X, shape, &zIndex, true, true, nil, balance)
}

func (m *Struct) GetStackAt(x int, y int) []shared.Shape {
	xs := m.HTree.Get(x)
	if len(xs) == 0 {
		return []shared.Shape{}
	}
	ys := m.VTree.Get(y)
	if len(ys) == 0 {
		return []shared.Shape{}
	}

	if len(xs)*len(ys) < 64 {
		result := make([]shared.Shape, 0, len(ys))
		for _, yVal := range ys {
			for _, xVal := range xs {
				if xVal == yVal {
					result = append(result, yVal)
					break
				}
			}
		}
		return result
	}

	seen := make(map[shared.Shape]struct{}, len(xs))
	for _, e := range xs {
		seen[e] = struct{}{}
	}
	result := make([]shared.Shape, 0, len(ys))
	for _, e := range ys {
		if _, ok := seen[e]; ok {
			result = append(result, e)
		}
	}
	return result
}

func (m *Struct) Unbalance() *Struct {
	m.Balanced = false
	return m
}

func New() *Struct {
	return &Struct{
		Balanced: true,
	}
}

func (n *Node) RemoveBetween(from, to int, s shared.Shape, depth int) (*Node, bool) {
	if n == nil {
		return n, false
	}
	r := n
	deleted := false
	if from <= r.Value && r.Value <= to {
		for {
			remaining := r.RemoveHere(s)
			if remaining == 0 {
				r = r.DeleteNode()
				deleted = true
				if r == nil {
					return nil, deleted
				}
				continue
			}
			break
		}
	}
	var nDepth = depth
	if depth >= 0 {
		nDepth = depth + 1
	}
	if r.Value > from {
		var d bool
		n.Children[0], d = r.Children[0].RemoveBetween(from, to, s, nDepth)
		deleted = deleted || d
	}
	if r.Value < to {
		var d bool
		r.Children[1], d = r.Children[1].RemoveBetween(from, to, s, nDepth)
		deleted = deleted || d
	}
	if depth >= 0 && deleted {
		r = r.AvlBalance(depth)
	}
	return r, deleted
}

func (n *Node) RemoveHere(s shared.Shape) int {
	shrink := 0
	for i := range n.Here {
		for i+shrink < len(n.Here) {
			if n.Here[i+shrink].Shape == s {
				shrink++
			} else {
				break
			}
		}
		if i+shrink >= len(n.Here) {
			break
		}
		if shrink == 0 {
			continue
		}
		n.Here[i] = n.Here[i+shrink]
	}
	n.Here = n.Here[:len(n.Here)-shrink]
	return len(n.Here)
}

func (n *Node) DeleteNode() *Node {
	if n.Children[0] == nil {
		return n.Children[1]
	}
	if n.Children[1] == nil {
		return n.Children[0]
	}
	nn, np := n.Children[1].Most(0)
	if np != nil {
		np.Children[0] = nn.Children[1]
	} else {
		n.Children[1] = nil
	}
	nn.Children = n.Children
	return nn
}

func (n *Node) Most(i int) (*Node, *Node) {
	if n == nil {
		return nil, nil
	}
	m, p := n.Children[i].Most(i)
	if m == nil {
		return n, nil
	}
	if p == nil {
		p = n
	}
	return m, p
}

func (m *Struct) Remove(shape shared.Shape) {
	b := shape.Bounds()
	var balance = -1
	if m.Balanced {
		balance = 0
	}
	m.VTree, _ = m.VTree.RemoveBetween(b.Min.Y, b.Max.Y, shape, balance)
	m.HTree, _ = m.HTree.RemoveBetween(b.Min.X, b.Max.X, shape, balance)
}

func (m *Struct) GetAt(x int, y int) shared.Shape {
	s := m.GetStackAt(x, y)
	if len(s) > 0 {
		return s[0]
	}
	return nil
}
