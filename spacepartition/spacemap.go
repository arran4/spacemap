package spaceparition

import (
	"image"
	"log"
	"sort"
	"spacemap/shared"
)

type Alignment int

const (
	NotAligned Alignment = iota
	Vertical
	Horizontal
)

type Split struct {
	Position  int
	BecauseOf []shared.Shape
	Alignment Alignment
}

type SplitCoordination struct {
	HSplit *Split
	VSplit *Split
}

func (c SplitCoordination) Sane() bool {
	if c.VSplit.Alignment != Vertical {
		return false
	}
	if c.HSplit.Alignment != Horizontal {
		return false
	}
	return true
}

func SC(HSplit *Split, VSplit *Split) SplitCoordination {
	coordination := SplitCoordination{
		HSplit,
		VSplit,
	}
	return coordination
}

type Struct struct {
	VSplits []*Split
	HSplits []*Split
	Stacks  map[SplitCoordination][]shared.Shape
}

func (m *Struct) AddAll(shapes ...shared.Shape) *Struct {
	for _, shape := range shapes {
		m.Add(shape)
	}
	return m
}

func (m *Struct) Add(shape shared.Shape) *Struct {
	b := shape.Bounds()
	minxi, minyi := m.GetXYPositions(b.Min)
	maxxi, maxyi := m.GetXYPositions(b.Max)
	{
		var maxxhs [2]*Split
		if maxxi == len(m.HSplits) || m.HSplits[maxxi].Position != b.Max.X {
			maxxhs[0] = &Split{
				Position:  b.Max.X,
				BecauseOf: []shared.Shape{shape},
				Alignment: Horizontal,
			}
			if maxxi >= 0 && maxxi < len(m.HSplits) {
				maxxhs[1] = m.HSplits[maxxi]
			}
			m.HSplits = append(m.HSplits[:maxxi], append([]*Split{maxxhs[0]}, m.HSplits[maxxi:]...)...)
		} else {
			m.HSplits[maxxi].BecauseOf = append(m.HSplits[maxxi].BecauseOf, shape)
		}
		if maxyi == len(m.VSplits) || m.VSplits[maxyi].Position != b.Max.Y {
			var ovs *Split
			vs := &Split{
				Position:  b.Max.Y,
				BecauseOf: []shared.Shape{shape},
				Alignment: Vertical,
			}
			if maxyi >= 0 && maxyi < len(m.VSplits) {
				ovs = m.VSplits[maxyi]
			}
			m.VSplits = append(m.VSplits[:maxyi], append([]*Split{vs}, m.VSplits[maxyi:]...)...)
			var lhs *Split
			for _, hs := range m.HSplits {
				_, exists := m.Stacks[SC(hs, vs)]
				if exists {
					lhs = hs
					continue
				}
				if ovs != nil {
					if hs != maxxhs[0] {
						m.Stacks[SC(hs, vs)] = append([]shared.Shape{}, m.Stacks[SC(hs, ovs)]...)
					} else if lhs != nil {
						m.Stacks[SC(hs, vs)] = append([]shared.Shape{}, m.Stacks[SC(lhs, ovs)]...)
					}
				}
				lhs = hs
			}
		} else {
			m.VSplits[maxyi].BecauseOf = append(m.VSplits[maxyi].BecauseOf, shape)
		}
		if maxxhs[0] != nil {
			for _, vs := range m.VSplits {
				_, exists := m.Stacks[SC(maxxhs[0], vs)]
				if exists {
					continue
				}
				if maxxhs[1] != nil {
					m.Stacks[SC(maxxhs[0], vs)] = append([]shared.Shape{}, m.Stacks[SC(maxxhs[1], vs)]...)
				}
			}
		}
	}
	{
		var minxhs [2]*Split
		if m.HSplits[minxi].Position != b.Min.X {
			minxhs[0] = &Split{
				Position:  b.Min.X,
				BecauseOf: []shared.Shape{shape},
				Alignment: Horizontal,
			}
			if minxi >= 0 && minxi < len(m.HSplits) {
				minxhs[1] = m.HSplits[minxi]
			}
			m.HSplits = append(m.HSplits[:minxi], append([]*Split{minxhs[0]}, m.HSplits[minxi:]...)...)
			maxxi++
		} else {
			m.HSplits[minxi].BecauseOf = append(m.HSplits[minxi].BecauseOf, shape)
		}
		if m.VSplits[minyi].Position != b.Min.Y {
			vs := &Split{
				Position:  b.Min.Y,
				BecauseOf: []shared.Shape{shape},
				Alignment: Vertical,
			}
			var pvs *Split
			if minyi >= 0 && minyi < len(m.VSplits) {
				pvs = m.VSplits[minyi]
			}
			m.VSplits = append(m.VSplits[:minyi], append([]*Split{vs}, m.VSplits[minyi:]...)...)
			var lhs *Split
			for _, hs := range m.HSplits {
				_, exists := m.Stacks[SC(hs, vs)]
				if exists {
					lhs = hs
					continue
				}
				if pvs != nil {
					if hs != minxhs[0] {
						m.Stacks[SC(hs, vs)] = append([]shared.Shape{}, m.Stacks[SC(hs, pvs)]...)
					} else if lhs != nil {
						m.Stacks[SC(hs, vs)] = append([]shared.Shape{}, m.Stacks[SC(lhs, pvs)]...)
					}
				}
				lhs = hs
			}
			maxyi++
		} else {
			m.VSplits[minxi].BecauseOf = append(m.VSplits[minxi].BecauseOf, shape)
		}
		if minxhs[0] != nil {
			for _, vs := range m.VSplits {
				_, exists := m.Stacks[SC(minxhs[0], vs)]
				if exists {
					continue
				}
				if minxhs[1] != nil {
					m.Stacks[SC(minxhs[0], vs)] = append([]shared.Shape{}, m.Stacks[SC(minxhs[1], vs)]...)
				}
			}
		}
	}
	for x := minxi; x <= maxxi; x++ {
		for y := minyi; y <= maxyi; y++ {
			hs := m.HSplits[x]
			vs := m.VSplits[y]
			m.Stacks[SC(hs, vs)] = append(m.Stacks[SC(hs, vs)], shape)
		}
	}
	return m
}

func (m *Struct) GetXYPositions(p image.Point) (int, int) {
	minxi := sort.Search(len(m.HSplits), func(i int) bool {
		return m.HSplits[i].Position >= p.X
	})
	minyi := sort.Search(len(m.VSplits), func(i int) bool {
		return m.VSplits[i].Position >= p.Y
	})
	return minxi, minyi
}

func (m *Struct) GetStackAt(x int, y int) []shared.Shape {
	xi, yi := m.GetXYPositions(image.Point{x, y})
	if xi >= 0 && yi >= 0 && xi <= len(m.HSplits) && yi <= len(m.VSplits) {
		var hs *Split = nil
		if xi < len(m.HSplits) {
			hs = m.HSplits[xi]
		}
		if (hs == nil || hs.Position != x) && xi > 0 {
			hs = m.HSplits[xi-1]
		}
		var vs *Split = nil
		if yi < len(m.VSplits) {
			vs = m.VSplits[yi]
		}
		if (vs == nil || vs.Position != y) && yi > 0 {
			vs = m.VSplits[yi-1]
		}
		if hs != nil && vs != nil && vs.Position <= y && hs.Position <= y {
			if s, ok := m.Stacks[SC(hs, vs)]; ok && s != nil {
				return s
			} else {
				log.Default()
			}
		}
	}
	return []shared.Shape{}
}

func New() *Struct {
	return &Struct{
		VSplits: []*Split{},
		HSplits: []*Split{},
		Stacks:  map[SplitCoordination][]shared.Shape{},
	}
}

type ShapeArray []shared.Shape

func (a ShapeArray) Remove(shape shared.Shape) ([]shared.Shape, int) {
	shrink := 0
	for i := range a {
		for a[i] == shape && len(a)-shrink > i {
			shrink++
			a[i] = a[len(a)-shrink]
		}
	}
	return a[:len(a)-shrink], shrink
}

type SplitArray []*Split

func (sa SplitArray) Remove(shape shared.Shape) ([]*Split, []*Split) {
	shrink := 0
	var removedFrom []*Split
	in := false
	for sai := range sa {
		var removeCount int
		sa[sai].BecauseOf, removeCount = ShapeArray(sa[sai].BecauseOf).Remove(shape)
		min, max := shape.Bounds().Min.Y, shape.Bounds().Max.Y
		switch sa[sai].Alignment {
		case Horizontal:
			min, max = shape.Bounds().Min.X, shape.Bounds().Max.X
		}
		if !in {
			in = min <= sa[sai].Position
		} else {
			in = max >= sa[sai].Position
		}
		if removeCount > 0 || in {
			removedFrom = append(removedFrom, sa[sai])
		}
		sa[sai-shrink] = sa[sai]
		if len(sa[sai].BecauseOf) == 0 {
			shrink++
		}
	}
	return sa[:len(sa)-shrink], removedFrom
}

func (m *Struct) Remove(shape shared.Shape) *Struct {
	var vRemovedFrom []*Split
	var hRemovedFrom []*Split
	m.VSplits, vRemovedFrom = SplitArray(m.VSplits).Remove(shape)
	m.HSplits, hRemovedFrom = SplitArray(m.HSplits).Remove(shape)
	var i int
	ks := map[SplitCoordination]struct{}{}
	for _, v := range vRemovedFrom {
		for _, h := range m.HSplits {
			k := SplitCoordination{HSplit: h, VSplit: v}
			ks[k] = struct{}{}
		}
	}
	for _, v := range m.VSplits {
		for _, h := range hRemovedFrom {
			k := SplitCoordination{HSplit: h, VSplit: v}
			ks[k] = struct{}{}
		}
	}
	for _, v := range vRemovedFrom {
		for _, h := range hRemovedFrom {
			k := SplitCoordination{HSplit: h, VSplit: v}
			ks[k] = struct{}{}
		}
	}
	for k := range ks {
		a, ok := m.Stacks[k]
		if !ok {
			continue
		}
		if len(k.HSplit.BecauseOf) == 0 || len(k.VSplit.BecauseOf) == 0 {
			delete(m.Stacks, k)
			continue
		}
		a, i = ShapeArray(a).Remove(shape)
		if i > 0 {
			m.Stacks[k] = a
		}
	}
	return m
}
