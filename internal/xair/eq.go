package xair

import (
	"fmt"
)

// Eq represents the EQ parameters.
type Eq struct {
	client      *Client
	baseAddress string
	AddressFunc func(fmtString string, args ...any) string
}

// Factory function to create Eq instance with optional configuration
func newEq(c *Client, baseAddress string, opts ...EqOption) *Eq {
	eq := &Eq{
		client:      c,
		baseAddress: fmt.Sprintf("%s/eq", baseAddress),
		AddressFunc: fmt.Sprintf,
	}

	for _, opt := range opts {
		opt(eq)
	}

	return eq
}

// On retrieves the on/off status of the EQ for a specific strip or bus (1-based indexing).
func (e *Eq) On(index int) (bool, error) {
	address := e.AddressFunc(e.baseAddress, index) + "/on"
	err := e.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for EQ on value")
	}
	return val != 0, nil
}

// SetOn sets the on/off status of the EQ for a specific strip or bus (1-based indexing).
func (e *Eq) SetOn(index int, on bool) error {
	address := e.AddressFunc(e.baseAddress, index) + "/on"
	var value int32
	if on {
		value = 1
	}
	return e.client.SendMessage(address, value)
}

func (e *Eq) Mode(index int) (string, error) {
	address := e.AddressFunc(e.baseAddress, index) + "/mode"
	err := e.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	possibleModes := []string{"peq", "geq", "teq"}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for EQ mode value")
	}
	return possibleModes[val], nil
}

func (e *Eq) SetMode(index int, mode string) error {
	address := e.AddressFunc(e.baseAddress, index) + "/mode"
	possibleModes := []string{"peq", "geq", "teq"}
	return e.client.SendMessage(address, int32(indexOf(possibleModes, mode)))
}

// Gain retrieves the gain for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Gain(index int, band int) (float64, error) {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/g", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ gain value")
	}
	return linGet(-15, 15, float64(val)), nil
}

// SetGain sets the gain for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetGain(index int, band int, gain float64) error {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/g", band)
	return e.client.SendMessage(address, float32(linSet(-15, 15, gain)))
}

// Frequency retrieves the frequency for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Frequency(index int, band int) (float64, error) {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/f", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ frequency value")
	}
	return logGet(20, 20000, float64(val)), nil
}

// SetFrequency sets the frequency for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetFrequency(index int, band int, frequency float64) error {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/f", band)
	return e.client.SendMessage(address, float32(logSet(20, 20000, frequency)))
}

// Q retrieves the Q factor for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Q(index int, band int) (float64, error) {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/q", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for EQ Q value")
	}
	return logGet(0.3, 10, 1.0-float64(val)), nil
}

// SetQ sets the Q factor for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetQ(index int, band int, q float64) error {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/q", band)
	return e.client.SendMessage(address, float32(1.0-logSet(0.3, 10, q)))
}

// Type retrieves the type for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) Type(index int, band int) (string, error) {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/type", band)
	err := e.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	possibleTypes := []string{"lcut", "lshv", "peq", "veq", "hshv", "hcut"}

	msg, err := e.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for EQ type value")
	}
	return possibleTypes[val], nil
}

// SetType sets the type for a specific EQ band on a strip or bus (1-based indexing).
func (e *Eq) SetType(index int, band int, eqType string) error {
	address := e.AddressFunc(e.baseAddress, index) + fmt.Sprintf("/%d/type", band)
	possibleTypes := []string{"lcut", "lshv", "peq", "veq", "hshv", "hcut"}
	return e.client.SendMessage(address, int32(indexOf(possibleTypes, eqType)))
}
