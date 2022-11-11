package btree

import (
	"github.com/google/go-cmp/cmp"
	"spacemap/shared"
	"testing"
)

func NSM(shapes ...shared.Shape) func() *SpaceMap {
	return func() *SpaceMap {
		return NewSpaceMap().AddAll(shapes...)
	}
}

/*func TestSpaceMap(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100)
	rect2 := shared.NewRectangle(40, 40, 60, 60)
	rect3 := shared.NewRectangle(10, 10, 60, 60)
	rect4 := shared.NewRectangle(60, 60, 100, 100)
	for _, test := range []struct {
		Name      string
		Stack     []shared.Shape
		Position  image.Point
		SpaceMap  func() *SpaceMap
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
*/

func TestHereInsert(t *testing.T) {
	tests := []struct {
		name  string
		n     []*Here
		wantN []*Here
		zi    int
		wantP int
	}{
		{
			name:  "Empty",
			n:     []*Here{},
			zi:    0,
			wantP: 0,
			wantN: []*Here{},
		},
		{
			name: "1 Before 2",
			n: []*Here{
				{
					ZIndex: 2,
				},
			},
			zi:    1,
			wantP: 0,
			wantN: []*Here{
				{
					ZIndex: 2,
				},
			},
		},
		{
			name: "2 Before 2",
			n: []*Here{
				{
					ZIndex: 2,
				},
			},
			zi:    2,
			wantP: 0,
			wantN: []*Here{
				{
					ZIndex: 3,
				},
			},
		},
		{
			name: "1 Between 0,1",
			n: []*Here{
				{
					ZIndex: 0,
				},
				{
					ZIndex: 1,
				},
			},
			zi:    1,
			wantP: 1,
			wantN: []*Here{
				{
					ZIndex: 0,
				},
				{
					ZIndex: 2,
				},
			},
		},
		{
			name: "0 Before 0,1",
			n: []*Here{
				{
					ZIndex: 0,
				},
				{
					ZIndex: 1,
				},
			},
			zi:    0,
			wantP: 0,
			wantN: []*Here{
				{
					ZIndex: 1,
				},
				{
					ZIndex: 2,
				},
			},
		},
		{
			name: "2 After 1",
			n: []*Here{
				{
					ZIndex: 0,
				},
				{
					ZIndex: 1,
				},
			},
			zi:    2,
			wantP: 2,
			wantN: []*Here{
				{
					ZIndex: 0,
				},
				{
					ZIndex: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotP := HereInsert(tt.n, tt.zi); gotP != tt.wantP {
				t.Errorf("HereInsert() = %v, want %v", gotP, tt.wantP)
			}
			if diff := cmp.Diff(tt.n, tt.wantN); diff != "" {
				t.Errorf("N Diff: \n%s", diff)
			}
		})
	}
}
