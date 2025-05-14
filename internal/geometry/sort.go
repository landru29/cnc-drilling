package geometry

type pointsPtr struct {
	data []*Point

	reference *Point
}

func newSorter(data []*Point, reference *Point) pointsPtr {
	return pointsPtr{
		data:      data,
		reference: reference,
	}
}

func (p pointsPtr) Len() int {
	return len(p.data)
}

func (p pointsPtr) Less(i, j int) bool {
	return p.reference.Weight(p.data[i].Coordinates) < p.reference.Weight(p.data[j].Coordinates)
}

func (p pointsPtr) Swap(i, j int) {
	p.data[i], p.data[j] = p.data[j], p.data[i]
}
