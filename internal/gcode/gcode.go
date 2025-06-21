package gcode

import (
	"fmt"
	"reflect"
)

// Marshaler is the gcode marshaler.
type Marshaler interface {
	MarshallGCode(configs ...Configurator) ([]byte, error)
}

// Configurator is the marshaler configuration.
type Configurator func(*Options)

type Options struct {
	Deep        float64
	Feed        float64
	SecurityZ   float64
	IgnoreStart bool
	IgnoreEnd   bool
	Offset      []float64
}

// WithDeep is a configuration point.
func WithDeep(deep float64) Configurator {
	return func(o *Options) {
		o.Deep = deep
	}
}

// WithFeed is a configuration point.
func WithFeed(feed float64) Configurator {
	return func(o *Options) {
		o.Feed = feed
	}
}

// WithSecurityZ is a configuration point.
func WithSecurityZ(securityZ float64) Configurator {
	return func(o *Options) {
		o.SecurityZ = securityZ
	}
}

// WithoutStart is a configuration point.
func WithoutStart() Configurator {
	return func(o *Options) {
		o.IgnoreStart = true
	}
}

// WithoutEnd is a configuration point.
func WithoutEnd() Configurator {
	return func(o *Options) {
		o.IgnoreEnd = true
	}
}

// WithOffset is a configuration point.
func WithOffset(offset []float64) Configurator {
	return func(o *Options) {
		o.Offset = offset
	}
}

// Marshal converts any data to gcode.
// data must implements the Marshaler interface.
func Marshal(data any, configs ...Configurator) ([]byte, error) {
	if marshaler, ok := data.(Marshaler); ok {
		return marshaler.MarshallGCode(configs...)
	}

	return nil, fmt.Errorf("%s does not implement gcode.Marchaler", reflect.TypeOf(data).Name())
}

// OffsetX is the tool offset.
func (o Options) OffsetX() float64 {
	if len(o.Offset) < 1 {
		return 0
	}

	return o.Offset[0]
}

// OffsetY is the tool offset.
func (o Options) OffsetY() float64 {
	if len(o.Offset) < 2 {
		return 0
	}

	return o.Offset[1]
}
