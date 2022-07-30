package spacemap

import (
	"image"
	"sort"
)

type Split struct {
	Position  int
	BecauseOf []Shape
}

type SpaceMap struct {
	VSplits []*Split
	HSplits []*Split
	Stacks  map[[2]*Split][]Shape
}

func (m *SpaceMap) Add(shape Shape) *SpaceMap {
	b := shape.Bounds()
	minxi, minyi := m.GetXYPositions(b.Min)
	maxxi, maxyi := m.GetXYPositions(b.Max)
	for x := minxi; x < maxxi; x++ {
		for y := minyi; y < maxyi; y++ {
			hs := m.HSplits[x]
			vs := m.VSplits[y]
			m.Stacks[[2]*Split{hs, vs}] = append(m.Stacks[[2]*Split{hs, vs}], shape)
		}
	}
	var maxvhs [2]*Split
	if maxxi == len(m.HSplits) || m.HSplits[maxxi].Position != b.Max.X {
		maxvhs[0] = &Split{
			Position:  b.Max.X,
			BecauseOf: []Shape{shape},
		}
		m.HSplits = append(append(m.HSplits[:maxxi], maxvhs[0]), m.HSplits[maxxi:]...)
		if maxxi-1 >= 0 {
			maxvhs[1] = m.HSplits[maxxi-1]
		}
	} else {
		m.HSplits[maxxi].BecauseOf = append(m.HSplits[maxxi].BecauseOf, shape)
	}
	if maxyi == len(m.VSplits) || m.VSplits[maxyi].Position != b.Max.Y {
		vs := &Split{
			Position:  b.Max.Y,
			BecauseOf: []Shape{shape},
		}
		m.VSplits = append(append(m.VSplits[:maxyi], vs), m.VSplits[maxyi:]...)
		var phs *Split = nil
		if maxyi-1 >= 0 {
			phs = m.VSplits[maxyi-1]
		}
		for x := 0; x < len(m.HSplits); x++ {
			ovs := m.HSplits[x]
			var ps = []Shape{}
			if phs != nil {
				ps = m.Stacks[[2]*Split{phs, ovs}]
			}
			m.Stacks[[2]*Split{vs, ovs}] = append(append([]Shape{}, ps...), shape)
		}
	} else {
		m.VSplits[maxxi].BecauseOf = append(m.VSplits[maxxi].BecauseOf, shape)
	}
	if maxvhs[0] != nil {
		for y := 0; y < len(m.VSplits); y++ {
			vs := m.VSplits[y]
			var ps = []Shape{}
			if maxvhs[1] != nil {
				ps = m.Stacks[[2]*Split{maxvhs[1], vs}]
			}
			m.Stacks[[2]*Split{maxvhs[0], vs}] = append(append([]Shape{}, ps...), shape)
		}
	}
	if m.HSplits[minxi].Position != b.Min.X {
		hs := &Split{
			Position:  b.Min.X,
			BecauseOf: []Shape{shape},
		}
		m.HSplits = append(append(m.HSplits[:minxi], hs), m.HSplits[minxi:]...)
		var phs *Split = nil
		if minxi-1 >= 0 {
			phs = m.HSplits[minxi-1]
		}
		for y := 0; y < len(m.VSplits); y++ {
			vs := m.VSplits[y]
			var ps = []Shape{}
			if phs != nil {
				ps = m.Stacks[[2]*Split{phs, vs}]
			}
			m.Stacks[[2]*Split{hs, vs}] = append(append([]Shape{}, ps...), shape)
		}
	} else {
		m.HSplits[minxi].BecauseOf = append(m.HSplits[minxi].BecauseOf, shape)
	}
	if m.VSplits[minyi].Position != b.Min.Y {
		vs := &Split{
			Position:  b.Min.Y,
			BecauseOf: []Shape{shape},
		}
		m.VSplits = append(append(m.VSplits[:minyi], vs), m.VSplits[minyi:]...)
		var phs *Split = nil
		if minyi-1 >= 0 {
			phs = m.VSplits[minyi-1]
		}
		for x := 0; x < len(m.HSplits); x++ {
			ovs := m.HSplits[x]
			var ps = []Shape{}
			if phs != nil {
				ps = m.Stacks[[2]*Split{phs, ovs}]
			}
			m.Stacks[[2]*Split{vs, ovs}] = append(append([]Shape{}, ps...), shape)
		}
	} else {
		m.VSplits[minxi].BecauseOf = append(m.VSplits[minxi].BecauseOf, shape)
	}
	var minvhs [2]*Split
	if minvhs[0] != nil {
		for y := 0; y < len(m.VSplits); y++ {
			vs := m.VSplits[y]
			var ps = []Shape{}
			if minvhs[1] != nil {
				ps = m.Stacks[[2]*Split{minvhs[1], vs}]
			}
			m.Stacks[[2]*Split{minvhs[0], vs}] = append(append([]Shape{}, ps...), shape)
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

func (m *SpaceMap) GetStackAt(x int, y int) []Shape {
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
		if hs != nil && vs != nil {
			if s, ok := m.Stacks[[2]*Split{hs, vs}]; ok && s != nil {
				return s
			}
		}
	}
	return []Shape{}
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{
		VSplits: []*Split{},
		HSplits: []*Split{},
		Stacks:  map[[2]*Split][]Shape{},
	}
}

type Shape interface {
	PointIn(x, y int) bool
	Bounds() image.Rectangle
}

type Rectangle image.Rectangle

func (r Rectangle) PointIn(x, y int) bool {
	return (image.Point{x, y}).In(image.Rectangle(r))
}

func (r Rectangle) Bounds() image.Rectangle {
	return image.Rectangle(r)
}

var _ Shape = (*Rectangle)(nil)

func NewRectangle(left, top, right, bottom int) *Rectangle {
	return &Rectangle{
		Min: image.Point{
			left,
			top,
		},
		Max: image.Point{
			right,
			bottom,
		},
	}
}
