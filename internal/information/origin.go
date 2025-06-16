package information

import (
	"fmt"
	"strings"

	"github.com/landru29/cnc-drilling/internal/geometry"
)

type OriginDetection struct {
	Value    geometry.Coordinates
	Relative bool
}

// String implements the pflag.Value interface.
func (o OriginDetection) String() string {
	prefix := ""
	if o.Relative {
		prefix = "@"
	}

	return fmt.Sprintf("%s%.01f, %.01f", prefix, o.Value.X, o.Value.Y)

}

// Set implements the pflag.Value interface.
func (o *OriginDetection) Set(data string) error {
	if strings.HasPrefix(data, "@") {
		o.Relative = true

		data = data[1:]
	}

	return o.Value.Set(data)
}

// Type implements the pflag.Value interface.
func (c OriginDetection) Type() string {
	return "Origin"
}
