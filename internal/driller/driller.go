package driller

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

	setOfPoints := []*entity.Point{}

	for _, geometry := range geometry.FilterEntities(drawing.Entities(), config.Layers...) {
		if point, ok := geometry.(*entity.Point); ok {
			setOfPoints = append(setOfPoints, point)
		}
	}

	for idx, point := range geometry.PointsFromDXFPoints(geometry.WithDXFPoints(setOfPoints...)) {
		tryDeeps := config.TryDeeps()

		for deepIndex, deep := range tryDeeps {
			code, err := gcode.Marshal(
				point,
				gcode.WithDeep(deep),
				gcode.WithFeed(config.SpeedMillimeterPerMinute),
				gcode.WithSecurityZ(config.SecurityZ),
			)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintf(
				out,
				";\n;=== Drilling #%d %d/%d ===\n%s",
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
