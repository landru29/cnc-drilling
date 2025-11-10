package geometry

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// Box is the cutting box.
type Box struct {
	Min Coordinates
	Max Coordinates
}

// Merge is the Union of many boxes.
func (b Box) Merge(others ...Box) Box {
	output := b

	for _, box := range others {
		output = Box{
			Min: Coordinates{
				X: math.Min(output.Min.X, box.Min.X),
				Y: math.Min(output.Min.Y, box.Min.Y),
			},
			Max: Coordinates{
				X: math.Max(output.Max.X, box.Max.X),
				Y: math.Max(output.Max.Y, box.Max.Y),
			},
		}
	}

	return output
}

// String implements the Stringer interface.
func (b Box) String() string {
	return fmt.Sprintf(
		"[(%.03f, %.03f), (%.03f, %.03f)]",
		b.Min.X,
		b.Min.Y,
		b.Max.X,
		b.Max.Y,
	)
}

func (a *Box) Set(value string) error {
	cleanedValue := strings.ReplaceAll(value, " ", "")

	re := regexp.MustCompile(`\[\(([+-]?\d*(\.\d+)?),([+-]?\d*(\.\d+)?)\),\(([+-]?\d*(\.\d+)?),([+-]?\d*(\.\d+)?)\)\]`)
	if !re.MatchString(cleanedValue) {
		return fmt.Errorf("invalid box format: [(minX,minY),(maxX,maxY)]")
	}

	var x1, x2, y1, y2 float64

	if _, err := fmt.Sscanf(cleanedValue, "[(%f,%f),(%f,%f)]",
		&x1, &y1, &x2, &y2); err != nil {
		return err
	}

	a.Min = Coordinates{X: math.Min(x1, x2), Y: math.Min(y1, y2)}
	a.Max = Coordinates{X: math.Max(x1, x2), Y: math.Max(y1, y2)}
	return nil
}

// Type returns the type of the shape.
func (a Box) Type() string {
	return "box"
}

// Height returns the height of the box.
func (a Box) Height() float64 {
	return a.Max.Y - a.Min.Y
}

// Width returns the width of the box.
func (a Box) Width() float64 {
	return a.Max.X - a.Min.X
}
