package engraver

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/gcode"
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

	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", securityZ); err != nil {
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

	for idx, path := range geometry.CurvesFromDXF(geometry.WithDXFLines(lines...), geometry.WithDXFArcs(arcs...)) {
		code, err := gcode.Marshal(
			path,
			gcode.WithDeep(drillingDeep),
			gcode.WithFeed(speedMillimeterPerMinute),
			gcode.WithSecurityZ(securityZ),
		)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(
			out,
			";--- Path #%d ---\n%s",
			idx,
			string(code),
		); err != nil {
			return err
		}
	}

	return nil
}
