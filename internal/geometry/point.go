package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/yofu/dxf/entity"
)

type Point struct {
	Coordinates
	Name string
}

func NewPointFromPoint(name string, data *entity.Point) *Point {
	return &Point{
		Name:        name,
		Coordinates: NewCoordinatesFromPoint(data),
	}
}

func NewPointFromVertex(name string, data *entity.Vertex) Point {
	return Point{
		Name:        name,
		Coordinates: NewCoordinatesFromVertex(data),
	}
}

// MarshallGCode implements the Marshaler interface.
func (p Point) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	return []byte(fmt.Sprintf(";------ Point %s\nG0 X%.01F Y%.01f\nG1 Z%.01f F%.01f; Tool down\nG0 Z%.01f; Tool up\n",
		p.Name,
		p.X-options.OffsetX(),
		p.Y-options.OffsetY(),
		-options.Deep,
		options.Feed,
		options.SecurityZ,
	)), nil
}
