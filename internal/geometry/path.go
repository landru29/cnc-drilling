package geometry

import (
	"fmt"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/yofu/dxf/entity"
)

type Path []Linker

func NewPathFromCircle(name string, data *entity.Circle) *Path {
	return &Path{
		&Curve{
			Name: fmt.Sprintf("%s (1/2)", name),
			Center: Coordinates{
				X: data.Center[0],
				Y: data.Center[1],
			},
			Radius: data.Radius,
			StartPoint: Coordinates{
				X: data.Center[0] + data.Radius,
				Y: data.Center[1],
			},
			EndPoint: Coordinates{
				X: data.Center[0] - data.Radius,
				Y: data.Center[1],
			},
		},
		&Curve{
			Name: fmt.Sprintf("%s (2/2)", name),
			Center: Coordinates{
				X: data.Center[0],
				Y: data.Center[1],
			},
			Radius: data.Radius,
			StartPoint: Coordinates{
				X: data.Center[0] - data.Radius,
				Y: data.Center[1],
			},
			EndPoint: Coordinates{
				X: data.Center[0] + data.Radius,
				Y: data.Center[1],
			},
		},
	}
}

// MarshallGCode implements the Marshaler interface.
func (p Path) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	var output string

	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	if !options.IgnoreStart {
		start := p.Start()
		output = fmt.Sprintf(
			"G0 X%.03f Y%.03f\nG1 Z%.03f F%.03f; Tool down\n",
			start.X-options.OffsetX(),
			start.Y-options.OffsetY(),
			-options.Deep,
			options.Feed,
		)
	}

	for _, segmentOrCurve := range p {
		localConf := append([]gcode.Configurator{gcode.WithoutStart(), gcode.WithoutEnd()}, configs...)

		out, err := gcode.Marshal(segmentOrCurve, localConf...)
		if err != nil {
			return nil, err
		}

		output += string(out)
	}

	if !options.IgnoreEnd {
		output += fmt.Sprintf("G0 Z%.03f; Tool up\n", options.SecurityZ)
	}

	return []byte(output), nil
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

// Box implements the Linker interface.
func (p Path) Box() Box {
	if len(p) == 0 {
		return Box{}
	}

	output := p[0].Box()

	for _, elt := range p[1:] {
		output = output.Merge(elt.Box())
	}

	return output
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
