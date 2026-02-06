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

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
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

// Mode retrieves the current mode of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Mode(index int) (string, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mode"
	err := c.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	possibleModes := []string{"comp", "exp"}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for Compressor mode value")
	}
	return possibleModes[val], nil
}

// SetMode sets the mode of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetMode(index int, mode string) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mode"
	possibleModes := []string{"comp", "exp"}
	return c.client.SendMessage(address, int32(indexOf(possibleModes, mode)))
}

// Threshold retrieves the threshold value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Threshold(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/thr"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor threshold value")
	}
	return linGet(-60, 0, float64(val)), nil
}

// SetThreshold sets the threshold value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetThreshold(index int, threshold float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/thr"
	return c.client.SendMessage(address, float32(linSet(-60, 0, threshold)))
}

// Ratio retrieves the ratio value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Ratio(index int) (float32, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/ratio"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	possibleValues := []float32{1.1, 1.3, 1.5, 2.0, 2.5, 3.0, 4.0, 5.0, 7.0, 10, 20, 100}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor ratio value")
	}

	return possibleValues[val], nil
}

// SetRatio sets the ratio value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetRatio(index int, ratio float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/ratio"
	possibleValues := []float32{1.1, 1.3, 1.5, 2.0, 2.5, 3.0, 4.0, 5.0, 7.0, 10, 20, 100}

	return c.client.SendMessage(address, int32(indexOf(possibleValues, float32(ratio))))
}

// Attack retrieves the attack time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Attack(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/attack"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor attack value")
	}
	return linGet(0, 120, float64(val)), nil
}

// SetAttack sets the attack time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetAttack(index int, attack float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/attack"
	return c.client.SendMessage(address, float32(linSet(0, 120, attack)))
}

// Hold retrieves the hold time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Hold(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/hold"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor hold value")
	}
	return logGet(0.02, 2000, float64(val)), nil
}

// SetHold sets the hold time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetHold(index int, hold float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/hold"
	return c.client.SendMessage(address, float32(logSet(0.02, 2000, hold)))
}

// Release retrieves the release time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Release(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/release"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor release value")
	}
	return logGet(4, 4000, float64(val)), nil
}

// SetRelease sets the release time of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetRelease(index int, release float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/release"
	return c.client.SendMessage(address, float32(logSet(4, 4000, release)))
}

// Makeup retrieves the makeup gain of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Makeup(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mgain"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor makeup gain value")
	}
	return linGet(0, 24, float64(val)), nil
}

// SetMakeup sets the makeup gain of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetMakeup(index int, makeup float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mgain"
	return c.client.SendMessage(address, float32(linSet(0, 24, makeup)))
}

// Mix retrieves the mix value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) Mix(index int) (float64, error) {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mix"
	err := c.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := c.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Compressor mix value")
	}
	return linGet(0, 100, float64(val)), nil
}

// SetMix sets the mix value of the Compressor for a specific strip or bus (1-based indexing).
func (c *Comp) SetMix(index int, mix float64) error {
	address := fmt.Sprintf(c.baseAddress, index) + "/dyn/mix"
	return c.client.SendMessage(address, float32(linSet(0, 100, mix)))
}
