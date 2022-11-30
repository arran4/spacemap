package simplearray

import (
	"github.com/google/go-cmp/cmp"
	"spacemap/shared"
	"testing"
)

func TestStruct_Remove(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(40, 40, 60, 60, shared.Name("rect2"))
	rect3 := shared.NewRectangle(10, 10, 60, 60, shared.Name("rect3"))
	rect4 := shared.NewRectangle(60, 60, 100, 100, shared.Name("rect4"))
	tests := []struct {
		name string
		*Struct
		shape shared.Shape
		want  *Struct
	}{
		{
			name: "Remove removes 1",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{},
			},
			shape: rect1,
		},
		{
			name: "Remove missing does nothing",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect1},
				},
			},
			shape: rect2,
		},
		{
			name: "Remove doesn't impact the following",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect2},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect2},
				},
			},
			shape: rect1,
		},
		{
			name: "Remove doesn't impact the preceding",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect2},
					&shared.Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect2},
				},
			},
			shape: rect1,
		},
		{
			name: "Multiple remove of a heterogeneous set is fine",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect2},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect3},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect4},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect3},
					&shared.Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect2},
					&shared.Point{Shape: rect3},
					&shared.Point{Shape: rect3},
					&shared.Point{Shape: rect4},
				},
			},
			shape: rect1,
		},
		{
			name: "Multiple remove of a homogenous set is fine",
			Struct: &Struct{
				Shapes: []*shared.Point{
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
					&shared.Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*shared.Point{},
			},
			shape: rect1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.Remove(tt.shape)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Remove() = \n%s", diff)
			}
		})
	}
}

func TestZIndex(t *testing.T) {
	rect1 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect1"))
	rect2 := shared.NewRectangle(10, 10, 100, 100, shared.Name("rect2"))
	tests := []struct {
		name        string
		Constructor func() *Struct
		want        shared.Shape
	}{
		{
			name: "r1&2 want r2",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect1, 1)
				s.Add(rect2, 0)
				return s
			},
			want: rect2,
		},
		{
			name: "r1&2 want r1",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect1, 1)
				s.Add(rect2, 2)
				return s
			},
			want: rect1,
		},
		{
			name: "r2&1 want r2",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect2, 0)
				s.Add(rect1, 1)
				return s
			},
			want: rect2,
		},
		{
			name: "r2&1 want r1",
			Constructor: func() *Struct {
				s := New()
				s.Add(rect2, 2)
				s.Add(rect1, 1)
				return s
			},
			want: rect1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			sm := test.Constructor()
			r := sm.GetAt(50, 50)
			if r != test.want {
				tt.Errorf("Failed got %s expected %s", r, test.want)
			}
		})
	}
}
