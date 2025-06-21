package configuration

import (
	"math"
)

// Config is the main application configuration.
type Config struct {
	Feed         float64         `default:"60"     json:"feed"           mapstructure:"feed"          yaml:"feed"`
	SecurityZ    float64         `default:"5"      json:"security_z"     mapstructure:"security_z"    yaml:"security_z"`
	Deepness     float64         `default:"1"      json:"deepness"       mapstructure:"deepness"      yaml:"deepness"`
	DeepPerTry   float64         `default:"0"      json:"deep_per_try"   mapstructure:"deep_per_try"  yaml:"deep_per_try"`
	Layers       []string        `                 json:"layers"         mapstructure:"layers"        yaml:"layers"`
	Origin       OriginDetection `                 json:"origin"         mapstructure:"origin"        yaml:"origin"`
	BeforeScript string          `default:""       json:"before_script"  mapstructure:"before_script" yaml:"before_script"`
	AfterScript  string          `default:"G0X0Y0" json:"after_script"   mapstructure:"after_script"  yaml:"after_script"`
}

// TryDeeps is the set of deeps during all tries.
func (c Config) TryDeeps() []float64 {
	if c.DeepPerTry <= 0 {
		return []float64{c.Deepness}
	}

	output := make([]float64, int(math.Ceil(c.Deepness/c.DeepPerTry)))

	maxFullTry := int(math.Floor(c.Deepness / c.DeepPerTry))

	output[len(output)-1] = math.Mod(c.Deepness, c.DeepPerTry) + float64(maxFullTry)*c.DeepPerTry

	for index := range maxFullTry {
		output[index] = float64(index+1) * c.DeepPerTry
	}

	return output
}
