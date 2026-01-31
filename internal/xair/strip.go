package xair

import "fmt"

type Strip struct {
	client Client
}

func NewStrip(c Client) *Strip {
	return &Strip{
		client: c,
	}
}

// Mute gets the mute status of the specified strip (1-based indexing).
func (s *Strip) Mute(strip int) (bool, error) {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	err := s.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-s.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for strip mute value")
	}
	return val == 0, nil
}

// SetMute sets the mute status of the specified strip (1-based indexing).
func (s *Strip) SetMute(strip int, muted bool) error {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	var value int32 = 0
	if !muted {
		value = 1
	}
	return s.client.SendMessage(address, value)
}

// Fader gets the fader level of the specified strip (1-based indexing).
func (s *Strip) Fader(strip int) (float64, error) {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-s.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetFader sets the fader level of the specified strip (1-based indexing).
func (s *Strip) SetFader(strip int, level float64) error {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	return s.client.SendMessage(address, float32(mustDbInto(level)))
}

// MicGain requests the phantom gain for a specific strip (1-based indexing).
func (s *Strip) MicGain(strip int) (float64, error) {
	address := fmt.Sprintf("/ch/%02d/mix/gain", strip)
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip gain request: %v", err)
	}

	resp := <-s.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip gain value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetMicGain sets the phantom gain for a specific strip (1-based indexing).
func (s *Strip) SetMicGain(strip int, gain float32) error {
	address := fmt.Sprintf("/ch/%02d/mix/gain", strip)
	return s.client.SendMessage(address, gain)
}

// Name requests the name for a specific strip
func (s *Strip) Name(strip int) (string, error) {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	err := s.client.SendMessage(address)
	if err != nil {
		return "", fmt.Errorf("failed to send strip name request: %v", err)
	}

	resp := <-s.client.respChan
	val, ok := resp.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for strip name value")
	}
	return val, nil
}

// SetName sets the name for a specific strip
func (s *Strip) SetName(strip int, name string) error {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	return s.client.SendMessage(address, name)
}

// Color requests the color for a specific strip
func (s *Strip) Color(strip int) (int32, error) {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip color request: %v", err)
	}

	resp := <-s.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip color value")
	}
	return val, nil
}

// SetColor sets the color for a specific strip (0-15)
func (s *Strip) SetColor(strip int, color int32) error {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	return s.client.SendMessage(address, color)
}
