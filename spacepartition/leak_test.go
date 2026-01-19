package spacepartition

import (
	"testing"
	"github.com/arran4/spacemap/shared"
)

func TestSplitArrayRemoveLeak(t *testing.T) {
	// Create a SplitArray with a few elements
	s1 := &Split{Position: 10, Alignment: Vertical, BecauseOf: []shared.Shape{shared.NewRectangle(0,0,1,1)}}

	shapeToRemove := shared.NewRectangle(0, 0, 100, 100)
	s2 := &Split{Position: 20, Alignment: Vertical, BecauseOf: []shared.Shape{shapeToRemove}}

	s3 := &Split{Position: 30, Alignment: Vertical, BecauseOf: []shared.Shape{shared.NewRectangle(0,0,1,1)}}

	sa := SplitArray{s1, s2, s3}

	// Perform remove
	newSa, _ := sa.Remove(shapeToRemove)

	// Check new length
	if len(newSa) != 2 {
		t.Fatalf("Expected length 2, got %d", len(newSa))
	}

	// Check for leak
	// The backing array has capacity 3.
	// After removing s2 (middle), the array becomes [s1, s3, s3].
	// The element at index 2 (last one) should be nil if we want to fix the leak.

	leakedSlice := newSa[:cap(newSa)]
	if leakedSlice[2] != nil {
		t.Logf("Memory leak detected: element at index 2 is %v", leakedSlice[2])
		// We fail if we want to assert the leak is present (for baseline) or absent (for verification).
		// Since this is a "reproduction test", failing means we successfully reproduced the issue.
		// However, for the final state, we want this to PASS.
		// So I will write this test to FAIL if the leak is present.
		t.Fail()
	}
}

func BenchmarkSplitArrayRemove(b *testing.B) {
	s1 := &Split{Position: 10, Alignment: Vertical, BecauseOf: []shared.Shape{shared.NewRectangle(0,0,1,1)}}
	shapeToRemove := shared.NewRectangle(0, 0, 100, 100)
	s2 := &Split{Position: 20, Alignment: Vertical, BecauseOf: []shared.Shape{shapeToRemove}}
	s3 := &Split{Position: 30, Alignment: Vertical, BecauseOf: []shared.Shape{shared.NewRectangle(0,0,1,1)}}

	baseSa := SplitArray{s1, s2, s3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Copy baseSa so we don't mutate the original in the loop
		sa := make(SplitArray, len(baseSa))
		copy(sa, baseSa)

		sa.Remove(shapeToRemove)
	}
}
