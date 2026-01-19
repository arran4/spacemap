package shared

import (
	"testing"
)

func TestPointArray_Remove_MemoryLeak(t *testing.T) {
	// Create some shapes
	s1 := NewRectangle(0, 0, 10, 10, Name("s1"))
	s2 := NewRectangle(10, 10, 20, 20, Name("s2"))
	s3 := NewRectangle(20, 20, 30, 30, Name("s3"))

	p1 := &Point{ZIndex: 1, Shape: s1}
	p2 := &Point{ZIndex: 2, Shape: s2}
	p3 := &Point{ZIndex: 3, Shape: s3}

	// Initialize PointArray with enough capacity
	pa := make(PointArray, 0, 10)
	pa = append(pa, p1, p2, p3)

	// Remove p2
	res, shrink := pa.Remove(s2)

	if shrink != 1 {
		t.Errorf("Expected shrink to be 1, got %d", shrink)
	}

	if len(res) != 2 {
		t.Errorf("Expected result length to be 2, got %d", len(res))
	}

	// Check underlying array
	// We need to re-slice to capacity to see the "removed" elements
	underlying := res[:cap(res)]

	leakIndex := len(res)
	// The element at leakIndex corresponds to what was at the end of the slice before removal/swap
	// In the swap algorithm: pa[i] = pa[len-shrink].
	// The element at len-shrink remains there unless cleared.

	if underlying[leakIndex] != nil {
		t.Logf("Memory leak detected! Item at index %d in underlying array is not nil: %v", leakIndex, underlying[leakIndex])
        t.Fail()
	}
}

func BenchmarkPointArray_Remove(b *testing.B) {
	s := NewRectangle(0, 0, 10, 10)
	p := &Point{ZIndex: 1, Shape: s}

    // Create a large array to make the effect measurable if we were testing speed,
    // but here we are testing overhead of clearing.
    // However, allocating in the loop might dominate.

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		pa := make(PointArray, 100)
		for j := 0; j < 100; j++ {
			pa[j] = p
		}
		b.StartTimer()
		pa.Remove(s)
	}
}
