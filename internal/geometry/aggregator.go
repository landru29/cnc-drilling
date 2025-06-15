package geometry

import (
	"math"
	"sort"
)

type Linker interface {
	Start() *Coordinates
	End() *Coordinates
	Revert()
	Weight(Linker) [2]float64
	Box() Box
}

type LinkerSort struct {
	data []Linker

	reference Linker
}

func (p LinkerSort) Len() int {
	return len(p.data)
}

func (p LinkerSort) Less(i, j int) bool {
	weightI := p.reference.Weight(p.data[i])
	weightJ := p.reference.Weight(p.data[j])

	return math.Min(weightI[0], weightI[1]) < math.Min(weightJ[0], weightJ[1])
}

func (p *LinkerSort) Swap(i, j int) {
	p.data[i], p.data[j] = p.data[j], p.data[i]
}

func nextEntity(entities []Linker, from Linker, filter func(from Linker, to Linker) bool) (Linker, []Linker, float64) {
	if len(entities) == 0 {
		return nil, nil, 0
	}

	sorter := LinkerSort{
		data:      make([]Linker, len(entities)),
		reference: from,
	}

	copy(sorter.data, entities)

	// for idx, entity := range entities {
	// 	sorter.data[idx] = entity
	// }

	sort.Sort(&sorter)

	for idx := range sorter.data {
		if filter(from, sorter.data[idx]) {
			weight := from.Weight(sorter.data[0])

			return sorter.data[0], sorter.data[1:], math.Min(weight[0], weight[1])
		}
	}

	return nil, sorter.data, 0
}

func SortEntities(entities []Linker, from *Coordinates, filter func(from Linker, to Linker) bool) ([]Linker, []Linker) {
	var (
		end    *Coordinates
		start  *Coordinates
		output []Linker

		current Linker
	)

	linkers := make([]Linker, len(entities))

	copy(linkers, entities)

	if from != nil {
		end = from
		start = from
		current = from
	} else {
		current = linkers[0]
		end = current.End()
		start = current.Start()
		output = append(output, current)
		linkers = linkers[1:]
	}

	for current != nil {
		after, linkersAfter, weightAfter := nextEntity(linkers, end, filter)
		before, linkersBefore, weightBefore := nextEntity(linkers, start, filter)

		switch {
		case before == nil && after == nil:
			current = nil
			continue

		case from != nil:
			current = after

			weight := end.Weight(current)

			if weight[0] < weight[1] {
				current.Revert()
			}

			output = append(output, current)

			linkers = linkersAfter

		case weightAfter > weightBefore || after == nil:
			current = before

			weight := start.Weight(current)

			if weight[1] > weight[0] {
				current.Revert()
			}

			output = append([]Linker{current}, output...)

			linkers = linkersBefore

		case weightAfter <= weightBefore || before == nil:
			current = after

			weight := end.Weight(current)

			if weight[0] > weight[1] {
				current.Revert()
			}

			output = append(output, current)

			linkers = linkersAfter
		}

		end = output[len(output)-1].End()
		start = output[0].Start()
	}

	return output, linkers
}
