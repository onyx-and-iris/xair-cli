package xair

import "time"

// EngineOption defines a functional option for configuring the engine.
type EngineOption func(*engine)

// WithTimeout sets the timeout duration for OSC message responses.
func WithTimeout(timeout time.Duration) EngineOption {
	return func(e *engine) {
		e.timeout = timeout
	}
}

// CompOption defines a functional option for configuring Comp parameters.
type CompOption func(*Comp)

// WithCompAddressFunc allows customization of the OSC address formatting for Comp parameters.
func WithCompAddressFunc(f func(fmtString string, args ...any) string) CompOption {
	return func(c *Comp) {
		c.AddressFunc = f
	}
}

// EqOption defines a functional option for configuring Eq parameters.
type EqOption func(*Eq)

// WithEqAddressFunc allows customization of the OSC address formatting for Eq parameters.
func WithEqAddressFunc(f func(fmtString string, args ...any) string) EqOption {
	return func(e *Eq) {
		e.AddressFunc = f
	}
}

// GateOption defines a functional option for configuring Gate parameters.
type GateOption func(*Gate)

// WithGateAddressFunc allows customization of the OSC address formatting for Gate parameters.
func WithGateAddressFunc(f func(fmtString string, args ...any) string) GateOption {
	return func(g *Gate) {
		g.AddressFunc = f
	}
}
