package geometry

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/yofu/dxf/entity"
)

// CoordinatesXY is 2D coordinates.
type CoordinatesXY struct {
	X float64
	Y float64
}

func (c CoordinatesXY) DistanceTo(other CoordinatesXY) float64 {
	dx := c.X - other.X
	dy := c.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// NewCoordinatesFromPoint is a coordinate builder.
func NewCoordinatesFromPoint(data *entity.Point) CoordinatesXY {
	if data == nil {
		return CoordinatesXY{}
	}

	return CoordinatesXY{
		X: data.Coord[0],
		Y: data.Coord[1],
	}
}

// NewCoordinatesFromVertex is a coordinate builder.
func NewCoordinatesFromVertex(data *entity.Vertex) CoordinatesXY {
	if data == nil {
		return CoordinatesXY{}
	}

	return CoordinatesXY{
		X: data.Coord[0],
		Y: data.Coord[1],
	}
}

// String implements the pflag.Value interface.
func (c CoordinatesXY) String() string {
	return fmt.Sprintf("(%.01f, %.01f)", c.X, c.Y)
}

// Set implements the pflag.Value interface.
func (c *CoordinatesXY) Set(data string) error {
	splitter := strings.Split(data, ",")
	if len(splitter) < 2 {
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
func (c CoordinatesXY) Type() string {
	return "CoordinateXY"
}

// Start implements the Linker interface.
func (c CoordinatesXY) Start() *CoordinatesXY {
	return &c
}

// End implements the Linker interface.
func (c CoordinatesXY) End() *CoordinatesXY {
	return &c
}

// Revert implements the Linker interface.
func (c CoordinatesXY) Revert() {}

func (c CoordinatesXY) Weight(other CoordinatesXY) float64 {
	return (c.X-other.X)*(c.X-other.X) + (c.Y-other.Y)*(c.Y-other.Y)
}

func (c CoordinatesXY) Equal(other CoordinatesXY) bool {
	return c.Weight(other) < 0.00001
}

// Box implements the Linker interface.
func (c CoordinatesXY) Box() Box {
	return Box{
		Min: c,
		Max: c,
	}
}
