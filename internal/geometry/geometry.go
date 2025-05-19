package geometry

import (
	"fmt"
	"math"

	"github.com/yofu/dxf/entity"
)

type dxf struct {
	points     []*entity.Point
	arcs       []*entity.Arc
	lines      []*entity.Line
	polyline   []*entity.Polyline
	lwPolyline []*entity.LwPolyline
}

type dxfConfigurator func(*dxf)

func WithDXFPoints(data ...*entity.Point) dxfConfigurator {
	return func(d *dxf) {
		d.points = append(d.points, data...)
	}
}

func WithDXFArcs(data ...*entity.Arc) dxfConfigurator {
	return func(d *dxf) {
		d.arcs = append(d.arcs, data...)
	}
}

func WithDXFLines(data ...*entity.Line) dxfConfigurator {
	return func(d *dxf) {
		d.lines = append(d.lines, data...)
	}
}

func WithDXFPolyline(data ...*entity.Polyline) dxfConfigurator {
	return func(d *dxf) {
		d.polyline = append(d.polyline, data...)
	}
}

func WithDXFLwPolyline(data ...*entity.LwPolyline) dxfConfigurator {
	return func(d *dxf) {
		d.lwPolyline = append(d.lwPolyline, data...)
	}
}

func PathsFromDXF(entities ...dxfConfigurator) []Path {
	dxfFile := dxf{}

	for _, entitie := range entities {
		entitie(&dxfFile)
	}

	output := make([]Linker, len(dxfFile.lines)+len(dxfFile.arcs)+len(dxfFile.polyline))

	for idx, dxfLine := range dxfFile.lines {
		output[idx] = &Segment{
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

	for idx, dxfPoly := range dxfFile.polyline {
		output[idx+len(dxfFile.lines)+len(dxfFile.arcs)] = newPolyline(dxfPoly, fmt.Sprintf("#%d", idx))
	}

	for idx, dxfPoly := range dxfFile.lwPolyline {
		output[idx+len(dxfFile.lines)+len(dxfFile.arcs)+len(dxfFile.polyline)] = newLightPolyline(dxfPoly, fmt.Sprintf("#%d", idx))
	}

	return buildPath(output)
}
