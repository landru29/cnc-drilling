package engraver

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/yofu/dxf"
	"github.com/yofu/dxf/entity"
)

// Process is the engraving process.
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

	arcs := []*entity.Arc{}
	lines := []*entity.Line{}
	lightPolylines := []*entity.LwPolyline{}
	polylines := []*entity.Polyline{}
	circles := []*entity.Circle{}

	var shapeBox *geometry.Box

	for _, geometryElement := range geometry.FilterEntities(drawing.Entities(), config.Layers...) {
		if arc, ok := geometryElement.(*entity.Arc); ok {
			arcs = append(arcs, arc)
		}

		if line, ok := geometryElement.(*entity.Line); ok {
			lines = append(lines, line)
		}

		if lightPolyline, ok := geometryElement.(*entity.LwPolyline); ok {
			lightPolylines = append(lightPolylines, lightPolyline)
		}

		if polyline, ok := geometryElement.(*entity.Polyline); ok {
			polylines = append(polylines, polyline)
		}

		if circle, ok := geometryElement.(*entity.Circle); ok {
			circles = append(circles, circle)
		}

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
				gcode.WithFeed(config.Feed),
				gcode.WithSecurityZ(config.SecurityZ),
				gcode.WithOffset(config.Origin.Computed(shapeBox)),
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

	if _, err := fmt.Fprintf(out, "%s\n", config.AfterScript); err != nil {
		return err
	}

	return nil
}
