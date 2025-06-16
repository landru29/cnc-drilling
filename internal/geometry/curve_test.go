package geometry_test

import (
	"testing"

	"github.com/landru29/cnc-drilling/internal/geometry"
	"github.com/stretchr/testify/assert"
)

func TestCurveBox(t *testing.T) {
	t.Run("first quarter", func(t *testing.T) {
		curve := geometry.Curve{
			StartPoint: geometry.Coordinates{
				X: 20,
				Y: 30,
			},
			EndPoint: geometry.Coordinates{
				X: 30,
				Y: 20,
			},
			Center: geometry.Coordinates{
				X: 20,
				Y: 20,
			},
			Radius: 10,
		}

		assert.Equal(
			t,
			geometry.Box{
				Min: geometry.Coordinates{
					X: 20,
					Y: 20,
				},
				Max: geometry.Coordinates{
					X: 30,
					Y: 30,
				},
			},
			curve.Box(),
		)
	})

	t.Run("2 first quarters", func(t *testing.T) {
		curve := geometry.Curve{
			StartPoint: geometry.Coordinates{
				X: 20,
				Y: 30,
			},
			EndPoint: geometry.Coordinates{
				X: 27,
				Y: 13,
			},
			Center: geometry.Coordinates{
				X: 20,
				Y: 20,
			},
			Radius: 10,
		}

		assert.Equal(
			t,
			geometry.Box{
				Min: geometry.Coordinates{
					X: 20,
					Y: 13,
				},
				Max: geometry.Coordinates{
					X: 30,
					Y: 30,
				},
			},
			curve.Box(),
		)
	})

	t.Run("3 first quarters", func(t *testing.T) {
		curve := geometry.Curve{
			StartPoint: geometry.Coordinates{
				X: 20,
				Y: 30,
			},
			EndPoint: geometry.Coordinates{
				X: 13,
				Y: 13,
			},
			Center: geometry.Coordinates{
				X: 20,
				Y: 20,
			},
			Radius: 10,
		}

		assert.Equal(
			t,
			geometry.Box{
				Min: geometry.Coordinates{
					X: 13,
					Y: 10,
				},
				Max: geometry.Coordinates{
					X: 30,
					Y: 30,
				},
			},
			curve.Box(),
		)
	})

	t.Run("4 first quarters", func(t *testing.T) {
		curve := geometry.Curve{
			StartPoint: geometry.Coordinates{
				X: 20,
				Y: 30,
			},
			EndPoint: geometry.Coordinates{
				X: 13,
				Y: 27,
			},
			Center: geometry.Coordinates{
				X: 20,
				Y: 20,
			},
			Radius: 10,
		}

		assert.Equal(
			t,
			geometry.Box{
				Min: geometry.Coordinates{
					X: 10,
					Y: 10,
				},
				Max: geometry.Coordinates{
					X: 30,
					Y: 30,
				},
			},
			curve.Box(),
		)
	})
}
