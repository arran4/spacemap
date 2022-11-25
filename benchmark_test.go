package spacemap

import (
	"image"
	"spacemap/shared"
	"spacemap/spacebtree"
	"spacemap/spacepartition"
	"testing"
)

var (
	benchShapes  = GenerateBenchShapes()
	spaceLookups = GenerateSpaceLookups()
)

func GenerateSpaceLookups() (result []image.Point) {
	for top := 10; top < 100; top += 8 {
		for left := 10; left < 100; left += 8 {
			result = append(result, image.Pt(top, left))
		}
	}
	return result
}

func GenerateBenchShapes() (result []shared.Shape) {
	for top := 10; top < 100; top += 10 {
		for left := 10; left < 100; left += 10 {
			for right := left + 10; right < left+100; right += 10 {
				for bottom := top + 10; bottom < 100+top; bottom += 10 {
					result = append(result, shared.NewRectangle(left, top, right, bottom))
				}
			}
		}
	}
	return result
}

type Interface[T Interface[T]] interface {
	AddAll(...shared.Shape) T
	Add(shared.Shape) T
	GetStackAt(x int, y int) []shared.Shape
}

var (
	_ Interface[*spacebtree.SpaceMap]    = (*spacebtree.SpaceMap)(nil)
	_ Interface[*spaceparition.SpaceMap] = (*spaceparition.SpaceMap)(nil)
)

func BenchmarkSpacePartitionAddSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := spaceparition.NewSpaceMap()
		sm.AddAll(benchShapes...)
		for _, l := range spaceLookups {
			sm.GetStackAt(l.X, l.Y)
		}
	}
}

func BenchmarkSpaceBTreeAddSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := spacebtree.NewSpaceMap()
		sm.AddAll(benchShapes...)
		for _, l := range spaceLookups {
			sm.GetStackAt(l.X, l.Y)
		}
	}
}
