package configuration

import (
	"math"
)

// Config is the main application configuration.
type Config struct {
	Feed         float64         `default:"60"     json:"feed"            mapstructure:"feed"            yaml:"feed"`
	SecurityZ    float64         `default:"5"      json:"security_z"      mapstructure:"security_z"      yaml:"security_z"`
	DeepZ        float64         `default:"1"      json:"deep_z"          mapstructure:"deep_z"          yaml:"deep_z"`
	DeepZPerTry  float64         `default:"0"      json:"deep_z_per_try"  mapstructure:"deep_z_per_try"  yaml:"deep_z_per_try"`
	DeepXY       float64         `default:"0"      json:"deep_xy"         mapstructure:"deep_xy"         yaml:"deep_xy"`
	DeepXYPerTry float64         `default:"0"      json:"deep_xy_per_try" mapstructure:"deep_xy_per_try" yaml:"deep_xy_per_try"`
	Layers       []string        `                 json:"layers"          mapstructure:"layers"          yaml:"layers"`
	Origin       OriginDetection `                 json:"origin"          mapstructure:"origin"          yaml:"origin"`
	BeforeScript string          `default:""       json:"before_script"   mapstructure:"before_script"   yaml:"before_script"`
	AfterScript  string          `default:"G0X0Y0" json:"after_script"    mapstructure:"after_script"    yaml:"after_script"`
}

// TryDeepsZ is the set of deeps on Z axis during all tries.
func (c Config) TryDeepsZ() []float64 {
	if c.DeepZPerTry <= 0 {
		return []float64{c.DeepZ}
	}
	output := make([]float64, int(math.Ceil((c.DeepZ+c.Origin.Z)/c.DeepZPerTry)))

	maxFullTry := int(math.Floor((c.DeepZ + c.Origin.Z) / c.DeepZPerTry))

	output[len(output)-1] = math.Mod(c.DeepZ+c.Origin.Z, c.DeepZPerTry) + float64(maxFullTry)*c.DeepZPerTry - c.Origin.Z

	for index := range maxFullTry {
		output[index] = float64(index+1)*c.DeepZPerTry - c.Origin.Z
	}

	return output
}

// TryDeepsXY is the set of deeps on XY axis during all tries.
func (c Config) TryDeepsXY() []float64 {
	if c.DeepXYPerTry <= 0 {
		return []float64{c.DeepXY}
	}

	output := make([]float64, int(math.Ceil((c.DeepXY)/c.DeepXYPerTry)))

	maxFullTry := int(math.Floor((c.DeepXY) / c.DeepXYPerTry))

	output[len(output)-1] = math.Mod(c.DeepXY, c.DeepXYPerTry) + float64(maxFullTry)*c.DeepXYPerTry

	for index := range maxFullTry {
		output[index] = float64(index+1) * c.DeepXYPerTry
	}

	return output
}
