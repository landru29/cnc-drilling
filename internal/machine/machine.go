package machine

import (
	"fmt"
	"io"
	"math"
	"time"

	"github.com/landru29/cnc-drilling/internal/geometry"
)

// Path is a gcode path.
type Path struct {
	CurrentPosition geometry.Coordinates
	CurrentZ        float64
	Distance        float64
	Duration        time.Duration
}

// NewPath is a builder.
func NewPath(x float64, y float64, z float64) *Path {
	return &Path{
		CurrentPosition: geometry.Coordinates{X: x, Y: y},
		CurrentZ:        z,
		Distance:        0,
	}
}

func duration(distance float64, feed float64) time.Duration {
	return time.Duration((distance * float64(time.Minute)) / feed)
}

// MoveTo moves to a position.
func (p *Path) MoveTo(x float64, y float64, z float64, feed float64, out io.Writer) error {
	distance := math.Sqrt((p.CurrentZ-z)*(p.CurrentZ-z) + (p.CurrentPosition.X-x)*(p.CurrentPosition.X-x) + (p.CurrentPosition.Y-y)*(p.CurrentPosition.Y-y))
	p.Distance += distance
	p.CurrentPosition.X = x
	p.CurrentPosition.Y = y
	p.CurrentZ = z
	p.Duration += duration(distance, feed)

	if _, err := fmt.Fprintf(out, "G1 X%.3f Y%.3f Z%.3f F%.0f\n", x, y, z, feed); err != nil {
		return err
	}

	return nil
}

// MoveToXY moves to a XY position.
func (p *Path) MoveToXY(x float64, y float64, feed float64, out io.Writer) error {
	distance := math.Sqrt((p.CurrentPosition.X-x)*(p.CurrentPosition.X-x) + (p.CurrentPosition.Y-y)*(p.CurrentPosition.Y-y))
	p.Distance += distance
	p.CurrentPosition.X = x
	p.CurrentPosition.Y = y
	p.Duration += duration(distance, feed)

	if _, err := fmt.Fprintf(out, "G1 X%.3f Y%.3f F%.0f\n", x, y, feed); err != nil {
		return err
	}

	return nil
}

func (p *Path) MoveToZ(z float64, feed float64, out io.Writer) error {
	distance := math.Abs(p.CurrentZ - z)
	p.Distance += distance
	p.CurrentZ = z
	p.Duration += duration(distance, feed)

	if _, err := fmt.Fprintf(out, "G1 Z%.3f F%.0f\n", z, feed); err != nil {
		return err
	}

	return nil
}

// ArcTo creates an arc to a position.
func (p *Path) ArcTo(x float64, y float64, centerX float64, centerY float64, clockwise bool, feed float64, out io.Writer) error {
	code := 2
	if clockwise {
		code = 3
	}

	angleStart := math.Atan2(p.CurrentPosition.Y-centerY, p.CurrentPosition.X-centerX)
	angleEnd := math.Atan2(y-centerY, x-centerX)
	radius := math.Sqrt((p.CurrentPosition.X-centerX)*(p.CurrentPosition.X-centerX) + (p.CurrentPosition.Y-centerY)*(p.CurrentPosition.Y-centerY))

	angleDiff := angleEnd - angleStart
	if clockwise && angleDiff > 0 {
		angleDiff -= 2 * math.Pi
	} else if !clockwise && angleDiff < 0 {
		angleDiff += 2 * math.Pi
	}

	arcLength := math.Abs(angleDiff * radius)
	p.Distance += arcLength
	p.CurrentPosition.X = x
	p.CurrentPosition.Y = y
	p.Duration += duration(arcLength, feed)

	if _, err := fmt.Fprintf(
		out,
		"G%d X%.03f Y%.03f I%.03f J%.03f F%.03f\n",
		code,
		x,
		y,
		centerX-p.CurrentPosition.X,
		centerY-p.CurrentPosition.Y,
		feed,
	); err != nil {
		return err
	}

	return nil
}

// Retract moves the tool to the security Z height.
func (p *Path) Retract(securityZ float64, feed float64, out io.Writer) error {
	if p.CurrentZ < securityZ {
		distance := securityZ - p.CurrentZ
		p.Distance += distance
		p.CurrentZ = securityZ
		p.Duration += duration(distance, feed)

		if _, err := fmt.Fprintf(out, "G1 Z%.3f F%.0f\n", securityZ, feed); err != nil {
			return err
		}
	}

	return nil
}
