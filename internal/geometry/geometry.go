package geometry

import (
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

// WithDXFPoints is a configuration point.
func WithDXFPoints(data ...*entity.Point) dxfConfigurator {
	return func(d *dxf) {
		d.points = append(d.points, data...)
	}
}

// WithDXFArcs is a configuration point.
func WithDXFArcs(data ...*entity.Arc) dxfConfigurator {
	return func(d *dxf) {
		d.arcs = append(d.arcs, data...)
	}
}

// WithDXFCircle is a configuration point.
func WithDXFCircle(data ...*entity.Circle) dxfConfigurator {
	return func(d *dxf) {
		d.circles = append(d.circles, data...)
	}
}

// WithDXFLines is a configuration point.
func WithDXFLines(data ...*entity.Line) dxfConfigurator {
	return func(d *dxf) {
		d.lines = append(d.lines, data...)
	}
}

// WithDXFPolyline is a configuration point.
func WithDXFPolyline(data ...*entity.Polyline) dxfConfigurator {
	return func(d *dxf) {
		d.polyline = append(d.polyline, data...)
	}
}

// WithDXFLwPolyline is a configuration point.
func WithDXFLwPolyline(data ...*entity.LwPolyline) dxfConfigurator {
	return func(d *dxf) {
		d.lwPolyline = append(d.lwPolyline, data...)
	}
}
