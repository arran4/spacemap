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

type SpaceMap struct {
	VSplits []*Split
	HSplits []*Split
	Stacks  map[SplitCoordination][]shared.Shape
}

func (m *SpaceMap) AddAll(shapes ...shared.Shape) *SpaceMap {
	for _, shape := range shapes {
		m.Add(shape)
	}
	return m
}

func (m *SpaceMap) Add(shape shared.Shape) *SpaceMap {
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

func (m *SpaceMap) GetXYPositions(p image.Point) (int, int) {
	minxi := sort.Search(len(m.HSplits), func(i int) bool {
		return m.HSplits[i].Position >= p.X
	})
	minyi := sort.Search(len(m.VSplits), func(i int) bool {
		return m.VSplits[i].Position >= p.Y
	})
	return minxi, minyi
}

func (m *SpaceMap) GetStackAt(x int, y int) []shared.Shape {
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

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{
		VSplits: []*Split{},
		HSplits: []*Split{},
		Stacks:  map[SplitCoordination][]shared.Shape{},
	}
}
