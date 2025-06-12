package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
)

type Curve struct {
	Name       string
	StartPoint Coordinates
	EndPoint   Coordinates
	Center     Coordinates
	Radius     float64
	Clockwise  bool
}

// Start implements the Linker interface.
func (c Curve) Start() *Coordinates {
	return &c.StartPoint
}

// End implements the Linker interface.
func (c Curve) End() *Coordinates {
	return &c.EndPoint
}

// Revert implements the Linker interface.
func (c *Curve) Revert() {
	c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
	c.Clockwise = !c.Clockwise
}

// Weight implements the Linker interface.
func (c Curve) Weight(other Linker) [2]float64 {
	return c.EndPoint.Weight(other)
}

// MarshallGCode implements the Marshaler interface.
func (c Curve) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	output := ";------ Curve " + c.Name + "\n"

	if !options.IgnoreStart {
		start := c.Start()
		output = fmt.Sprintf(
			"G0 X%.01f Y%.01f\nG1 Z%.01f F%.01f ; Tool down\n",
			start.X,
			start.Y,
			-options.Deep,
			options.Feed,
		)
	}

	code := 2
	if c.Clockwise {
		code = 3
	}

	output += fmt.Sprintf(
		"G%d X%.01f Y%.01f I%.01f J%.01f F%.01f\n",
		code,
		c.EndPoint.X,
		c.EndPoint.Y,
		c.Center.X-c.StartPoint.X,
		c.Center.Y-c.StartPoint.Y,
		options.Feed,
	)

	return []byte(output), nil
}
