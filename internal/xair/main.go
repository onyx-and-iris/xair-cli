package xair

import "fmt"

type Main struct {
	client      *Client
	baseAddress string
	Eq          *Eq
	Comp        *Comp
}

// newMainStereo creates a new Main instance for stereo main output
func newMainStereo(c *Client) *Main {
	return &Main{
		client:      c,
		baseAddress: c.addressMap["main"],
		Eq:          newEqForMain(c, c.addressMap["main"]),
		Comp:        newCompForMain(c, c.addressMap["main"]),
	}
}

// newMainMono creates a new MainMono instance for mono main output (X32 only)
func newMainMono(c *Client) *Main {
	return &Main{
		baseAddress: c.addressMap["mainmono"],
		client:      c,
		Eq:          newEqForMain(c, c.addressMap["mainmono"]),
		Comp:        newCompForMain(c, c.addressMap["mainmono"]),
	}
}

// Fader requests the current main L/R fader level
func (m *Main) Fader() (float64, error) {
	address := m.baseAddress + "/mix/fader"
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
		return 0, fmt.Errorf("unexpected argument type for main LR fader value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetFader sets the main L/R fader level
func (m *Main) SetFader(level float64) error {
	address := m.baseAddress + "/mix/fader"
	return m.client.SendMessage(address, float32(mustDbInto(level)))
}

// Mute requests the current main L/R mute status
func (m *Main) Mute() (bool, error) {
	address := m.baseAddress + "/mix/on"
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
		return false, fmt.Errorf("unexpected argument type for main LR mute value")
	}
	return val == 0, nil
}

// SetMute sets the main L/R mute status
func (m *Main) SetMute(muted bool) error {
	address := m.baseAddress + "/mix/on"
	var value int32
	if !muted {
		value = 1
	}
	return m.client.SendMessage(address, value)
}
