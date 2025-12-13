package surfacer

import (
	"github.com/landru29/cnc-drilling/internal/geometry"
)

func translateTo(currentPosition *geometry.CoordinatesXY, targetPosition geometry.CoordinatesXY) float64 {
	moveDistance := currentPosition.DistanceTo(targetPosition)

	currentPosition.X = targetPosition.X
	currentPosition.Y = targetPosition.Y

	return moveDistance
}
