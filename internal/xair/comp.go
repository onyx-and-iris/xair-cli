package xair

import "fmt"

type Comp struct {
	client      *Client
	baseAddress string
}

// Factory function to create Comp instance for Strip
func newCompForStrip(c *Client) *Comp {
	return &Comp{
		client:      c,
		baseAddress: c.addressMap["strip"],
	}
}

// Factory function to create Comp instance for Bus
func newCompForBus(c *Client) *Comp {
	return &Comp{
		client:      c,
		baseAddress: c.addressMap["bus"],
	}
}

// On retrieves the on/off status of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) On(index int) (bool, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/on"
	err := c.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-c.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for Compressor on value")
	}
	return val != 0, nil
}

// SetOn sets the on/off status of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetOn(index int, on bool) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/on"
	var value int32
	if on {
		value = 1
	}
	return c.client.SendMessage(address, value)
}
