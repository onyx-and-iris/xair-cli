package xair

import "fmt"

type Main struct {
	client Client
}

func newMain(c Client) *Main {
	return &Main{
		client: c,
	}
}

// Fader requests the current main L/R fader level
func (m *Main) Fader() (float64, error) {
	err := m.client.SendMessage("/lr/mix/fader")
	if err != nil {
		return 0, err
	}

	resp := <-m.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for main LR fader value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetFader sets the main L/R fader level
func (m *Main) SetFader(level float64) error {
	return m.client.SendMessage("/lr/mix/fader", float32(mustDbInto(level)))
}

// Mute requests the current main L/R mute status
func (m *Main) Mute() (bool, error) {
	err := m.client.SendMessage("/lr/mix/on")
	if err != nil {
		return false, err
	}

	resp := <-m.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for main LR mute value")
	}
	return val == 0, nil
}

// SetMute sets the main L/R mute status
func (m *Main) SetMute(muted bool) error {
	var value int32
	if !muted {
		value = 1
	}
	return m.client.SendMessage("/lr/mix/on", value)
}
