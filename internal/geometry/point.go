package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
)

type Point struct {
	Coordinates
	Name string
}

// MarshallGCode implements the Marshaler interface.
func (p Point) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	return []byte(fmt.Sprintf(";------ Point %s\nG0 X%.01F Y%.01f\nG1 Z%.01f F%.01f; Tool down\nG0 Z%.01f; Tool up\n",
		p.Name,
		p.X, p.Y,
		-options.Deep,
		options.Feed,
		options.SecurityZ,
	)), nil
}

func PointsFromDXFPoints(configs ...dxfConfigurator) []Point {
	dxfFile := dxf{}

	for _, config := range configs {
		config(&dxfFile)
	}

	inputPoints := make([]Linker, len(dxfFile.points))

	for idx, dxfPoint := range dxfFile.points {
		inputPoints[idx] = Point{
			Name: fmt.Sprintf("#%d / Layer %s", idx, dxfPoint.Layer().Name()),
			Coordinates: Coordinates{
				X: dxfPoint.Coord[0],
				Y: dxfPoint.Coord[1],
			},
		}
	}

	points, _ := SortEntities(inputPoints, &Coordinates{X: 0, Y: 0}, func(from, to Linker) bool {
		return true
	})

	output := []Point{}
	for _, point := range points {
		if value, ok := point.(Point); ok {
			output = append(output, value)
		}
	}

	return output
}
