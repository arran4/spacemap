package spacemap

import "spacemap/shared"

type Interface[T any] interface {
	Add(shape shared.Shape, zIndex int) T
	Remove(shape shared.Shape) T
	GetStackAt(x int, y int) []shared.Shape
	GetAt(x int, y int) shared.Shape
}
