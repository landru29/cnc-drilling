package surfacer

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
)

// Process is the surfacing process.
func Process(box geometry.Box, step float64, out io.Writer, config configuration.Config) error {
	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.BeforeScript); err != nil {
		return err
	}

	tryDeeps := config.TryDeeps()

	for deepIndex, deep := range tryDeeps {
		if _, err := fmt.Fprintf(
			out,
			";\n;=== Surface %d/%d ===\n",
			deepIndex+1,
			len(tryDeeps),
		); err != nil {
			return err
		}

		if err := surfaceArea(box, step, out, config, deep); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.AfterScript); err != nil {
		return err
	}

	return nil
}

func surfaceArea(box geometry.Box, step float64, out io.Writer, config configuration.Config, deep float64) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", box.Min.X, box.Min.Y); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G1 Z-%.01f F%.01f\n", deep, config.Feed); err != nil {
		return err
	}

	positiveX := true
	for y := box.Min.Y; y <= box.Max.Y; y += step {
		if _, err := fmt.Fprintf(out, "G1 Y%.01f F%.01f\n", y, config.Feed); err != nil {
			return err
		}

		if positiveX {
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", box.Max.X, y, config.Feed); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", box.Min.X, y, config.Feed); err != nil {
				return err
			}
		}

		positiveX = !positiveX
	}

	return nil
}
