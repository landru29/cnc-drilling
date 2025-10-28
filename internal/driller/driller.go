package driller

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

// Process is the drilling process.
func Process(in io.Reader, out io.Writer, config configuration.Config) error {
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

	if _, err := fmt.Fprintf(out, "%s\n", config.BeforeScript); err != nil {
		return err
	}

	setOfPoints := []*entity.Point{}
	var shapeBox *geometry.Box

	for _, geometryElement := range geometry.FilterEntities(drawing.Entities(), config.Layers...) {
		if point, ok := geometryElement.(*entity.Point); ok {
			setOfPoints = append(setOfPoints, point)

			data := geometry.NewLinker("", geometryElement)
			if data == nil {
				continue
			}

			currentBox := data.Box()

			if shapeBox == nil {
				shapeBox = &currentBox

				continue
			}

			currentBox = currentBox.Merge(*shapeBox)
			shapeBox = &currentBox
		}
	}

	tryDeeps := config.TryDeeps()

	for deepIndex, deep := range tryDeeps {

		for idx, point := range geometry.PointsFromDXFPoints(geometry.WithDXFPoints(setOfPoints...)) {

			code, err := gcode.Marshal(
				point,
				gcode.WithDeep(deep),
				gcode.WithFeed(config.Feed),
				gcode.WithSecurityZ(config.SecurityZ),
				gcode.WithOffset(config.Origin.Computed(shapeBox)),
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

	if _, err := fmt.Fprintf(out, "%s\n", config.AfterScript); err != nil {
		return err
	}

	return nil
}
