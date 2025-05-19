package geometry

import (
	"fmt"

	"github.com/yofu/dxf/entity"
)

type Polyline []Segment

func (p Polyline) Start() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return p[0].Start()
}

func (p Polyline) End() *Coordinates {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1].End()
}

func (p *Polyline) Revert() {
	for incIdx, decIdx := 0, len(*p)-1; incIdx < decIdx; incIdx, decIdx = incIdx+1, decIdx-1 {
		(*p)[incIdx], (*p)[decIdx] = (*p)[decIdx], (*p)[incIdx]
	}

	for idx := range *p {
		(*p)[idx].Revert()
	}
}

func (p Polyline) Weight(other Linker) [2]float64 {
	if len(p) == 0 {
		return [2]float64{0, 0}
	}

	return p.End().Weight(other)
}

func newPolyline(polyline *entity.Polyline, name string) *Polyline {
	if polyline == nil || len(polyline.Vertices) < 2 {
		return nil
	}

	output := make(Polyline, len(polyline.Vertices)-1)
	currentVertex := polyline.Vertices[0]

	for idx, vertex := range polyline.Vertices[1:] {
		output[idx] = Segment{
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

	return &output
}

func newLightPolyline(polyline *entity.LwPolyline, name string) *Polyline {
	if polyline == nil || len(polyline.Vertices) < 2 {
		return nil
	}

	output := make(Polyline, len(polyline.Vertices)-1)
	currentVertex := polyline.Vertices[0]

	for idx, vertex := range polyline.Vertices[1:] {
		output[idx] = Segment{
			Name: fmt.Sprintf("%s #%d", name, idx),
			StartPoint: Coordinates{
				X: currentVertex[0],
				Y: currentVertex[0],
			},
			EndPoint: Coordinates{
				X: vertex[0],
				Y: vertex[0],
			},
		}

		currentVertex = vertex
	}

	return &output
}
