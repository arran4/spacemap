# SpaceMap

SpaceMap is a small Go library that explores several algorithms for spatial indexing.  It is designed to answer the question **"what shapes exist at a given coordinate"** while taking each shapes Z-index into account.  The library originated as an experiment for building fast UI and game image maps but can be used in any scenario where spatial queries with layering are required.

The project focuses on simple rectangle shapes and provides a consistent interface so that different implementations can be swapped without changing your application code.

## When Would You Use This?

- **Interactive UIs** – efficiently detect which widget is under the mouse cursor when multiple elements overlap.
- **2D games** – quick hit-testing of sprites and detecting render order at a point.
- **Editors and design tools** – keeping track of layered objects on a canvas and retrieving them by position.

If you only have a small number of elements you can rely on the simple array implementation.  For larger scenes where performance matters you can switch to a more advanced algorithm without altering the rest of your code.

## Available Implementations

The library contains three sub-packages that implement the `spacemap.Interface`:

| Package | Highlights |
|---------|-----------|
| `simplearray` | Straight forward slice of shapes.  Easy to understand but not optimised. |
| `space2trees` | Uses two binary trees (one per axis) to index the space.  Balances itself using an AVL tree approach. |
| `spacepartition` | Splits the plane into partitions as shapes are added.  This is efficient for large numbers of static or rarely moved shapes. |

Each package exposes a `New()` constructor returning a structure that implements the interface shown below.

```go
package spacemap

type Interface interface {
    Add(shape shared.Shape, zIndex int)
    Remove(shape shared.Shape)
    GetStackAt(x int, y int) []shared.Shape
    GetAt(x int, y int) shared.Shape
}
```

## Shapes

The `shared` package defines primitive shapes.  Currently only rectangles are implemented but other shapes can be added by implementing the `shared.Shape` interface.

```go
// Rectangle returns a new rectangle using screen coordinates
r := shared.NewRectangle(left, top, right, bottom, shared.Name("player"))
```

`Shape.PointIn` tests whether a coordinate falls within a shape, while `Bounds` returns its bounding box.

## Quick Start

1. Install the library:

```bash
go get github.com/arran4/spacemap
```

2. Choose an implementation and create a space map:

```go
import (
    "github.com/arran4/spacemap/simplearray" // or space2trees, spacepartition
    "github.com/arran4/spacemap/shared"
)

sm := simplearray.New()

// Add shapes with optional z-index (0 is default)
sm.Add(shared.NewRectangle(10, 10, 100, 100), 0)
sm.Add(shared.NewRectangle(40, 40, 60, 60), 1)

shape := sm.GetAt(50, 50)     // highest Z at that point
stack := sm.GetStackAt(50,50) // all shapes at the point
```

3. Removing shapes is just as simple:

```go
sm.Remove(shape)
```

## Running Tests and Benchmarks

Tests cover all implementations and can be run with the standard Go tooling:

```bash
go test ./...
```

The `benchmark_test.go` file contains micro benchmarks comparing the performance of each implementation.  Run them with:

```bash
go test -bench .
```

## Status

This project is experimental but functional.  The algorithms are independent and may evolve separately. Contributions and bug reports are welcome.

