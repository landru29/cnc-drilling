package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/landru29/cnc-drilling/internal/geometry"
)

// OriginDetection is the tool origin coordinates. it can be absolute
// or relative of the cutting box.
type OriginDetection struct {
	Value    geometry.CoordinatesXY
	Z        float64
	Relative bool
}

// String implements the pflag.Value interface.
func (o OriginDetection) String() string {
	prefix := ""
	if o.Relative {
		prefix = "@"
	}

	return fmt.Sprintf("%s%.01f, %.01f, %.01f", prefix, o.Value.X, o.Value.Y, o.Z)
}

// Set implements the pflag.Value interface.
func (o *OriginDetection) Set(data string) error {
	if strings.HasPrefix(data, "@") {
		o.Relative = true

		data = data[1:]
	}

	splitter := strings.Split(data, ",")
	if len(splitter) < 2 {
		return errors.New("coordinates must be 0.0,0.0,0.0")
	}

	xValue, err := strconv.ParseFloat(strings.TrimSpace(splitter[0]), 64)
	if err != nil {
		return err
	}

	yValue, err := strconv.ParseFloat(strings.TrimSpace(splitter[1]), 64)
	if err != nil {
		return err
	}

	if len(splitter) > 2 {
		zValue, err := strconv.ParseFloat(strings.TrimSpace(splitter[2]), 64)
		if err != nil {
			return err
		}

		o.Z = zValue
	}

	o.Value.X = xValue
	o.Value.Y = yValue

	return nil
}

// Type implements the pflag.Value interface.
func (o OriginDetection) Type() string {
	return "Origin"
}

// Computed gets the real absolute coordinates of the origin.
func (o OriginDetection) Computed(shapeBox *geometry.Box) []float64 {
	if o.Relative && shapeBox != nil {
		return []float64{o.Value.X + shapeBox.Min.X, o.Value.Y + shapeBox.Min.Y}
	}

	return []float64{o.Value.X, o.Value.Y}
}

// UnmarshalJSON implements the JSON Unmarshaler interface.
func (o *OriginDetection) UnmarshalJSON(data []byte) error {
	var strData string

	if err := json.Unmarshal(data, &strData); err != nil {
		return err
	}

	return o.Set(strData)
}

// MarshalJSON implements the JSON Marshaler interface.
func (o OriginDetection) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// UnmarshalYAML implements the YAML Unmarshaler interface.
func (o *OriginDetection) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var strData string

	if err := unmarshal(&strData); err != nil {
		return err
	}

	return o.Set(strData)
}

// MarshalYAML implements the YAML Marshaler interface.
func (o OriginDetection) MarshalYAML() (any, error) {
	return o.String(), nil
}

// DecodeOrigin is the decoder form mapstructure.
func DecodeOrigin(f reflect.Type,
	t reflect.Type,
	data interface{},
) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	fieldName := t.Name()
	if fieldName != "OriginDetection" {
		return data, nil
	}

	output := OriginDetection{}
	if err := output.Set(data.(string)); err != nil {
		return nil, err
	}

	return output, nil
}
