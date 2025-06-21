package geometry

import "github.com/yofu/dxf/entity"

// Linker is a geometry glue.
type Linker interface {
	Start() *Coordinates
	End() *Coordinates
	Revert()
	Weight(Linker) [2]float64
	Box() Box
}

// NewLinker is a builder.
func NewLinker(name string, dxfEntity entity.Entity) Linker {
	switch data := dxfEntity.(type) {
	case *entity.Point:
		return NewPointFromPoint(name, data)
	case *entity.Vertex:
		return NewPointFromVertex(name, data)
	case *entity.Line:
		return NewSgmentFromLine(name, data)
	case *entity.Arc:
		return NewCurveFromArc(name, data)
	case *entity.Circle:
		return NewPathFromCircle(name, data)
	case *entity.Polyline:
		return NewPathFromPolyline(name, data)
	case *entity.LwPolyline:
		return NewPathFromLightPolyline(name, data)
	}

	return nil
}
