package xair

type Comp struct {
	client      *Client
	baseAddress string
}

func newCompForStrip(c *Client) *Comp {
	return &Comp{
		client:      c,
		baseAddress: c.addressMap["strip"],
	}
}

func newCompForBus(c *Client) *Comp {
	return &Comp{
		client:      c,
		baseAddress: c.addressMap["bus"],
	}
}
