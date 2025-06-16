package geometry

import (
	"fmt"
	"math"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/yofu/dxf/entity"
)

type Segment struct {
	Name       string
	StartPoint Coordinates
	EndPoint   Coordinates
}

func NewSgmentFromPoints(name string, from *entity.Point, to *entity.Point) *Segment {
	return &Segment{
		Name:       name,
		StartPoint: NewCoordinatesFromPoint(from),
		EndPoint:   NewCoordinatesFromPoint(to),
	}
}

func NewSgmentFromLine(name string, data *entity.Line) *Segment {
	return &Segment{
		Name: name,
		StartPoint: Coordinates{
			X: data.Start[0],
			Y: data.Start[1],
		},
		EndPoint: Coordinates{
			X: data.End[0],
			Y: data.End[1],
		},
	}
}

func (s Segment) Start() *Coordinates {
	return &s.StartPoint
}

func (s Segment) End() *Coordinates {
	return &s.EndPoint
}

func (s *Segment) Revert() {
	s.StartPoint, s.EndPoint = s.EndPoint, s.StartPoint
}

func (s Segment) Weight(other Linker) [2]float64 {
	return s.EndPoint.Weight(other)
}

// Box implements the Linker interface.
func (s Segment) Box() Box {
	return Box{
		Min: Coordinates{
			X: math.Min(s.StartPoint.X, s.EndPoint.X),
			Y: math.Min(s.StartPoint.Y, s.EndPoint.Y),
		},
		Max: Coordinates{
			X: math.Max(s.StartPoint.X, s.EndPoint.X),
			Y: math.Max(s.StartPoint.Y, s.EndPoint.Y),
		},
	}
}

// MarshallGCode implements the Marshaler interface.
func (s Segment) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	output := ";------ Segment " + s.Name + "\n"

	if !options.IgnoreStart {
		start := s.Start()
		output = fmt.Sprintf(
			"; * %s\nG0 X%.03f Y%.03f\nG1 Z%.03f F%.03f; Tool down\n",
			s.Name,
			start.X-options.OffsetX(),
			start.Y-options.OffsetY(),
			-options.Deep,
			options.Feed,
		)
	}

	output += fmt.Sprintf(
		"G1 X%.03f Y%.03f F%.03f\n",
		s.EndPoint.X-options.OffsetX(),
		s.EndPoint.Y-options.OffsetY(),
		options.Feed,
	)

	return []byte(output), nil
}
