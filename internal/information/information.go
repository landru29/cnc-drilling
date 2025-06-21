package information

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

type counters struct {
	points         int
	lines          int
	circles        int
	ars            int
	polylines      int
	lightpolylines int
	vertices       int
}

func Process(in io.Reader, out io.Writer, config configuration.Config) error {
	drawing, err := dxf.FromReader(in)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(drawing)

	if _, err := fmt.Fprintf(out, "%d layer(s) found:\n", len(drawing.Layers)); err != nil {
		return err
	}

	layers := config.Layers
	if len(layers) == 0 {
		for _, layer := range drawing.Layers {
			layers = append(layers, layer.Name())
		}
	}

	for _, layer := range layers {
		isDefault := ""

		if drawing.CurrentLayer.Name() == layer {
			isDefault = " [default]"
		}

		if _, err := fmt.Fprintf(out, "\t* %s%s\n", layer, isDefault); err != nil {
			return err
		}

		var (
			box           *geometry.Box
			entityCounter counters
		)

		for idx, dxfEntity := range geometry.FilterEntities(drawing.Entities(), layer) {
			switch dxfEntity.(type) {
			case *entity.Point:
				entityCounter.points++
			case *entity.Vertex:
				entityCounter.vertices++
			case *entity.Line:
				entityCounter.lines++
			case *entity.Arc:
				entityCounter.ars++
			case *entity.Circle:
				entityCounter.circles++
			case *entity.Polyline:
				entityCounter.polylines++
			case *entity.LwPolyline:
				entityCounter.lightpolylines++
			}

			data := geometry.NewLinker(fmt.Sprintf("#%d", idx), dxfEntity)
			if data == nil {
				continue
			}

			currentBox := data.Box()

			if box == nil {
				box = &currentBox

				continue
			}

			currentBox = currentBox.Merge(*box)
			box = &currentBox
		}

		if entityCounter.points != 0 {
			if _, err := fmt.Fprintf(out, "\t\tPoints: %d\n", entityCounter.points); err != nil {
				return err
			}
		}

		if entityCounter.lines != 0 {
			if _, err := fmt.Fprintf(out, "\t\tLines: %d\n", entityCounter.lines); err != nil {
				return err
			}
		}

		if entityCounter.circles != 0 {
			if _, err := fmt.Fprintf(out, "\t\tCircles: %d\n", entityCounter.circles); err != nil {
				return err
			}
		}

		if entityCounter.ars != 0 {
			if _, err := fmt.Fprintf(out, "\t\tArcs: %d\n", entityCounter.ars); err != nil {
				return err
			}
		}

		if entityCounter.polylines != 0 {
			if _, err := fmt.Fprintf(out, "\t\tPolylines: %d\n", entityCounter.polylines); err != nil {
				return err
			}
		}

		if entityCounter.lightpolylines != 0 {
			if _, err := fmt.Fprintf(out, "\t\tLight polylines: %d\n", entityCounter.lightpolylines); err != nil {
				return err
			}
		}

		if entityCounter.vertices != 0 {
			if _, err := fmt.Fprintf(out, "\t\tVertices: %d\n", entityCounter.vertices); err != nil {
				return err
			}
		}

		if box != nil {
			if _, err := fmt.Fprintf(out, "\t\tBox %s\n", box); err != nil {
				return err
			}
		}
	}

	return nil
}
