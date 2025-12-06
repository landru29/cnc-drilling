package outcutter

import (
	"fmt"
	"io"
	"time"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/machine"
)

// Process is the surfacing process.
func Process(box geometry.Box, out io.Writer, info io.Writer, config configuration.Config) error {
	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.BeforeScript); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.DeepZ); err != nil {
		return err
	}

	tryDeeps := config.TryDeepsXY()

	var (
		distance float64       = 0
		duration time.Duration = 0
	)

	for deepIndex, deep := range tryDeeps {
		cuttingBox := geometry.Box{
			Min: geometry.Coordinates{
				X: box.Min.X + deep,
				Y: box.Min.Y + deep,
			},
			Max: geometry.Coordinates{
				X: box.Max.X - deep,
				Y: box.Max.Y - deep,
			},
		}

		if _, err := fmt.Fprintf(
			out,
			";\n;=== Outcut %d/%d ===\n;New box: [(%.01f, %.01f), (%.01f, %.01f)] => [%.01f, %.01f]\n",
			deepIndex+1,
			len(tryDeeps),
			cuttingBox.Min.X,
			cuttingBox.Min.Y,
			cuttingBox.Max.X,
			cuttingBox.Max.Y,
			cuttingBox.Max.X-cuttingBox.Min.X,
			cuttingBox.Max.Y-cuttingBox.Min.Y,
		); err != nil {
			return err
		}

		if err := cutOut(cuttingBox, out, config, &distance, &duration); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\nZ%.01f\n", box.Min.X, box.Min.Y, config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.AfterScript); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(info, "; Total distance: %.01f mm\n; Total time: %s\n", distance, duration.Round(time.Second).String()); err != nil {
		return err
	}

	return nil
}

func cutOut(
	box geometry.Box,
	out io.Writer,
	config configuration.Config,
	distance *float64,
	duration *time.Duration,
) error {

	path := machine.NewPath(box.Min.X, box.Min.Y, config.DeepZ)

	// Move left
	if err := path.MoveToXY(box.Min.X, box.Max.Y, config.Feed, out); err != nil {
		return err
	}

	// Move down
	if err := path.MoveToXY(box.Max.X, box.Max.Y, config.Feed, out); err != nil {
		return err
	}

	// Move right
	if err := path.MoveToXY(box.Max.X, box.Min.Y, config.Feed, out); err != nil {
		return err
	}

	// Move up
	if err := path.MoveToXY(box.Min.X, box.Min.Y, config.Feed, out); err != nil {
		return err
	}

	*distance += path.Distance
	*duration += path.Duration

	return nil
}
