package surfacer

import "fmt"

// Method is the surfacing method.
type Method int

const (
	// MethodZigzag is the zigzag surfacing method.
	MethodZigzag Method = iota

	// MethodSpiral is the spiral surfacing method.
	MethodSpiral

	// MethodSpiralInverted is the spiral inverted surfacing method.
	MethodSpiralInverted

	// MethodSpiralFromCenter is the spiral surfacing method starting from the center.
	MethodSpiralFromCenter

	// MethodSpiralFromCenterInverted is the spiral inverted surfacing method starting from the center.
	MethodSpiralFromCenterInverted
)

// String implements the pflag.Value interface.
func (m Method) String() string {
	switch m {
	case MethodZigzag:
		return "zigzag"
	case MethodSpiral:
		return "spiral"
	case MethodSpiralInverted:
		return "spiral-inverted"
	case MethodSpiralFromCenter:
		return "spiral-from-center"
	case MethodSpiralFromCenterInverted:
		return "spiral-from-center-inverted"
	default:
		return "unknown"
	}
}

// Set implements the pflag.Value interface.
func (m *Method) Set(value string) error {
	switch value {
	case "zigzag":
		*m = MethodZigzag
	case "spiral":
		*m = MethodSpiral
	case "spiral-inverted":
		*m = MethodSpiralInverted
	case "spiral-from-center":
		*m = MethodSpiralFromCenter
	case "spiral-from-center-inverted":
		*m = MethodSpiralFromCenterInverted
	default:
		return fmt.Errorf("unknown surfacer method: %s", value)
	}
	return nil
}

// Type implements the pflag.Value interface.
func (m Method) Type() string {
	return "surfacerMethod"
}
