package xair

import "fmt"

type Strip struct {
	client      *Client
	baseAddress string
	Gate        *Gate
	Eq          *Eq
	Comp        *Comp
}

func newStrip(c *Client) *Strip {
	return &Strip{
		client:      c,
		baseAddress: c.addressMap["strip"],
		Gate:        newGateForStrip(c, c.addressMap["strip"]),
		Eq:          newEqForStrip(c, c.addressMap["strip"]),
		Comp:        newCompForStrip(c, c.addressMap["strip"]),
	}
}

// Mute gets the mute status of the specified strip (1-based indexing).
func (s *Strip) Mute(index int) (bool, error) {
	address := fmt.Sprintf(s.baseAddress, index) + "/mix/on"
	err := s.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for strip mute value")
	}
	return val == 0, nil
}

// SetMute sets the mute status of the specified strip (1-based indexing).
func (s *Strip) SetMute(strip int, muted bool) error {
	address := fmt.Sprintf(s.baseAddress, strip) + "/mix/on"
	var value int32 = 0
	if !muted {
		value = 1
	}
	return s.client.SendMessage(address, value)
}

// Fader gets the fader level of the specified strip (1-based indexing).
func (s *Strip) Fader(strip int) (float64, error) {
	address := fmt.Sprintf(s.baseAddress, strip) + "/mix/fader"
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetFader sets the fader level of the specified strip (1-based indexing).
func (s *Strip) SetFader(strip int, level float64) error {
	address := fmt.Sprintf(s.baseAddress, strip) + "/mix/fader"
	return s.client.SendMessage(address, float32(mustDbInto(level)))
}

// Name requests the name for a specific strip
func (s *Strip) Name(strip int) (string, error) {
	address := fmt.Sprintf(s.baseAddress, strip) + "/config/name"
	err := s.client.SendMessage(address)
	if err != nil {
		return "", fmt.Errorf("failed to send strip name request: %v", err)
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for strip name value")
	}
	return val, nil
}

// SetName sets the name for a specific strip
func (s *Strip) SetName(strip int, name string) error {
	address := fmt.Sprintf(s.baseAddress, strip) + "/config/name"
	return s.client.SendMessage(address, name)
}

// Color requests the color for a specific strip
func (s *Strip) Color(strip int) (int32, error) {
	address := fmt.Sprintf(s.baseAddress, strip) + "/config/color"
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip color request: %v", err)
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip color value")
	}
	return val, nil
}

// SetColor sets the color for a specific strip (0-15)
func (s *Strip) SetColor(strip int, color int32) error {
	address := fmt.Sprintf(s.baseAddress, strip) + "/config/color"
	return s.client.SendMessage(address, color)
}

// Sends requests the sends level for a mixbus.
func (s *Strip) SendLevel(strip int, bus int) (float64, error) {
	address := fmt.Sprintf(s.baseAddress, strip) + fmt.Sprintf("/mix/%02d/level", bus)
	err := s.client.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip send level request: %v", err)
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip send level value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetSendLevel sets the sends level for a mixbus.
func (s *Strip) SetSendLevel(strip int, bus int, level float64) error {
	address := fmt.Sprintf(s.baseAddress, strip) + fmt.Sprintf("/mix/%02d/level", bus)
	return s.client.SendMessage(address, float32(mustDbInto(level)))
}
