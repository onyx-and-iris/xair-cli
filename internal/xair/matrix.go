package xair

import "fmt"

type Matrix struct {
	client      *Client
	baseAddress string
	Eq          *Eq
	Comp        *Comp
}

// newMatrix creates a new Matrix instance
func newMatrix(c *Client) *Matrix {
	return &Matrix{
		client:      c,
		baseAddress: c.addressMap["matrix"],
		Eq:          newEq(c, c.addressMap["matrix"]),
		Comp:        newComp(c, c.addressMap["matrix"]),
	}
}

// Fader requests the current main L/R fader level
func (m *Matrix) Fader(index int) (float64, error) {
	address := fmt.Sprintf(m.baseAddress, index) + "/mix/fader"
	err := m.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := m.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for matrix fader value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetFader sets the matrix fader level
func (m *Matrix) SetFader(index int, level float64) error {
	address := fmt.Sprintf(m.baseAddress, index) + "/mix/fader"
	return m.client.SendMessage(address, float32(mustDbInto(level)))
}

// Mute requests the current matrix mute status
func (m *Matrix) Mute(index int) (bool, error) {
	address := fmt.Sprintf(m.baseAddress, index) + "/mix/on"
	err := m.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := m.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for matrix mute value")
	}
	return val == 0, nil
}

// SetMute sets the matrix mute status
func (m *Matrix) SetMute(index int, muted bool) error {
	address := fmt.Sprintf(m.baseAddress, index) + "/mix/on"
	var value int32
	if !muted {
		value = 1
	}
	return m.client.SendMessage(address, value)
}
