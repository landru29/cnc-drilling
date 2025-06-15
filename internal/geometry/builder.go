package geometry

import (
	"fmt"
)

func PathsFromDXF(entities ...dxfConfigurator) []Path {
	dxfFile := dxf{}

	for _, entitie := range entities {
		entitie(&dxfFile)
	}

	output := []Linker{}

	for _, dxfLine := range dxfFile.lines {
		output = append(output, &Segment{
			Name: fmt.Sprintf("#%d / Layer %s", len(output), dxfLine.Layer().Name()),
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
		output = append(output, NewCurveFromArc(fmt.Sprintf("#%d / layer %s", len(output), dxfArc.Layer().Name()), dxfArc))
	}

	for _, dxfCircle := range dxfFile.circles {
		output = append(output,
			NewPathFromCircle(fmt.Sprintf("#%d / Layer %s", len(output), dxfCircle.Layer().Name()), dxfCircle),
		)
	}

	for _, dxfPoly := range dxfFile.polyline {
		output = append(output, NewPathFromPolyline(fmt.Sprintf("#%d / Layer %s", len(output), dxfPoly.Layer().Name()), dxfPoly))
	}

	for _, dxfPoly := range dxfFile.lwPolyline {
		output = append(output, NewPathFromLightPolyline(fmt.Sprintf("#%d / Layer %s", len(output), dxfPoly.Layer().Name()), dxfPoly))
	}

	return buildPath(output)
}

func PointsFromDXFPoints(configs ...dxfConfigurator) []Point {
	dxfFile := dxf{}

	for _, config := range configs {
		config(&dxfFile)
	}

	inputPoints := make([]Linker, len(dxfFile.points))

	for idx, dxfPoint := range dxfFile.points {
		inputPoints[idx] = Point{
			Name:        fmt.Sprintf("#%d / Layer %s", idx, dxfPoint.Layer().Name()),
			Coordinates: NewCoordinatesFromPoint(dxfPoint),
		}
	}

	points, _ := SortEntities(inputPoints, &Coordinates{X: 0, Y: 0}, func(from, to Linker) bool {
		return true
	})

	output := []Point{}
	for _, point := range points {
		if value, ok := point.(Point); ok {
			output = append(output, value)
		}
	}

	return output
}
