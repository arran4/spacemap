package spacemap

import (
	"image"
	"spacemap/shared"
	"spacemap/simplearray"
	"spacemap/space2trees"
	"spacemap/spacepartition"
	"testing"
)

var (
	benchShapes  = GenerateBenchShapes(100)
	spaceLookups = GenerateSpaceLookups(100)
)

func GenerateSpaceLookups(limit int) (result []image.Point) {
	for top := 10; top < limit; top += 8 {
		for left := 10; left < limit; left += 8 {
			result = append(result, image.Pt(top, left))
		}
	}
	return result
}

func GenerateBenchShapes(limit int) (result []shared.Shape) {
	for top := 10; top < limit; top += 10 {
		for left := 10; left < limit; left += 10 {
			for right := left + 10; right < left+limit; right += 10 {
				for bottom := top + 10; bottom < limit+top; bottom += 10 {
					result = append(result, shared.NewRectangle(left, top, right, bottom))
				}
			}
		}
	}
	return result
}

var (
	_ Interface[*space2trees.Struct]   = (*space2trees.Struct)(nil)
	_ Interface[*spaceparition.Struct] = (*spaceparition.Struct)(nil)
	_ Interface[*simplearray.Struct]   = (*simplearray.Struct)(nil)
)

func BenchmarkSpacePartitionAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := spaceparition.New()
		sm.AddAll(benchShapes...)
	}
}

func BenchmarkSpaceBTreeAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := space2trees.New()
		sm.AddAll(benchShapes...)
	}
}

func BenchmarkSimpleArrayAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := simplearray.New()
		sm.AddAll(benchShapes...)
	}
}

func BenchmarkSpacePartitionAddSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := spaceparition.New()
		sm.AddAll(benchShapes...)
		for _, l := range spaceLookups {
			sm.GetStackAt(l.X, l.Y)
		}
	}
}

func BenchmarkSpaceBTreeAddSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := space2trees.New()
		sm.AddAll(benchShapes...)
		for _, l := range spaceLookups {
			sm.GetStackAt(l.X, l.Y)
		}
	}
}

func BenchmarkSimpleArrayAddSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := simplearray.New()
		sm.AddAll(benchShapes...)
		for _, l := range spaceLookups {
			sm.GetStackAt(l.X, l.Y)
		}
	}
}

func BenchmarkSpacePartitionAddDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := spaceparition.New()
		sm.AddAll(benchShapes...)
		for _, l := range benchShapes {
			sm.Remove(l)
		}
	}
}

func BenchmarkSpaceBTreeAddDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := space2trees.New()
		sm.AddAll(benchShapes...)
		for _, l := range benchShapes {
			sm.Remove(l)
		}
	}
}

func BenchmarkSimpleArrayAddDelete(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sm := simplearray.New()
		sm.AddAll(benchShapes...)
		for _, l := range benchShapes {
			sm.Remove(l)
		}
	}
}
