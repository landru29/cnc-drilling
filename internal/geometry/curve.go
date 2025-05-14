package geometry

import (
	"fmt"
	"math"

	"github.com/yofu/dxf/entity"
)

type Curve struct {
	Name      string
	Start     Coordinates
	End       Coordinates
	Center    Coordinates
	Radius    float64
	Clockwise bool

	processed bool
}

type Path []Curve

func CurvesFromDXF(lines []*entity.Line, arcs []*entity.Arc) []Path {
	output := make([]Curve, len(lines)+len(arcs))

	for idx, dxfLine := range lines {
		output[idx] = Curve{
			Name: fmt.Sprintf("#%d", idx),
			Start: Coordinates{
				X: dxfLine.Start[0],
				Y: dxfLine.Start[1],
			},
			End: Coordinates{
				X: dxfLine.End[0],
				Y: dxfLine.End[1],
			},
		}
	}

	for idx, dxfArc := range arcs {
		output[idx+len(lines)] = Curve{
			Name: fmt.Sprintf("#%d", idx+len(lines)),
			Center: Coordinates{
				X: dxfArc.Center[0],
				Y: dxfArc.Center[1],
			},
			Start: Coordinates{
				X: math.Cos(dxfArc.Angle[0]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[0],
				Y: math.Sin(dxfArc.Angle[0]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[1],
			},
			End: Coordinates{
				X: math.Cos(dxfArc.Angle[1]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[0],
				Y: math.Sin(dxfArc.Angle[1]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[1],
			},
			Clockwise: math.Mod((dxfArc.Angle[1]+360.0-dxfArc.Angle[0]), 360.0) > 0,
			Radius:    dxfArc.Radius,
		}
	}

	return buildPath(output)
}

func (p Path) Start() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return &p[0].Start
}

func (p Path) End() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return &p[len(p)-1].End
}

func (p *Path) next(data []Curve) bool {
	end := p.End()
	for idx := range data {
		if data[idx].processed {
			continue
		}

		if end.Equal(data[idx].Start) {
			data[idx].processed = true
			*p = append(*p, data[idx])
			return true
		}

		if end.Equal(data[idx].End) {
			data[idx].processed = true
			data[idx].Start, data[idx].End = data[idx].End, data[idx].Start
			data[idx].Clockwise = !data[idx].Clockwise
			*p = append(*p, data[idx])
			return true
		}
	}

	start := p.Start()
	for idx := range data {
		if data[idx].processed {
			continue
		}

		if start.Equal(data[idx].End) {
			data[idx].processed = true
			*p = append(Path{data[idx]}, *p...)
			return true
		}

		if end.Equal(data[idx].Start) {
			data[idx].processed = true
			data[idx].Start, data[idx].End = data[idx].End, data[idx].Start
			data[idx].Clockwise = !data[idx].Clockwise
			*p = append(Path{data[idx]}, *p...)
			return true
		}
	}

	return false
}

func availableCurve(input []Curve) *Curve {
	for idx, curve := range input {
		if !curve.processed {
			input[idx].processed = true
			return &input[idx]
		}
	}

	return nil
}

func buildPath(input []Curve) []Path {
	output := []Path{}

	for first := availableCurve(input); first != nil; first = availableCurve(input) {
		current := Path{*first}

		for current.next(input) {
		}

		output = append(output, current)
	}

	return output
}
