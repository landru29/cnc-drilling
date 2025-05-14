package driller

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

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

	setOfPoints := []*entity.Point{}

	for _, geometry := range drawing.Entities() {
		if point, ok := geometry.(*entity.Point); ok {
			setOfPoints = append(setOfPoints, point)
		}
	}

	for idx, point := range geometry.FromDXF(setOfPoints) {
		if _, err := fmt.Fprintf(
			out,
			";--- Drilling #%d ---\nG0 X%.01F Y%.01f\nG1 Z%.01f F%.01f\nG0 Z%.01f\n",
			idx,
			point.X, point.Y,
			-drillingDeep,
			speedMillimeterPerMinute,
			securityZ,
		); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(out, "G0 X0 Y0\n"); err != nil {
		return err
	}

	return nil
}
