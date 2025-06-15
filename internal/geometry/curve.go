package geometry

import (
	"fmt"
	"math"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/yofu/dxf/entity"
)

type Curve struct {
	Name       string
	StartPoint Coordinates
	EndPoint   Coordinates
	Center     Coordinates
	Radius     float64
	Clockwise  bool
}

func NewCurveFromArc(name string, data *entity.Arc) *Curve {
	return &Curve{
		Name: name,
		Center: Coordinates{
			X: data.Center[0],
			Y: data.Center[1],
		},
		StartPoint: Coordinates{
			X: math.Cos(data.Angle[1]*math.Pi/180)*data.Radius + data.Center[0],
			Y: math.Sin(data.Angle[1]*math.Pi/180)*data.Radius + data.Center[1],
		},
		EndPoint: Coordinates{
			X: math.Cos(data.Angle[0]*math.Pi/180)*data.Radius + data.Center[0],
			Y: math.Sin(data.Angle[0]*math.Pi/180)*data.Radius + data.Center[1],
		},
		Clockwise: math.Mod((data.Angle[1]+360.0-data.Angle[0]), 360.0) < 0,
		Radius:    data.Radius,
	}
}

// Start implements the Linker interface.
func (c Curve) Start() *Coordinates {
	return &c.StartPoint
}

// End implements the Linker interface.
func (c Curve) End() *Coordinates {
	return &c.EndPoint
}

// Revert implements the Linker interface.
func (c *Curve) Revert() {
	c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
	c.Clockwise = !c.Clockwise
}

// Weight implements the Linker interface.
func (c Curve) Weight(other Linker) [2]float64 {
	return c.EndPoint.Weight(other)
}

func quarter(center Coordinates, point Coordinates) int {
	xSign := math.Signbit(point.X - center.X)
	ySign := math.Signbit(point.Y - center.Y)

	return map[bool]map[bool]int{
		false: {
			false: 1,
			true:  2,
		},
		true: {
			true:  3,
			false: 4,
		},
	}[xSign][ySign]
}

// Box implements the Linker interface.
func (c Curve) Box() Box {
	currentCurve := c

	if c.Clockwise {
		currentCurve = Curve{
			StartPoint: c.EndPoint,
			EndPoint:   c.StartPoint,
			Center:     c.Center,
			Radius:     c.Radius,
		}
	}

	startQuarter := quarter(currentCurve.Center, currentCurve.StartPoint)
	endQuarter := quarter(currentCurve.Center, currentCurve.EndPoint)

	maxX := math.Max(currentCurve.StartPoint.X, currentCurve.EndPoint.X)
	maxY := math.Max(currentCurve.StartPoint.Y, currentCurve.EndPoint.Y)
	minX := math.Min(currentCurve.StartPoint.X, currentCurve.EndPoint.X)
	minY := math.Min(currentCurve.StartPoint.Y, currentCurve.EndPoint.Y)

	if startQuarter == endQuarter {
		return Box{
			Min: Coordinates{
				X: minX,
				Y: minY,
			},
			Max: Coordinates{
				X: maxX,
				Y: maxY,
			},
		}
	}

	if startQuarter == 1 || endQuarter == 2 {
		maxX = currentCurve.Center.X + currentCurve.Radius
	}

	if startQuarter == 2 || endQuarter == 3 {
		minY = currentCurve.Center.Y - currentCurve.Radius
	}

	if startQuarter == 3 || endQuarter == 4 {
		minX = currentCurve.Center.X - currentCurve.Radius
	}

	if startQuarter == 4 || endQuarter == 1 {
		maxY = currentCurve.Center.Y + currentCurve.Radius
	}

	if (startQuarter == 1 || startQuarter == 2) && endQuarter == 4 {
		minY = currentCurve.Center.Y - currentCurve.Radius
	}

	return Box{
		Min: Coordinates{
			X: minX,
			Y: minY,
		},
		Max: Coordinates{
			X: maxX,
			Y: maxY,
		},
	}
}

// MarshallGCode implements the Marshaler interface.
func (c Curve) MarshallGCode(configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	output := ";------ Curve " + c.Name + "\n"

	if !options.IgnoreStart {
		start := c.Start()
		output = fmt.Sprintf(
			"G0 X%.01f Y%.01f\nG1 Z%.01f F%.01f ; Tool down\n",
			start.X-options.OffsetX(),
			start.Y-options.OffsetY(),
			-options.Deep,
			options.Feed,
		)
	}

	code := 2
	if c.Clockwise {
		code = 3
	}

	output += fmt.Sprintf(
		"G%d X%.01f Y%.01f I%.01f J%.01f F%.01f\n",
		code,
		c.EndPoint.X-options.OffsetX(),
		c.EndPoint.Y-options.OffsetY(),
		c.Center.X-c.StartPoint.X,
		c.Center.Y-c.StartPoint.Y,
		options.Feed,
	)

	return []byte(output), nil
}
