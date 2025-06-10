package geometry

import (
	"fmt"
	"math"

	"github.com/yofu/dxf/entity"
)

type dxf struct {
	points     []*entity.Point
	arcs       []*entity.Arc
	circles    []*entity.Circle
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

func WithDXFCircle(data ...*entity.Circle) dxfConfigurator {
	return func(d *dxf) {
		d.circles = append(d.circles, data...)
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

	output := []Linker{}

	for _, dxfLine := range dxfFile.lines {
		output = append(output, &Segment{
			Name: fmt.Sprintf("#%d", len(output)),
			StartPoint: Coordinates{
				X: dxfLine.Start[0],
				Y: dxfLine.Start[1],
			},
			EndPoint: Coordinates{
				X: dxfLine.End[0],
				Y: dxfLine.End[1],
			},
		})
	}

	for _, dxfArc := range dxfFile.arcs {
		output = append(output, &Curve{
			Name: fmt.Sprintf("#%d", len(output)),
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
		})
	}

	for _, dxfCircle := range dxfFile.circles {
		output = append(output, &Path{
			&Curve{
				Name: fmt.Sprintf("#%d (1/2)", len(output)),
				Center: Coordinates{
					X: dxfCircle.Center[0],
					Y: dxfCircle.Center[1],
				},
				Radius: dxfCircle.Radius,
				StartPoint: Coordinates{
					X: dxfCircle.Center[0] + dxfCircle.Radius,
					Y: dxfCircle.Center[1],
				},
				EndPoint: Coordinates{
					X: dxfCircle.Center[0] - dxfCircle.Radius,
					Y: dxfCircle.Center[1],
				},
			},
			&Curve{
				Name: fmt.Sprintf("#%d (2/2)", len(output)),
				Center: Coordinates{
					X: dxfCircle.Center[0],
					Y: dxfCircle.Center[1],
				},
				Radius: dxfCircle.Radius,
				StartPoint: Coordinates{
					X: dxfCircle.Center[0] - dxfCircle.Radius,
					Y: dxfCircle.Center[1],
				},
				EndPoint: Coordinates{
					X: dxfCircle.Center[0] + dxfCircle.Radius,
					Y: dxfCircle.Center[1],
				},
			},
		})
	}

	for _, dxfPoly := range dxfFile.polyline {
		output = append(output, newPolyline(dxfPoly, fmt.Sprintf("#%d", len(output))))
	}

	for _, dxfPoly := range dxfFile.lwPolyline {
		output = append(output, newLightPolyline(dxfPoly, fmt.Sprintf("#%d", len(output))))
	}

	return buildPath(output)
}
