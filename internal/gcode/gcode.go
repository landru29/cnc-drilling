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

func Marshal(data any, configs ...Configurator) ([]byte, error) {
	if marshaler, ok := data.(Marshaler); ok {
		return marshaler.MarshallGCode(configs...)
	}

	return nil, fmt.Errorf("%s does not implement gcode.Marchaler", reflect.TypeOf(data).Name())
}
