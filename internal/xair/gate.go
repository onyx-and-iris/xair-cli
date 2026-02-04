package xair

import "fmt"

type Gate struct {
	client      *Client
	baseAddress string
}

func newGate(c *Client) *Gate {
	return &Gate{client: c, baseAddress: c.addressMap["strip"]}
}

// On retrieves the on/off status of the Gate for a specific strip (1-based indexing).
func (g *Gate) On(index int) (bool, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/on"
	err := g.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for Gate on value")
	}
	return val != 0, nil
}

// SetOn sets the on/off status of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetOn(index int, on bool) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/on"
	var value int32
	if on {
		value = 1
	}
	return g.client.SendMessage(address, value)
}

// Mode retrieves the current mode of the Gate for a specific strip (1-based indexing).
func (g *Gate) Mode(index int) (string, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/mode"
	err := g.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	possibleModes := []string{"exp2", "exp3", "exp4", "gate", "duck"}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for Gate mode value")
	}
	return possibleModes[val], nil
}

// SetMode sets the mode of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetMode(index int, mode string) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/mode"
	possibleModes := []string{"exp2", "exp3", "exp4", "gate", "duck"}

	return g.client.SendMessage(address, int32(indexOf(possibleModes, mode)))
}

// Threshold retrieves the threshold value of the Gate for a specific strip (1-based indexing).
func (g *Gate) Threshold(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/thr"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate threshold value")
	}
	return linGet(-80, 0, float64(val)), nil
}

// SetThreshold sets the threshold value of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetThreshold(index int, threshold float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/thr"
	return g.client.SendMessage(address, float32(linSet(-80, 0, threshold)))
}

// Range retrieves the range value of the Gate for a specific strip (1-based indexing).
func (g *Gate) Range(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/range"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate range value")
	}
	return linGet(3, 60, float64(val)), nil
}

// SetRange sets the range value of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetRange(index int, rangeVal float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/range"
	return g.client.SendMessage(address, float32(linSet(3, 60, rangeVal)))
}

// Attack retrieves the attack time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Attack(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/attack"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate attack value")
	}
	return linGet(0, 120, float64(val)), nil
}

// SetAttack sets the attack time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetAttack(index int, attack float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/attack"
	return g.client.SendMessage(address, float32(linSet(0, 120, attack)))
}

// Hold retrieves the hold time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Hold(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/hold"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate hold value")
	}
	return logGet(0.02, 2000, float64(val)), nil
}

// SetHold sets the hold time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetHold(index int, hold float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/hold"
	return g.client.SendMessage(address, float32(logSet(0.02, 2000, hold)))
}

// Release retrieves the release time of the Gate for a specific strip (1-based indexing).
func (g *Gate) Release(index int) (float64, error) {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/release"
	err := g.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-g.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for Gate release value")
	}
	return logGet(5, 4000, float64(val)), nil
}

// SetRelease sets the release time of the Gate for a specific strip (1-based indexing).
func (g *Gate) SetRelease(index int, release float64) error {
	address := fmt.Sprintf(g.baseAddress, index) + "/gate/release"
	return g.client.SendMessage(address, float32(logSet(5, 4000, release)))
}
