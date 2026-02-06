package xair

import "time"

type Option func(*engine)

func WithKind(kind string) Option {
	return func(e *engine) {
		e.Kind = MixerKind(kind)
		e.addressMap = addressMapForMixerKind(e.Kind)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(e *engine) {
		e.timeout = timeout
	}
}
