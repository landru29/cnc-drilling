package edger

import (
	"fmt"
	"io"
	"time"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/machine"
)

func Process(box geometry.Box, out io.Writer, info io.Writer, config configuration.Config, clockwise bool) error {
	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.BeforeScript); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", box.Min.X, box.Min.Y); err != nil {
		return err
	}

	tryDeeps := config.TryDeepsZ()

	var (
		distance float64       = 0
		duration time.Duration = 0
	)

	for deepIndex, deep := range tryDeeps {
		if _, err := fmt.Fprintf(
			out,
			";\n;=== Edge %d/%d ===\n",
			deepIndex+1,
			len(tryDeeps),
		); err != nil {
			return err
		}

		if err := edgePerimeter(box, out, config, deep, clockwise, &distance, &duration); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
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

func edgePerimeter(box geometry.Box, out io.Writer, config configuration.Config, deep float64, clockwise bool, distance *float64, duration *time.Duration) error {
	path := machine.NewPath(box.Min.X, box.Min.Y, config.SecurityZ)

	if err := path.MoveToZ(-deep, config.Feed, out); err != nil {
		return err
	}

	if clockwise {

		if err := path.MoveToXY(box.Min.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		if err := path.MoveToXY(box.Max.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		if err := path.MoveToXY(box.Max.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}

		if err := path.MoveToXY(box.Min.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}
	} else {
		if err := path.MoveToXY(box.Min.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		if err := path.MoveToXY(box.Max.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		if err := path.MoveToXY(box.Max.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}
		if err := path.MoveToXY(box.Min.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	*distance += path.Distance
	*duration += path.Duration

	return nil
}
