package geometry

import (
	"fmt"
	"math"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/yofu/dxf/entity"
)

type Curve struct {
	marker

	Name       string
	StartPoint Coordinates
	EndPoint   Coordinates
	Center     Coordinates
	Radius     float64
	Clockwise  bool
}

// Start implements the Linker interface.
func (c Curve) Start() *Coordinates {
	return &c.StartPoint
}

// End implements the Linker interface.
func (c Curve) End() *Coordinates {
	return &c.EndPoint
}

// Revert implements the Linker interface.
func (c *Curve) Revert() {
	c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
	c.Clockwise = !c.Clockwise
}

// Weight implements the Linker interface.
func (c Curve) Weight(other Linker) [2]float64 {
	return c.EndPoint.Weight(other)
}

// MarshallGCode implements the Marshaler interface.
func (c Curve) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	if c.Radius == 0 {
		return []byte(fmt.Sprintf(
			"G1 X%.01f Y%.01f F%.01f\n",
			c.EndPoint.X,
			c.EndPoint.Y,
			options.Feed,
		)), nil
	} else {
		code := 2
		if c.Clockwise {
			code = 3
		}

		return []byte(fmt.Sprintf(
			"G%d X%.01f Y%.01f I%.01f J%.01f F%.01f\n",
			code,
			c.EndPoint.X,
			c.EndPoint.Y,
			c.Center.X-c.StartPoint.X,
			c.Center.Y-c.StartPoint.Y,
			options.Feed,
		)), nil
	}
}

type Path []Curve

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
		out, err := segment.MarshallGCode(configs...)
		if err != nil {
			return nil, err
		}

		output += string(out)
	}

	return []byte(fmt.Sprintf("%sG0 Z%.01f\n", output, options.SecurityZ)), nil
}

func CurvesFromDXF(entities ...dxfConfigurator) []Path {
	dxfFile := dxf{}

	for _, entitie := range entities {
		entitie(&dxfFile)
	}

	output := make([]Linker, len(dxfFile.lines)+len(dxfFile.arcs)+len(dxfFile.polyline))

	for idx, dxfLine := range dxfFile.lines {
		output[idx] = &Curve{
			Name: fmt.Sprintf("#%d", idx),
			StartPoint: Coordinates{
				X: dxfLine.Start[0],
				Y: dxfLine.Start[1],
			},
			EndPoint: Coordinates{
				X: dxfLine.End[0],
				Y: dxfLine.End[1],
			},
		}
	}

	for idx, dxfArc := range dxfFile.arcs {
		output[idx+len(dxfFile.lines)] = &Curve{
			Name: fmt.Sprintf("#%d", idx+len(dxfFile.lines)),
			Center: Coordinates{
				X: dxfArc.Center[0],
				Y: dxfArc.Center[1],
			},
			StartPoint: Coordinates{
				X: math.Cos(dxfArc.Angle[0]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[0],
				Y: math.Sin(dxfArc.Angle[0]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[1],
			},
			EndPoint: Coordinates{
				X: math.Cos(dxfArc.Angle[1]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[0],
				Y: math.Sin(dxfArc.Angle[1]*math.Pi/180)*dxfArc.Radius + dxfArc.Center[1],
			},
			Clockwise: math.Mod((dxfArc.Angle[1]+360.0-dxfArc.Angle[0]), 360.0) > 0,
			Radius:    dxfArc.Radius,
		}
	}

	for idx, dxfArc := range dxfFile.polyline {
		output[idx+len(dxfFile.lines)+len(dxfFile.arcs)] = pathFromPolyline(dxfArc, fmt.Sprintf("#%d", idx))
	}

	return buildPath(output)
}

func pathFromPolyline(polyline *entity.Polyline, name string) Path {
	if polyline == nil || len(polyline.Vertices) < 2 {
		return nil
	}

	output := make(Path, len(polyline.Vertices)-1)
	currentVertex := polyline.Vertices[0]

	for idx, vertex := range polyline.Vertices[1:] {
		output[idx] = Curve{
			Name: fmt.Sprintf("%s #%d", name, idx),
			StartPoint: Coordinates{
				X: currentVertex.Coord[0],
				Y: currentVertex.Coord[0],
			},
			EndPoint: Coordinates{
				X: vertex.Coord[0],
				Y: vertex.Coord[0],
			},
		}

		currentVertex = vertex
	}

	return output
}

// Start implements the Linker interface.
func (p Path) Start() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return &p[0].StartPoint
}

// End implements the Linker interface.
func (p Path) End() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return &p[len(p)-1].EndPoint
}

// Revert implements the Linker interface.
func (p Path) Revert() {
	for incIdx, decIdx := 0, len(p)-1; incIdx < decIdx; incIdx, decIdx = incIdx+1, decIdx-1 {
		p[incIdx], p[decIdx] = p[decIdx], p[incIdx]
	}

	for idx := range p {
		p[idx].Clockwise = !p[idx].Clockwise
		p[idx].StartPoint, p[idx].EndPoint = p[idx].EndPoint, p[idx].StartPoint
	}
}

// Weight implements the Linker interface.
func (p Path) Weight(other Linker) [2]float64 {
	if len(p) == 0 {
		return [2]float64{0, 0}
	}

	return p[len(p)-1].Weight(other)
}

func (c *Curve) AddTo(path *Path) bool {
	end := path.End()
	start := path.Start()

	switch {
	case !c.available():
		return false

	case end.Equal(c.StartPoint):
		c.setUnavailable()
		*path = append(*path, *c)
		return true

	case end.Equal(c.EndPoint):
		c.setUnavailable()
		c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
		c.Clockwise = !c.Clockwise
		*path = append(*path, *c)
		return true

	case start.Equal(c.EndPoint):
		c.setUnavailable()
		*path = append(Path{*c}, *path...)
		return true

	case end.Equal(c.StartPoint):
		c.setUnavailable()
		c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
		c.Clockwise = !c.Clockwise
		*path = append(Path{*c}, *path...)
		return true
	}

	return false
}

func (p *Path) next(data []Curve) bool {
	for idx := range data {
		if data[idx].AddTo(p) {
			return true
		}
	}

	return false
}

func availableCurve(input []Curve) *Curve {
	for idx, curve := range input {
		if curve.available() {
			input[idx].setUnavailable()
			return &input[idx]
		}
	}

	return nil
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
				path = append(path, *value)
			}
		}

		output = append(output, path)
	}

	return output
}
