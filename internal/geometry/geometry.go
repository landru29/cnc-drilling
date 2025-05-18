package geometry

import "github.com/yofu/dxf/entity"

type marker struct {
	processed bool
}

func (m marker) available() bool {
	return !m.processed
}

func (m *marker) setUnavailable() {
	m.processed = true
}

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
