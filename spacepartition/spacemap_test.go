package spaceparition

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"image"
	"spacemap/shared"
	"strings"
	"testing"
)

func NSM(shapes ...shared.Shape) func() *Struct {
	return func() *Struct {
		return New().AddAll(shapes...)
	}
}

func TestSpaceMap(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100)
	rect2 := shared.NewRectangle(40, 40, 60, 60)
	rect3 := shared.NewRectangle(10, 10, 60, 60)
	rect4 := shared.NewRectangle(60, 60, 100, 100)
	for _, test := range []struct {
		Name      string
		Stack     []shared.Shape
		Position  image.Point
		SpaceMap  func() *Struct
		NumberMap [][]int
	}{
		{
			Name:      "Hit",
			Stack:     []shared.Shape{rect1},
			Position:  image.Point{20, 20},
			SpaceMap:  NSM(rect1),
			NumberMap: [][]int{{1, 1}, {1, 1}},
		},
		{
			Name:      "Hit Low Border",
			Stack:     []shared.Shape{rect1},
			Position:  rect1.Min,
			SpaceMap:  NSM(rect1),
			NumberMap: [][]int{{1, 1}, {1, 1}},
		},
		{
			Name:      "Hit High Border",
			Stack:     []shared.Shape{rect1},
			Position:  rect1.Max,
			SpaceMap:  NSM(rect1),
			NumberMap: [][]int{{1, 1}, {1, 1}},
		},
		{
			Name:      "Miss -- Near Hit High Border",
			Stack:     []shared.Shape{rect1},
			Position:  rect1.Max.Add(image.Pt(1, 1)),
			SpaceMap:  NSM(rect1),
			NumberMap: [][]int{{1, 1}, {1, 1}},
		},
		{
			Name:      "Miss",
			Stack:     []shared.Shape{},
			Position:  image.Point{-20, -20},
			SpaceMap:  NSM(rect1),
			NumberMap: [][]int{{1, 1}, {1, 1}},
		},
		{
			Name:      "Hit first with 2 overlapping",
			Stack:     []shared.Shape{rect1},
			Position:  image.Point{20, 20},
			SpaceMap:  NSM(rect1, rect2),
			NumberMap: [][]int{{1, 1, 1, 1}, {1, 2, 2, 1}, {1, 2, 2, 1}, {1, 1, 1, 1}},
		},
		{
			Name:      "Hit both with 2 overlapping",
			Stack:     []shared.Shape{rect1, rect2},
			Position:  image.Point{50, 50},
			SpaceMap:  NSM(rect1, rect2),
			NumberMap: [][]int{{1, 1, 1, 1}, {1, 2, 2, 1}, {1, 2, 2, 1}, {1, 1, 1, 1}},
		},
		{
			Name:      "Hit both with 2 overlapping same start",
			Stack:     []shared.Shape{rect1, rect3},
			Position:  image.Point{20, 20},
			SpaceMap:  NSM(rect1, rect3),
			NumberMap: [][]int{{2, 2, 1}, {2, 2, 1}, {1, 1, 1}},
		},
		{
			Name:      "Hit both with 2 overlapping same end",
			Stack:     []shared.Shape{rect1, rect4},
			Position:  image.Point{90, 90},
			SpaceMap:  NSM(rect1, rect4),
			NumberMap: [][]int{{1, 1, 1}, {1, 2, 2}, {1, 2, 2}},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			sm := test.SpaceMap()
			if len(sm.VSplits)*len(sm.HSplits) != len(sm.Stacks) {
				t.Errorf("Incorrect number of cells expected %d ( %d * %d ) but got %d",
					len(sm.VSplits)*len(sm.HSplits), len(sm.VSplits),
					len(sm.HSplits), len(sm.Stacks))
			}
			for k := range sm.Stacks {
				t.Run(fmt.Sprintf("Map check: %v %v", k.HSplit.Position, k.VSplit.Position), func(t *testing.T) {
					if !k.Sane() {
						t.Errorf("Not sane")
					}
				})
			}
			if test.NumberMap != nil {
				numberMap := NumberMapper(sm)
				if s := cmp.Diff(numberMap, test.NumberMap); len(s) > 0 {
					t.Errorf("Failed stacks differ: %s", s)
				}
			}
			stack := sm.GetStackAt(test.Position.X, test.Position.Y)
			if s := cmp.Diff(stack, test.Stack); len(s) > 0 {
				t.Errorf("Failed stacks differ: %s", s)
			}
			if t.Failed() {
				LogStructure(t, sm)
				LogStructureContents(t, sm)
			}
		})
	}
}

func NumberMapper(spaceMap *Struct) (result [][]int) {
	result = make([][]int, len(spaceMap.HSplits))
	for hi, h := range spaceMap.HSplits {
		result[hi] = make([]int, len(spaceMap.VSplits))
		for vi, v := range spaceMap.VSplits {
			if sq, ok := spaceMap.Stacks[SplitCoordination{h, v}]; ok {
				result[hi][vi] = len(sq)
			}
		}
	}
	return
}

func LogStructure(t *testing.T, spaceMap *Struct) {
	b := strings.Builder{}
	b.WriteString("\n____\t")
	for _, v := range spaceMap.VSplits {
		b.WriteString(fmt.Sprintf("p: %d,\t", v.Position))
	}
	b.WriteString("\n")
	for _, h := range spaceMap.HSplits {
		b.WriteString(fmt.Sprintf("p: %d,\t", h.Position))
		for _, v := range spaceMap.VSplits {
			if sq, ok := spaceMap.Stacks[SplitCoordination{h, v}]; ok {
				b.WriteString(fmt.Sprintf("c: %d,\t", len(sq)))
			} else {
				b.WriteString("null,\t")
			}
		}
		b.WriteString("\n")
	}
	t.Log(b.String())
}

func LogStructureContents(t *testing.T, spaceMap *Struct) {
	b := strings.Builder{}
	const format = "%40s"
	b.WriteString(fmt.Sprintf("\n"+format, "____"))
	for _, v := range spaceMap.VSplits {
		b.WriteString(fmt.Sprintf(format, fmt.Sprintf("p: %d,", v.Position)))
	}
	b.WriteString("\n")
	for _, h := range spaceMap.HSplits {
		b.WriteString(fmt.Sprintf(format, fmt.Sprintf("p: %d,", h.Position)))
		for _, v := range spaceMap.VSplits {
			if sq, ok := spaceMap.Stacks[SplitCoordination{h, v}]; ok {
				bb := strings.Builder{}
				for _, sqe := range sq {
					bb.WriteString(fmt.Sprintf("%s,", sqe.Bounds()))
				}
				b.WriteString(fmt.Sprintf(format, bb.String()))
			} else {
				b.WriteString(fmt.Sprintf(format, "null,"))
			}
		}
		b.WriteString("\n")
	}
	t.Log(b.String())
}

func TestStruct_Remove(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(40, 40, 60, 60, shared.Name("rect2"))
	rect3 := shared.NewRectangle(10, 10, 60, 60)
	tests := []struct {
		name        string
		Constructor func() *Struct
		shape       shared.Shape
		want        *Struct
		NumberMap   [][]int
		StackCount  int
	}{
		{
			name:        "Remove just one and all",
			Constructor: NSM(rect1),
			shape:       rect1,
			NumberMap:   [][]int{},
			StackCount:  0,
		},
		{
			name:        "Remove just one, one remains, no edge overlap",
			Constructor: NSM(rect1, rect2),
			shape:       rect1,
			NumberMap:   [][]int{{1, 1}, {1, 1}},
			StackCount:  4,
		},
		{
			name:        "Remove failed one, two remain, no edge overlap",
			Constructor: NSM(rect1, rect2),
			shape:       rect3,
			NumberMap:   [][]int{{1, 1, 1, 1}, {1, 2, 2, 1}, {1, 2, 2, 1}, {1, 1, 1, 1}},
			StackCount:  16,
		},
		{
			name:        "Remove just one, one remains, one edge overlap",
			Constructor: NSM(rect1, rect3),
			shape:       rect1,
			NumberMap:   [][]int{{1, 1}, {1, 1}},
			StackCount:  4,
		},
		{
			name:        "Remove failed one, two remain, one edge overlap",
			Constructor: NSM(rect1, rect3),
			shape:       rect2,
			NumberMap:   [][]int{{2, 2, 1}, {2, 2, 1}, {1, 1, 1}},
			StackCount:  9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.Constructor().Remove(tt.shape)
			if tt.NumberMap != nil {
				numberMap := NumberMapper(got)
				if s := cmp.Diff(numberMap, tt.NumberMap); len(s) > 0 {
					t.Errorf("Failed stacks differ:\n%s", s)
				}
			}
			if tt.StackCount != len(got.Stacks) {
				t.Errorf("Stack count mismatch: expected %d got %d", tt.StackCount, len(got.Stacks))
			}
		})
	}
}
