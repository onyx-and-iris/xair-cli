package xair

import "time"

type Option func(*engine)

func WithTimeout(timeout time.Duration) Option {
	return func(e *engine) {
		e.timeout = timeout
	}
}
