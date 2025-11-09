package surfacer

import (
	"github.com/landru29/cnc-drilling/internal/geometry"
)

func translateTo(currentPosition *geometry.Coordinates, targetPosition geometry.Coordinates) float64 {
	moveDistance := currentPosition.DistanceTo(targetPosition)

	currentPosition.X = targetPosition.X
	currentPosition.Y = targetPosition.Y

	return moveDistance
}
