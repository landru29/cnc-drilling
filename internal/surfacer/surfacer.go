package surfacer

import (
	"fmt"
	"io"
	"time"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/machine"
)

// Process is the surfacing process.
func Process(box geometry.Box, step float64, out io.Writer, info io.Writer, config configuration.Config, method Method) error {
	if _, err := fmt.Fprintf(out, "G90\nG21\nG0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "%s\n", config.BeforeScript); err != nil {
		return err
	}

	tryDeeps := config.TryDeeps()

	var (
		distance float64       = 0
		duration time.Duration = 0
	)

	for deepIndex, deep := range tryDeeps {
		if _, err := fmt.Fprintf(
			out,
			";\n;=== Surface %d/%d ===\n",
			deepIndex+1,
			len(tryDeeps),
		); err != nil {
			return err
		}

		switch method {
		case MethodZigzag:
			if err := surfaceAreaZigzag(box, step, out, config, deep, &distance, &duration); err != nil {
				return err
			}
		case MethodSpiral:
			if err := surfaceAreaSpiral(box, step, out, config, deep, true, &distance, &duration); err != nil {
				return err
			}
		case MethodSpiralInverted:
			if err := surfaceAreaSpiral(box, step, out, config, deep, false, &distance, &duration); err != nil {
				return err
			}
		case MethodSpiralFromCenter:
			if err := surfaceAreaSpiralFromCenter(box, step, out, config, deep, true, &distance, &duration); err != nil {
				return err
			}
		case MethodSpiralFromCenterInverted:
			if err := surfaceAreaSpiralFromCenter(box, step, out, config, deep, false, &distance, &duration); err != nil {
				return err
			}
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

func surfaceAreaZigzag(
	box geometry.Box,
	step float64,
	out io.Writer,
	config configuration.Config,
	deep float64,
	distance *float64,
	duration *time.Duration,
) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", box.Min.X, box.Min.Y); err != nil {
		return err
	}

	path := machine.NewPath(box.Min.X, box.Min.Y, config.SecurityZ)

	if err := path.MoveToZ(-deep, config.Feed, out); err != nil {
		return err
	}

	positiveX := true
	for y := box.Min.Y; y <= box.Max.Y; y += step {
		if err := path.MoveToXY(path.CurrentPosition.X, y, config.Feed, out); err != nil {
			return err
		}

		if positiveX {
			if err := path.MoveToXY(box.Max.X, y, config.Feed, out); err != nil {
				return err
			}
		} else {
			if err := path.MoveToXY(box.Min.X, y, config.Feed, out); err != nil {
				return err
			}
		}

		positiveX = !positiveX
	}

	*distance += path.Distance
	*duration += path.Duration

	return nil
}

func surfaceAreaSpiral(
	box geometry.Box,
	step float64,
	out io.Writer,
	config configuration.Config,
	deep float64,
	clockwise bool,
	distance *float64,
	duration *time.Duration,
) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", box.Min.X, box.Min.Y); err != nil {
		return err
	}

	path := machine.NewPath(box.Min.X, box.Min.Y, config.SecurityZ)

	if err := path.MoveToZ(-deep, config.Feed, out); err != nil {
		return err
	}

	minX := box.Min.X
	maxX := box.Max.X
	minY := box.Min.Y
	maxY := box.Max.Y

	for minX < maxX && minY < maxY {
		if clockwise {
			// Move right
			if err := path.MoveToXY(maxX, minY, config.Feed, out); err != nil {
				return err
			}
			minY += step

			// Move down
			if err := path.MoveToXY(maxX, maxY, config.Feed, out); err != nil {
				return err
			}
			maxX -= step

			// Move left
			if minY < maxY {
				if err := path.MoveToXY(minX, maxY, config.Feed, out); err != nil {
					return err
				}
				maxY -= step
			}

			// Move up
			if minX < maxX {
				if err := path.MoveToXY(minX, minY, config.Feed, out); err != nil {
					return err
				}
				minX += step
			}
		} else {
			// Move left
			if err := path.MoveToXY(minX, minY, config.Feed, out); err != nil {
				return err
			}
			minY += step

			// Move down
			if err := path.MoveToXY(minX, maxY, config.Feed, out); err != nil {
				return err
			}
			maxX -= step

			// Move right
			if minY < maxY {
				if err := path.MoveToXY(maxX, maxY, config.Feed, out); err != nil {
					return err
				}
				maxY -= step
			}

			// Move up
			if minX < maxX {
				if err := path.MoveToXY(maxX, minY, config.Feed, out); err != nil {
					return err
				}
				minX += step
			}
		}
	}

	*distance += path.Distance
	*duration += path.Duration

	return nil
}

func surfaceAreaSpiralFromCenter(
	box geometry.Box,
	step float64,
	out io.Writer,
	config configuration.Config,
	deep float64,
	clockwise bool,
	distance *float64,
	duration *time.Duration,
) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	schema := geometry.Box{
		Min: geometry.Coordinates{X: (box.Min.X + box.Max.X) / 2, Y: (box.Min.Y + box.Max.Y) / 2},
		Max: geometry.Coordinates{X: (box.Min.X + box.Max.X) / 2, Y: (box.Min.Y + box.Max.Y) / 2},
	}

	diff := box.Width() - box.Height()
	if diff > 0 {
		schema = geometry.Box{
			Min: geometry.Coordinates{X: (box.Min.X+box.Max.X)/2 - diff/2, Y: (box.Min.Y + box.Max.Y) / 2},
			Max: geometry.Coordinates{X: (box.Min.X+box.Max.X)/2 + diff/2, Y: (box.Min.Y + box.Max.Y) / 2},
		}
	} else {
		schema = geometry.Box{
			Min: geometry.Coordinates{X: (box.Min.X + box.Max.X) / 2, Y: (box.Min.Y+box.Max.Y)/2 + diff/2},
			Max: geometry.Coordinates{X: (box.Min.X + box.Max.X) / 2, Y: (box.Min.Y+box.Max.Y)/2 - diff/2},
		}
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", schema.Min.X, schema.Min.Y); err != nil {
		return err
	}

	path := machine.NewPath(schema.Min.X, schema.Min.Y, config.SecurityZ)
	if err := path.MoveToZ(-deep, config.Feed, out); err != nil {
		return err
	}

	for {
		if err := drawBox(schema, out, config, clockwise, path); err != nil {
			return err
		}

		if schema.Min.X == box.Min.X && schema.Max.X == box.Max.X && schema.Min.Y == box.Min.Y && schema.Max.Y == box.Max.Y {
			break
		}

		schema.Min.X -= step
		schema.Min.Y -= step
		schema.Max.X += step
		schema.Max.Y += step

		if schema.Min.X < box.Min.X {
			schema.Min.X = box.Min.X
		}
		if schema.Min.Y < box.Min.Y {
			schema.Min.Y = box.Min.Y
		}
		if schema.Max.X > box.Max.X {
			schema.Max.X = box.Max.X
		}
		if schema.Max.Y > box.Max.Y {
			schema.Max.Y = box.Max.Y
		}
	}

	*distance += path.Distance
	*duration += path.Duration

	return nil
}

func drawBox(
	box geometry.Box,
	out io.Writer,
	config configuration.Config,
	clockwise bool,
	path *machine.Path,
) error {
	path.MoveToXY(box.Min.X, box.Min.Y, config.Feed, out)
	if !clockwise {
		// Move right
		if err := path.MoveToXY(box.Max.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}

		// Move down
		if err := path.MoveToXY(box.Max.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		// Move left
		if err := path.MoveToXY(box.Min.X, box.Max.Y, config.Feed, out); err != nil {
			return err
		}

		// Move up
		if err := path.MoveToXY(box.Min.X, box.Min.Y, config.Feed, out); err != nil {
			return err
		}
	} else {
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
	}

	return nil
}
