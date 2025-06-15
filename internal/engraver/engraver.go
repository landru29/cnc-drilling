package engraver

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/information"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

func Process(in io.Reader, out io.Writer, config information.Config) error {
	drawing, err := dxf.FromReader(in)
	if err != nil {
		return err
	}

	defer func(closer io.Closer) {
		_ = closer.Close()
	}(drawing)

	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	arcs := []*entity.Arc{}
	lines := []*entity.Line{}
	lightPolylines := []*entity.LwPolyline{}
	polylines := []*entity.Polyline{}
	circles := []*entity.Circle{}

	for _, geometry := range geometry.FilterEntities(drawing.Entities(), config.Layers...) {
		if arc, ok := geometry.(*entity.Arc); ok {
			arcs = append(arcs, arc)
		}

		if line, ok := geometry.(*entity.Line); ok {
			lines = append(lines, line)
		}

		if lightPolyline, ok := geometry.(*entity.LwPolyline); ok {
			lightPolylines = append(lightPolylines, lightPolyline)
		}

		if polyline, ok := geometry.(*entity.Polyline); ok {
			polylines = append(polylines, polyline)
		}

		if circle, ok := geometry.(*entity.Circle); ok {
			circles = append(circles, circle)
		}
	}

	for idx, path := range geometry.PathsFromDXF(
		geometry.WithDXFLines(lines...),
		geometry.WithDXFArcs(arcs...),
		geometry.WithDXFLwPolyline(lightPolylines...),
		geometry.WithDXFPolyline(polylines...),
		geometry.WithDXFCircle(circles...),
	) {
		tryDeeps := config.TryDeeps()

		for deepIndex, deep := range tryDeeps {
			code, err := gcode.Marshal(
				path,
				gcode.WithDeep(deep),
				gcode.WithFeed(config.SpeedMillimeterPerMinute),
				gcode.WithSecurityZ(config.SecurityZ),
				gcode.WithOffset([]float64{config.Origin.X, config.Origin.Y}),
			)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(
				out,
				";\n;=== Path #%d %d/%d ===\n%s",
				idx,
				deepIndex+1,
				len(tryDeeps),
				string(code),
			); err != nil {
				return err
			}
		}
	}

	if _, err := fmt.Fprintf(out, "G0 X0 Y0\n"); err != nil {
		return err
	}

	return nil
}
