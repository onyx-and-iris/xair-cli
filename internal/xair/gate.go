package xair

import "fmt"

type Gate struct {
	client      *Client
	baseAddress string
}

// Factory function to create Gate instance for Strip
func newGateForStrip(c *Client, baseAddress string) *Gate {
	return &Gate{
		client:      c,
		baseAddress: fmt.Sprintf("%s/gate", baseAddress),
	}
}

// On retrieves the on/off status of the Gate for a specific strip (1-based indexing).
func (g *Gate) On(index int) (bool, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/on"
	err := g.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for Gate on value")
	}
	return val != 0, nil
}

// SetOn sets the on/off status of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetOn(index int, on bool) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/on"
	var value int32
	if on {
		value = 1
	}
	return g.client.SendMessage(address, value)
}

// Mode retrieves the current mode of the Gate for a specific strip (1-based indexing).
func (g *Gate) Mode(index int) (string, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/mode"
	err := g.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	possibleModes := []string{"exp2", "exp3", "exp4", "gate", "duck"}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for Gate mode value")
	}
	return possibleModes[val], nil
}

// SetMode sets the mode of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetMode(index int, mode string) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/mode"
	possibleModes := []string{"exp2", "exp3", "exp4", "gate", "duck"}

	return g.client.SendMessage(address, int32(indexOf(possibleModes, mode)))
}

// Threshold retrieves the threshold value of the Gate for a specific strip (1-based indexing).
func (g *Gate) Threshold(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/thr"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate threshold value")
	}
	return linGet(-80, 0, float64(val)), nil
}

// SetThreshold sets the threshold value of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetThreshold(index int, threshold float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/thr"
	return g.client.SendMessage(address, float32(linSet(-80, 0, threshold)))
}

// Range retrieves the range value of the Gate for a specific strip (1-based indexing).
func (g *Gate) Range(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/range"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate range value")
	}
	return linGet(3, 60, float64(val)), nil
}

// SetRange sets the range value of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetRange(index int, rangeVal float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/range"
	return g.client.SendMessage(address, float32(linSet(3, 60, rangeVal)))
}

// Attack retrieves the attack time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Attack(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/attack"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate attack value")
	}
	return linGet(0, 120, float64(val)), nil
}

// SetAttack sets the attack time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetAttack(index int, attack float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/attack"
	return g.client.SendMessage(address, float32(linSet(0, 120, attack)))
}

// Hold retrieves the hold time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Hold(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/hold"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate hold value")
	}
	return logGet(0.02, 2000, float64(val)), nil
}

// SetHold sets the hold time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetHold(index int, hold float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/hold"
	return g.client.SendMessage(address, float32(logSet(0.02, 2000, hold)))
}

// Release retrieves the release time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Release(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/release"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := g.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate release value")
	}
	return logGet(5, 4000, float64(val)), nil
}

// SetRelease sets the release time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetRelease(index int, release float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/release"
	return g.client.SendMessage(address, float32(logSet(5, 4000, release)))
}
