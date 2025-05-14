package geometry

import (
	"fmt"
	"sort"

	"github.com/yofu/dxf/entity"
)

type Point struct {
	Name string
	X    float64
	Y    float64

	linked    []*Point
	processed bool
}

type Points []Point

func (p Point) Weight(other *Point) float64 {
	return (p.X-other.X)*(p.X-other.X) + (p.Y-other.Y)*(p.Y-other.Y)
}

func (p Point) next() *Point {
	for idx := range p.linked {
		if !p.linked[idx].processed {
			p.linked[idx].processed = true

			return p.linked[idx]
		}
	}

	return nil
}

func FromDXF(input []*entity.Point) []Point {
	output := make([]Point, len(input)+1)

	for idx, dxfPoint := range input {
		output[idx+1] = Point{
			Name: fmt.Sprintf("#%d", idx),
			X:    dxfPoint.Coord[0],
			Y:    dxfPoint.Coord[1],
		}
	}

	for idx := range output {
		for index := range output {
			if output[idx].Weight(&output[index]) > 0 {
				output[idx].linked = append(output[idx].linked, &output[index])
			}
		}

		// Sort links by distance.
		sorter := newSorter(output[idx].linked, &output[idx])

		sort.Sort(sorter)

		output[idx].linked = sorter.data
	}

	return shorterPath(output)[1:]
}

func shorterPath(points []Point) []Point {
	points[0].processed = true
	currentPoint := &points[0]

	chain := []*Point{}

	for currentPoint != nil {
		chain = append(chain, currentPoint)
		currentPoint = currentPoint.next()
	}

	cpy := make([]Point, len(points))

	// copy values from pointers.
	for idx := range chain {
		cpy[idx] = *chain[idx]
	}

	return cpy
}
