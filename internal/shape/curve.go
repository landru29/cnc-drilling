package shape

import (
	"fmt"
	"math"

	"github.com/landru29/cnc-drilling/internal/gcode"
	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/landru29/cnc-drilling/internal/machine"
	"github.com/yofu/dxf/entity"
)

// Curve is a curved segment.
type Curve struct {
	Name       string
	StartPoint geometry.CoordinatesXY
	EndPoint   geometry.CoordinatesXY
	Center     geometry.CoordinatesXY
	Radius     float64
	Clockwise  bool
}

// NewCurveFromArc is a builder.
func NewCurveFromArc(name string, data *entity.Arc) *Curve {
	return &Curve{
		Name: name,
		Center: geometry.CoordinatesXY{
			X: data.Center[0],
			Y: data.Center[1],
		},
		StartPoint: geometry.CoordinatesXY{
			X: math.Cos(data.Angle[1]*math.Pi/180)*data.Radius + data.Center[0],
			Y: math.Sin(data.Angle[1]*math.Pi/180)*data.Radius + data.Center[1],
		},
		EndPoint: geometry.CoordinatesXY{
			X: math.Cos(data.Angle[0]*math.Pi/180)*data.Radius + data.Center[0],
			Y: math.Sin(data.Angle[0]*math.Pi/180)*data.Radius + data.Center[1],
		},
		Clockwise: math.Mod((data.Angle[1]+360.0-data.Angle[0]), 360.0) < 0,
		Radius:    data.Radius,
	}
}

// Start implements the Linker interface.
func (c Curve) Start() *geometry.CoordinatesXY {
	return &c.StartPoint
}

// End implements the Linker interface.
func (c Curve) End() *geometry.CoordinatesXY {
	return &c.EndPoint
}

// Revert implements the Linker interface.
func (c *Curve) Revert() {
	c.StartPoint, c.EndPoint = c.EndPoint, c.StartPoint
	c.Clockwise = !c.Clockwise
}

// Weight implements the Linker interface.
func (c Curve) Weight(other Linker) [2]float64 {
	return [2]float64{
		c.EndPoint.Weight(*other.Start()),
		c.EndPoint.Weight(*other.End()),
	}
}

func quarter(center geometry.CoordinatesXY, point geometry.CoordinatesXY) int {
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
func (c Curve) Box() geometry.Box {
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
		return geometry.Box{
			Min: geometry.CoordinatesXY{
				X: minX,
				Y: minY,
			},
			Max: geometry.CoordinatesXY{
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

	return geometry.Box{
		Min: geometry.CoordinatesXY{
			X: minX,
			Y: minY,
		},
		Max: geometry.CoordinatesXY{
			X: maxX,
			Y: maxY,
		},
	}
}

// MarshallGCode implements the Marshaler interface.
func (c Curve) MarshallGCode(state *machine.Path, configs ...gcode.Configurator) ([]byte, error) {
	options := gcode.Options{}
	for _, config := range configs {
		config(&options)
	}

	output := ";------ Curve " + c.Name + "\n"

	if !options.IgnoreStart {
		start := c.Start()
		output = fmt.Sprintf(
			"G0 X%.03f Y%.03f\nG1 Z%.03f F%.03f ; Tool down\n",
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
		"G%d X%.03f Y%.03f I%.03f J%.03f F%.03f\n",
		code,
		c.EndPoint.X-options.OffsetX(),
		c.EndPoint.Y-options.OffsetY(),
		c.Center.X-c.StartPoint.X,
		c.Center.Y-c.StartPoint.Y,
		options.Feed,
	)

	return []byte(output), nil
}
