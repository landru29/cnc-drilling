package geometry

type Coordinates struct {
	X float64
	Y float64
}

// Start implements the Linker interface.
func (c Coordinates) Start() *Coordinates {
	return &c
}

// End implements the Linker interface.
func (c Coordinates) End() *Coordinates {
	return &c
}

// Revert implements the Linker interface.
func (c Coordinates) Revert() {}

// Weight implements the Linker interface.
func (c Coordinates) Weight(other Linker) [2]float64 {
	output := [2]float64{0, 0}

	if start := other.Start(); start != nil {
		output[0] = c.weight(*start)
	}

	if end := other.End(); end != nil {
		output[1] = c.weight(*end)
	}

	return output

}

func (c Coordinates) weight(other Coordinates) float64 {
	return (c.X-other.X)*(c.X-other.X) + (c.Y-other.Y)*(c.Y-other.Y)
}

func (c Coordinates) Equal(other Coordinates) bool {
	return c.weight(other) < 0.00001
}

// Box implements the Linker interface.
func (c Coordinates) Box() Box {
	return Box{
		Min: c,
		Max: c,
	}
}
