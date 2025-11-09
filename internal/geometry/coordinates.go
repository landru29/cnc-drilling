package geometry

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/yofu/dxf/entity"
)

// Coordinates is 2D coordinates.
type Coordinates struct {
	X float64
	Y float64
}

func (c Coordinates) DistanceTo(other Coordinates) float64 {
	dx := c.X - other.X
	dy := c.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// NewCoordinatesFromPoint is a coordinate builder.
func NewCoordinatesFromPoint(data *entity.Point) Coordinates {
	if data == nil {
		return Coordinates{}
	}

	return Coordinates{
		X: data.Coord[0],
		Y: data.Coord[1],
	}
}

// NewCoordinatesFromVertex is a coordinate builder.
func NewCoordinatesFromVertex(data *entity.Vertex) Coordinates {
	if data == nil {
		return Coordinates{}
	}

	return Coordinates{
		X: data.Coord[0],
		Y: data.Coord[1],
	}
}

// String implements the pflag.Value interface.
func (c Coordinates) String() string {
	return fmt.Sprintf("(%.01f, %.01f)", c.X, c.Y)
}

// Set implements the pflag.Value interface.
func (c *Coordinates) Set(data string) error {
	splitter := strings.Split(data, ",")
	if len(splitter) != 2 {
		return errors.New("coordinates must be 0.0,0.0")
	}

	xValue, err := strconv.ParseFloat(strings.TrimSpace(splitter[0]), 64)
	if err != nil {
		return err
	}

	yValue, err := strconv.ParseFloat(strings.TrimSpace(splitter[1]), 64)
	if err != nil {
		return err
	}

	c.X = xValue
	c.Y = yValue

	return nil
}

// Type implements the pflag.Value interface.
func (c Coordinates) Type() string {
	return "Coordinate"
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
