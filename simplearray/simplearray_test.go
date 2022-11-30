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
				Shapes: []*Point{
					&Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*Point{},
			},
			shape: rect1,
		},
		{
			name: "Remove missing does nothing",
			Struct: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect1},
				},
			},
			shape: rect2,
		},
		{
			name: "Remove doesn't impact the following",
			Struct: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect1},
					&Point{Shape: rect2},
				},
			},
			want: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect2},
				},
			},
			shape: rect1,
		},
		{
			name: "Remove doesn't impact the preceding",
			Struct: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect2},
					&Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect2},
				},
			},
			shape: rect1,
		},
		{
			name: "Multiple remove of a heterogeneous set is fine",
			Struct: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect2},
					&Point{Shape: rect1},
					&Point{Shape: rect3},
					&Point{Shape: rect1},
					&Point{Shape: rect4},
					&Point{Shape: rect1},
					&Point{Shape: rect3},
					&Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect2},
					&Point{Shape: rect3},
					&Point{Shape: rect3},
					&Point{Shape: rect4},
				},
			},
			shape: rect1,
		},
		{
			name: "Multiple remove of a homogenous set is fine",
			Struct: &Struct{
				Shapes: []*Point{
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
					&Point{Shape: rect1},
				},
			},
			want: &Struct{
				Shapes: []*Point{},
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
