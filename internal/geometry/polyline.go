package geometry

import (
	"fmt"

	"github.com/yofu/dxf/entity"
)

func newPolyline(polyline *entity.Polyline, name string) *Path {
	if polyline == nil || len(polyline.Vertices) < 2 {
		return nil
	}

	output := make(Path, len(polyline.Vertices)-1)
	currentVertex := polyline.Vertices[0]

	for idx, vertex := range polyline.Vertices[1:] {
		center, ray := polyline.Bulge(idx)

		if center != nil && ray != 0 {
			output[idx] = &Curve{
				Name: fmt.Sprintf("%s #%d / Layer %s", name, idx, polyline.Layer().Name()),
				StartPoint: Coordinates{
					X: currentVertex.Coord[0],
					Y: currentVertex.Coord[1],
				},
				EndPoint: Coordinates{
					X: vertex.Coord[0],
					Y: vertex.Coord[1],
				},
				Center: Coordinates{
					X: center[0],
					Y: center[1],
				},
				Radius:    ray,
				Clockwise: vertex.Buldge > 0,
			}

			currentVertex = vertex

			continue
		}

		output[idx] = &Segment{
			Name: fmt.Sprintf("%s #%d", name, idx),
			StartPoint: Coordinates{
				X: currentVertex.Coord[0],
				Y: currentVertex.Coord[1],
			},
			EndPoint: Coordinates{
				X: vertex.Coord[0],
				Y: vertex.Coord[1],
			},
		}

		currentVertex = vertex
	}

	return &output
}

func newLightPolyline(polyline *entity.LwPolyline, name string) *Path {
	if polyline == nil || len(polyline.Vertices) < 2 {
		return nil
	}

	output := make(Path, len(polyline.Vertices)-1)
	currentVertex := polyline.Vertices[0]

	for idx, vertex := range polyline.Vertices[1:] {
		center, ray := polyline.Bulge(idx + 1)
		if center != nil && ray != 0 {
			output[idx] = &Curve{
				Name: fmt.Sprintf("%s #%d / Layer %s", name, idx, polyline.Layer().Name()),
				StartPoint: Coordinates{
					X: currentVertex[0],
					Y: currentVertex[1],
				},
				EndPoint: Coordinates{
					X: vertex[0],
					Y: vertex[1],
				},
				Center: Coordinates{
					X: center[0],
					Y: center[1],
				},
				Radius:    ray,
				Clockwise: polyline.Bulges[idx+1] > 0,
			}

			currentVertex = vertex

			continue
		}

		output[idx] = &Segment{
			Name: fmt.Sprintf("%s #%d", name, idx),
			StartPoint: Coordinates{
				X: currentVertex[0],
				Y: currentVertex[1],
			},
			EndPoint: Coordinates{
				X: vertex[0],
				Y: vertex[1],
			},
		}

		currentVertex = vertex
	}

	return &output
}
