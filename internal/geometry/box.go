package geometry

import "math"

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
				X: math.Max(output.Min.X, box.Min.X),
				Y: math.Max(output.Min.Y, box.Min.Y),
			},
		}
	}

	return output
}
