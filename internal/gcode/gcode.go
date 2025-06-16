package gcode

import (
	"fmt"
	"reflect"
)

type Marshaler interface {
	MarshallGCode(configs ...Configurator) ([]byte, error)
}

type Configurator func(*Options)

type Options struct {
	Deep        float64
	Feed        float64
	SecurityZ   float64
	IgnoreStart bool
	IgnoreEnd   bool
	Offset      []float64
}

func WithDeep(deep float64) Configurator {
	return func(o *Options) {
		o.Deep = deep
	}
}

func WithFeed(feed float64) Configurator {
	return func(o *Options) {
		o.Feed = feed
	}
}

func WithSecurityZ(securityZ float64) Configurator {
	return func(o *Options) {
		o.SecurityZ = securityZ
	}
}

func WithoutStart() Configurator {
	return func(o *Options) {
		o.IgnoreStart = true
	}
}

func WithoutEnd() Configurator {
	return func(o *Options) {
		o.IgnoreEnd = true
	}
}

func WithOffset(offset []float64) Configurator {
	return func(o *Options) {
		o.Offset = offset
	}
}

func Marshal(data any, configs ...Configurator) ([]byte, error) {
	if marshaler, ok := data.(Marshaler); ok {
		return marshaler.MarshallGCode(configs...)
	}

	return nil, fmt.Errorf("%s does not implement gcode.Marchaler", reflect.TypeOf(data).Name())
}

func (o Options) OffsetX() float64 {
	if len(o.Offset) < 1 {
		return 0
	}

	return o.Offset[0]
}

func (o Options) OffsetY() float64 {
	if len(o.Offset) < 2 {
		return 0
	}

	return o.Offset[1]
}
