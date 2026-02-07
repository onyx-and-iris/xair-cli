package xair

import "time"

type EngineOption func(*engine)

// WithTimeout sets the timeout duration for OSC message responses
func WithTimeout(timeout time.Duration) EngineOption {
	return func(e *engine) {
		e.timeout = timeout
	}
}

type CompOption func(*Comp)

// WithCompAddressFunc allows customization of the OSC address formatting for Comp parameters
func WithCompAddressFunc(f func(fmtString string, args ...any) string) CompOption {
	return func(c *Comp) {
		c.AddressFunc = f
	}
}

type EqOption func(*Eq)

// WithEqAddressFunc allows customization of the OSC address formatting for Eq parameters
func WithEqAddressFunc(f func(fmtString string, args ...any) string) EqOption {
	return func(e *Eq) {
		e.AddressFunc = f
	}
}

type GateOption func(*Gate)

// WithGateAddressFunc allows customization of the OSC address formatting for Gate parameters
func WithGateAddressFunc(f func(fmtString string, args ...any) string) GateOption {
	return func(g *Gate) {
		g.AddressFunc = f
	}
}
