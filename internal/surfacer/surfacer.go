package surfacer

import (
	"fmt"
	"io"

	"github.com/landru29/cnc-drilling/internal/configuration"
	"github.com/landru29/cnc-drilling/internal/geometry"
)

// Process is the surfacing process.
func Process(box geometry.Box, step float64, out io.Writer, config configuration.Config, method Method) error {
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

		switch method {
		case MethodZigzag:
			if err := surfaceAreaZigzag(box, step, out, config, deep); err != nil {
				return err
			}
		case MethodSpiral:
			if err := surfaceAreaSpiral(box, step, out, config, deep, true); err != nil {
				return err
			}
		case MethodSpiralInverted:
			if err := surfaceAreaSpiral(box, step, out, config, deep, false); err != nil {
				return err
			}
		case MethodSpiralFromCenter:
			if err := surfaceAreaSpiralFromCenter(box, step, out, config, deep, true); err != nil {
				return err
			}
		case MethodSpiralFromCenterInverted:
			if err := surfaceAreaSpiralFromCenter(box, step, out, config, deep, false); err != nil {
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

	return nil
}

func surfaceAreaZigzag(box geometry.Box, step float64, out io.Writer, config configuration.Config, deep float64) error {
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

func surfaceAreaSpiral(box geometry.Box, step float64, out io.Writer, config configuration.Config, deep float64, clockwise bool) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", box.Min.X, box.Min.Y); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G1 Z-%.01f F%.01f\n", deep, config.Feed); err != nil {
		return err
	}

	minX := box.Min.X
	maxX := box.Max.X
	minY := box.Min.Y
	maxY := box.Max.Y

	for minX < maxX && minY < maxY {
		if clockwise {
			// Move right
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, minY, config.Feed); err != nil {
				return err
			}
			minY += step

			// Move down
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, maxY, config.Feed); err != nil {
				return err
			}
			maxX -= step

			// Move left
			if minY < maxY {
				if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, maxY, config.Feed); err != nil {
					return err
				}
				maxY -= step
			}

			// Move up
			if minX < maxX {
				if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, minY, config.Feed); err != nil {
					return err
				}
				minX += step
			}
		} else {
			// Move left
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, minY, config.Feed); err != nil {
				return err
			}
			minY += step

			// Move down
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, maxY, config.Feed); err != nil {
				return err
			}
			maxX -= step

			// Move right
			if minY < maxY {
				if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, maxY, config.Feed); err != nil {
					return err
				}
				maxY -= step
			}

			// Move up
			if minX < maxX {
				if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, minY, config.Feed); err != nil {
					return err
				}
				minX += step
			}
		}
	}

	return nil
}

func surfaceAreaSpiralFromCenter(box geometry.Box, step float64, out io.Writer, config configuration.Config, deep float64, clockwise bool) error {
	if _, err := fmt.Fprintf(out, "G0 Z%.01f\n", config.SecurityZ); err != nil {
		return err
	}

	centerX := (box.Min.X + box.Max.X) / 2
	centerY := (box.Min.Y + box.Max.Y) / 2

	if _, err := fmt.Fprintf(out, "G0 X%.01f Y%.01f\n", centerX, centerY); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(out, "G1 Z-%.01f F%.01f\n", deep, config.Feed); err != nil {
		return err
	}

	minX := centerX
	maxX := centerX
	minY := centerY
	maxY := centerY

	for minX > box.Min.X || maxX < box.Max.X || minY > box.Min.Y || maxY < box.Max.Y {
		// Expand bounds
		if minX > box.Min.X {
			minX -= step
		}
		if maxX < box.Max.X {
			maxX += step
		}
		if minY > box.Min.Y {
			minY -= step
		}
		if maxY < box.Max.Y {
			maxY += step
		}

		if !clockwise {
			// Move right
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, minY, config.Feed); err != nil {
				return err
			}

			// Move down
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, maxY, config.Feed); err != nil {
				return err
			}

			// Move left
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, maxY, config.Feed); err != nil {
				return err
			}

			// Move up
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, minY, config.Feed); err != nil {
				return err
			}
		} else {
			// Move left
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, minY, config.Feed); err != nil {
				return err
			}

			// Move down
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", minX, maxY, config.Feed); err != nil {
				return err
			}

			// Move right
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, maxY, config.Feed); err != nil {
				return err
			}

			// Move up
			if _, err := fmt.Fprintf(out, "G1 X%.01f Y%.01f F%.01f\n", maxX, minY, config.Feed); err != nil {
				return err
			}
		}
	}

	return nil
}
