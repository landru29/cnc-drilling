package shape

import (
	"fmt"
	"math"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/machine"
	"github.com/yofu/dxf/entity"
)

// Point is a 2D point.
type Point struct {
	geometry.Coordinates
	Name string
}

// DistanceTo computes the distance to another point.
func (p Point) DistanceTo(other Point) float64 {
	dx := other.X - p.X
	dy := other.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// NewPointFromPoint is a builder.
func NewPointFromPoint(name string, data *entity.Point) *Point {
	return &Point{
		Name:        name,
		Coordinates: geometry.NewCoordinatesFromPoint(data),
	}
}

// NewPointFromVertex is a builder.
func NewPointFromVertex(name string, data *entity.Vertex) Point {
	return Point{
		Name:        name,
		Coordinates: geometry.NewCoordinatesFromVertex(data),
	}
}

// MarshallGCode implements the Marshaler interface.
func (p Point) MarshallGCode(state *machine.Path, configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	return []byte(fmt.Sprintf(";------ Point %s\nG0 X%.03f Y%.03f\nG1 Z%.03f F%.03f; Tool down\nG0 Z%.03f; Tool up\n",
		p.Name,
		p.X-options.OffsetX(),
		p.Y-options.OffsetY(),
		-options.Deep,
		options.Feed,
		options.SecurityZ,
	)), nil
}

// Weight implements the Linker interface.
func (p Point) Weight(other Linker) [2]float64 {
	output := [2]float64{0, 0}

	if start := other.Start(); start != nil {
		output[0] = p.Coordinates.Weight(*start)
	}

	if end := other.End(); end != nil {
		output[1] = p.Coordinates.Weight(*end)
	}

	return output

}
