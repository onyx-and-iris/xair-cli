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

func (e *Eq) Mode(index int) (int, error) {
	address := fmt.Sprintf(e.baseAddress, index) + "/eq/mode"
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ mode value")
	}
	return int(val), nil
}

func (e *Eq) SetMode(index int, mode int) error {
	address := fmt.Sprintf(e.baseAddress, index) + "/eq/mode"
	return e.client.SendMessage(address, int32(mode))
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
	return linGet(-15, 15, float64(val)), nil
}

// SetGain sets the gain for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetGain(index int, band int, gain float64) error {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/g", band)
	return e.client.SendMessage(address, float32(linSet(-15, 15, gain)))
}

// Frequency retrieves the frequency for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Frequency(index int, band int) (float64, error) {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/f", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ frequency value")
	}
	return logGet(20, 20000, float64(val)), nil
}

// SetFrequency sets the frequency for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetFrequency(index int, band int, frequency float64) error {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/f", band)
	return e.client.SendMessage(address, float32(logSet(20, 20000, frequency)))
}

// Q retrieves the Q factor for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Q(index int, band int) (float64, error) {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/q", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ Q value")
	}
	return logGet(0.3, 10, 1.0-float64(val)), nil
}

// SetQ sets the Q factor for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetQ(index int, band int, q float64) error {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/q", band)
	return e.client.SendMessage(address, float32(1.0-logSet(0.3, 10, q)))
}

// Type retrieves the type for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Type(index int, band int) (int, error) {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/type", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-e.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ type value")
	}
	return int(val), nil
}

// SetType sets the type for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetType(index int, band int, eqType int) error {
	address := fmt.Sprintf(e.baseAddress, index) + fmt.Sprintf("/eq/%d/type", band)
	return e.client.SendMessage(address, int32(eqType))
}
