package engraver

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

type curve struct {
}

func Process(in io.Reader, out io.Writer, speedMillimeterPerMinute float64, drillingDeep float64, securityZ float64) error {
	drawing, err := dxf.FromReader(in)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(drawing)

	if _, err := fmt.Fprintf(out, "G90\nG0 Z%.01f\n", securityZ); err != nil {
		return err
	}

	arcs := []*entity.Arc{}
	lines := []*entity.Line{}

	for _, geometry := range drawing.Entities() {
		if arc, ok := geometry.(*entity.Arc); ok {
			arcs = append(arcs, arc)
		}

		if line, ok := geometry.(*entity.Line); ok {
			lines = append(lines, line)
		}
	}

	for idx, path := range geometry.CurvesFromDXF(lines, arcs) {
		// Start of path
		start := path.Start()
		if _, err := fmt.Fprintf(out, "; path #%d\nG0 X%.01f Y%.01f\nG1 Z%.01f F%.01f \n", idx, start.X, start.Y, drillingDeep, speedMillimeterPerMinute); err != nil {
			return err
		}

		for _, segment := range path {
			if segment.Radius == 0 {
				if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", segment.End.X, segment.End.Y, speedMillimeterPerMinute); err != nil {
					return err
				}
			} else {
				code := 3
				if segment.Clockwise {
					code = 2
				}

				if _, err := fmt.Fprintf(out, "G%d X%.01f Y%.01f I%.01f J%.01f F%.01f\n", code, segment.End.X, segment.End.Y, segment.Center.X, segment.Start.Y, speedMillimeterPerMinute); err != nil {
					return err
				}
			}
		}

		// End of path
		if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", securityZ); err != nil {
			return err
		}
	}

	return nil
}
