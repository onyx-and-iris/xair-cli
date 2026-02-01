package xair

import "fmt"

type Eq struct {
	client      *Client
	baseAddress string
}

// Factory function to create Eq instance for Strip
func newEqForStrip(c *Client) *Eq {
	return &Eq{
		client:      c,
		baseAddress: c.addressMap["strip"],
	}
}

// Factory function to create Eq instance for Bus
func newEqForBus(c *Client) *Eq {
	return &Eq{
		client:      c,
		baseAddress: c.addressMap["bus"],
	}
}

// On retrieves the on/off status of the EQ for a specific strip or bus (1-based indexing).
func (e *Eq) On(index int) (bool, error) {
	address := fmt.Sprintf(e.baseAddress, index) + "/eq/on"
	err := e.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for EQ on value")
	}
	return val != 0, nil
}

// SetOn sets the on/off status of the EQ for a specific strip or bus (1-based indexing).
func (e *Eq) SetOn(index int, on bool) error {
	address := fmt.Sprintf(e.baseAddress, index) + "/eq/on"
	var value int32
	if on {
		value = 1
	}
	return e.client.SendMessage(address, value)
}

// Gain retrieves the gain for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Gain(index int, band int) (float64, error) {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/g", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ gain value")
	}
	return float64(val), nil
}

// SetGain sets the gain for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetGain(index int, band int, gain float64) error {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/g", band)
	return e.client.SendMessage(address, float32(gain))
}
