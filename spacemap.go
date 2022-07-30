package spacemap

import (
	"image"
	"log"
	"sort"
)

type Alignment int

const (
	NotAligned Alignment = iota
	Vertical
	Horizontal
)

type Split struct {
	Position  int
	BecauseOf []Shape
	Alignment Alignment
}

type SplitCoordination struct {
	HSplit *Split
	VSplit *Split
}

func (c SplitCoordination) Check() {
	if c.VSplit.Alignment != Vertical {
		panic("Alignment check fail")
	}
	if c.HSplit.Alignment != Horizontal {
		panic("Alignment check fail")
	}
}

func SC(HSplit *Split, VSplit *Split) SplitCoordination {
	coordination := SplitCoordination{
		HSplit,
		VSplit,
	}
	coordination.Check()
	return coordination
}

type SpaceMap struct {
	VSplits []*Split
	HSplits []*Split
	Stacks  map[SplitCoordination][]Shape
}

func (m *SpaceMap) Add(shape Shape) *SpaceMap {
	b := shape.Bounds()
	minxi, minyi := m.GetXYPositions(b.Min)
	maxxi, maxyi := m.GetXYPositions(b.Max)
	for x := minxi; x < maxxi; x++ {
		for y := minyi; y < maxyi; y++ {
			hs := m.HSplits[x]
			vs := m.VSplits[y]
			m.Stacks[SC(hs, vs)] = append(m.Stacks[SC(hs, vs)], shape)
		}
	}
	{
		var maxxhs [2]*Split
		if maxxi == len(m.HSplits) || m.HSplits[maxxi].Position != b.Max.X {
			maxxhs[0] = &Split{
				Position:  b.Max.X,
				BecauseOf: []Shape{shape},
				Alignment: Horizontal,
			}
			m.HSplits = append(m.HSplits[:maxxi], append([]*Split{maxxhs[0]}, m.HSplits[maxxi:]...)...)
			if maxxi-1 >= 0 {
				maxxhs[1] = m.HSplits[maxxi-1]
			}
		} else {
			m.HSplits[maxxi].BecauseOf = append(m.HSplits[maxxi].BecauseOf, shape)
		}
		if maxyi == len(m.VSplits) || m.VSplits[maxyi].Position != b.Max.Y {
			vs := &Split{
				Position:  b.Max.Y,
				BecauseOf: []Shape{shape},
				Alignment: Vertical,
			}
			m.VSplits = append(m.VSplits[:maxyi], append([]*Split{vs}, m.VSplits[maxyi:]...)...)
			var pvs *Split = nil
			if maxyi-1 >= 0 {
				pvs = m.VSplits[maxyi-1]
			}
			for x := 0; x < len(m.HSplits); x++ {
				hs := m.HSplits[x]
				var ps = []Shape{}
				if pvs != nil {
					ps = m.Stacks[SC(hs, pvs)]
				}
				m.Stacks[SC(hs, vs)] = append(append([]Shape{}, ps...), shape)
			}
		} else {
			m.VSplits[maxxi].BecauseOf = append(m.VSplits[maxxi].BecauseOf, shape)
		}
		if maxxhs[0] != nil {
			for y := 0; y < len(m.VSplits); y++ {
				vs := m.VSplits[y]
				var ps = []Shape{}
				if maxxhs[1] != nil {
					ps = m.Stacks[SC(maxxhs[1], vs)]
				}
				m.Stacks[SC(maxxhs[0], vs)] = append(append([]Shape{}, ps...), shape)
			}
		}
	}
	{
		var minxhs [2]*Split
		if m.HSplits[minxi].Position != b.Min.X {
			minxhs[0] = &Split{
				Position:  b.Min.X,
				BecauseOf: []Shape{shape},
				Alignment: Horizontal,
			}
			m.HSplits = append(m.HSplits[:minxi], append([]*Split{minxhs[0]}, m.HSplits[minxi:]...)...)
			if minxi-1 >= 0 {
				minxhs[1] = m.HSplits[minxi-1]
			}
		} else {
			m.HSplits[minxi].BecauseOf = append(m.HSplits[minxi].BecauseOf, shape)
		}
		if m.VSplits[minyi].Position != b.Min.Y {
			vs := &Split{
				Position:  b.Min.Y,
				BecauseOf: []Shape{shape},
				Alignment: Vertical,
			}
			m.VSplits = append(m.VSplits[:minyi], append([]*Split{vs}, m.VSplits[minyi:]...)...)
			var pvs *Split = nil
			if minyi-1 >= 0 {
				pvs = m.VSplits[minyi-1]
			}
			for x := 0; x < len(m.HSplits); x++ {
				ohs := m.HSplits[x]
				var ps = []Shape{}
				if pvs != nil {
					ps = m.Stacks[SC(ohs, pvs)]
				}
				m.Stacks[SC(ohs, vs)] = append(append([]Shape{}, ps...), shape)
			}
		} else {
			m.VSplits[minxi].BecauseOf = append(m.VSplits[minxi].BecauseOf, shape)
		}
		if minxhs[0] != nil {
			for y := 0; y < len(m.VSplits); y++ {
				vs := m.VSplits[y]
				var ps = []Shape{}
				if minxhs[1] != nil {
					ps = m.Stacks[SC(minxhs[1], vs)]
				}
				m.Stacks[SC(minxhs[0], vs)] = append(append([]Shape{}, ps...), shape)
			}
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
			if s, ok := m.Stacks[SC(hs, vs)]; ok && s != nil {
				return s
			} else {
				log.Default()
			}
		}
	}
	return []Shape{}
}

func NewSpaceMap() *SpaceMap {
	return &SpaceMap{
		VSplits: []*Split{},
		HSplits: []*Split{},
		Stacks:  map[SplitCoordination][]Shape{},
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
