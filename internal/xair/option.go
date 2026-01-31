package xair

import "strings"

type Option func(*engine)

func WithKind(kind string) Option {
	if strings.EqualFold(kind, "x32") {
		return func(c *engine) {
			c.Kind = kind
			c.addressMap = x32AddressMap
		}
	}
	return func(c *engine) {
		c.Kind = "xair"
		c.addressMap = xairAddressMap
	}
}
