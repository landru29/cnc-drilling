package geometry

import (
	"fmt"
	"math"
)

// Box is the cutting box.
type Box struct {
	Min Coordinates
	Max Coordinates
}

func (b Box) Merge(others ...Box) Box {
	output := b

	for _, box := range others {
		output = Box{
			Min: Coordinates{
				X: math.Min(output.Min.X, box.Min.X),
				Y: math.Min(output.Min.Y, box.Min.Y),
			},
			Max: Coordinates{
				X: math.Max(output.Max.X, box.Max.X),
				Y: math.Max(output.Max.Y, box.Max.Y),
			},
		}
	}

	return output
}

// String implements the Stringer interface.
func (b Box) String() string {
	return fmt.Sprintf(
		"[(%.03f, %.03f), (%.03f, %.03f)]",
		b.Min.X,
		b.Min.Y,
		b.Max.X,
		b.Max.Y,
	)
}
