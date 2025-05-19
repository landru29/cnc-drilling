package geometry_test

import (
	"testing"

	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/stretchr/testify/assert"
)

func TestPolylineRevert(t *testing.T) {
	t.Run("odd", func(t *testing.T) {
		polyline := geometry.Polyline{
			{
				StartPoint: geometry.Coordinates{
					X: 0,
					Y: 0,
				},
				EndPoint: geometry.Coordinates{
					X: 2,
					Y: 1,
				},
			},
			{
				StartPoint: geometry.Coordinates{
					X: 2,
					Y: 1,
				},
				EndPoint: geometry.Coordinates{
					X: 5,
					Y: 2,
				},
			},
			{
				StartPoint: geometry.Coordinates{
					X: 5,
					Y: 2,
				},
				EndPoint: geometry.Coordinates{
					X: 2,
					Y: 3,
				},
			},
		}

		polyline.Revert()

		assert.Equal(t,
			geometry.Polyline{
				{
					StartPoint: geometry.Coordinates{
						X: 2,
						Y: 3,
					},
					EndPoint: geometry.Coordinates{
						X: 5,
						Y: 2,
					},
				},
				{
					StartPoint: geometry.Coordinates{
						X: 5,
						Y: 2,
					},
					EndPoint: geometry.Coordinates{
						X: 2,
						Y: 1,
					},
				},
				{
					StartPoint: geometry.Coordinates{
						X: 2,
						Y: 1,
					},
					EndPoint: geometry.Coordinates{
						X: 0,
						Y: 0,
					},
				},
			},
			polyline,
		)
	})

	t.Run("even", func(t *testing.T) {
		polyline := geometry.Polyline{
			{
				StartPoint: geometry.Coordinates{
					X: 0,
					Y: 0,
				},
				EndPoint: geometry.Coordinates{
					X: 2,
					Y: 1,
				},
			},
			{
				StartPoint: geometry.Coordinates{
					X: 2,
					Y: 1,
				},
				EndPoint: geometry.Coordinates{
					X: 5,
					Y: 2,
				},
			},
			{
				StartPoint: geometry.Coordinates{
					X: 5,
					Y: 2,
				},
				EndPoint: geometry.Coordinates{
					X: 2,
					Y: 3,
				},
			},
			{
				StartPoint: geometry.Coordinates{
					X: 2,
					Y: 3,
				},
				EndPoint: geometry.Coordinates{
					X: 1,
					Y: 1,
				},
			},
		}

		polyline.Revert()

		assert.Equal(t,
			geometry.Polyline{
				{
					StartPoint: geometry.Coordinates{
						X: 1,
						Y: 1,
					},
					EndPoint: geometry.Coordinates{
						X: 2,
						Y: 3,
					},
				},
				{
					StartPoint: geometry.Coordinates{
						X: 2,
						Y: 3,
					},
					EndPoint: geometry.Coordinates{
						X: 5,
						Y: 2,
					},
				},
				{
					StartPoint: geometry.Coordinates{
						X: 5,
						Y: 2,
					},
					EndPoint: geometry.Coordinates{
						X: 2,
						Y: 1,
					},
				},
				{
					StartPoint: geometry.Coordinates{
						X: 2,
						Y: 1,
					},
					EndPoint: geometry.Coordinates{
						X: 0,
						Y: 0,
					},
				},
			},
			polyline,
		)
	})
}
