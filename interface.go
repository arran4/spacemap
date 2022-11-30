package spacemap

import "github.com/arran4/spacemap/shared"

type Interface interface {
	Add(shape shared.Shape, zIndex int)
	Remove(shape shared.Shape)
	GetStackAt(x int, y int) []shared.Shape
	GetAt(x int, y int) shared.Shape
}
