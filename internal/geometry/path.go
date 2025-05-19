package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
)

type Path []Linker

// MarshallGCode implements the Marshaler interface.
func (p Path) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	start := p.Start()
	output := fmt.Sprintf(
		"G0 X%.01f Y%.01f\nG1 Z%.01f F%.01f \n",
		start.X,
		start.Y,
		-options.Deep,
		options.Feed,
	)

	for _, segment := range p {
		out, err := gcode.Marshal(segment, configs...)
		if err != nil {
			return nil, err
		}

		output += string(out)
	}

	return []byte(fmt.Sprintf("%sG0 Z%.01f\n", output, options.SecurityZ)), nil
}

// Start implements the Linker interface.
func (p Path) Start() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return p[0].Start()
}

// End implements the Linker interface.
func (p Path) End() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1].End()
}

// Revert implements the Linker interface.
func (p Path) Revert() {
	for incIdx, decIdx := 0, len(p)-1; incIdx < decIdx; incIdx, decIdx = incIdx+1, decIdx-1 {
		p[incIdx], p[decIdx] = p[decIdx], p[incIdx]
	}

	for idx := range p {
		p[idx].Revert()
	}
}

// Weight implements the Linker interface.
func (p Path) Weight(other Linker) [2]float64 {
	if len(p) == 0 {
		return [2]float64{0, 0}
	}

	return p[len(p)-1].Weight(other)
}

func buildPath(input []Linker) []Path {
	output := []Path{}
	var (
		curveList []Linker
	)

	rest := input

	for len(rest) > 0 {
		curveList, rest = SortEntities(rest, nil, func(from, to Linker) bool {
			return from.End().Equal(*to.Start()) || from.Start().Equal(*to.End())
		})

		path := Path{}

		for _, elt := range curveList {
			if value, ok := elt.(*Curve); ok {
				path = append(path, value)
			}

			if value, ok := elt.(*Segment); ok {
				path = append(path, value)
			}

			if value, ok := elt.(*Path); ok {
				path = append(path, value)
			}
		}

		output = append(output, path)
	}

	return output
}
