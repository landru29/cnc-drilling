package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
)

type Segment struct {
	Name       string
	StartPoint Coordinates
	EndPoint   Coordinates
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
			"; * %s\nG0 X%.01f Y%.01f\nG1 Z%.01f F%.01f; Tool down\n",
			s.Name,
			start.X,
			start.Y,
			-options.Deep,
			options.Feed,
		)
	}

	output += fmt.Sprintf(
		"G1 X%.01f Y%.01f F%.01f\n",
		s.EndPoint.X,
		s.EndPoint.Y,
		options.Feed,
	)

	return []byte(output), nil
}
