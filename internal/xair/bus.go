package xair

import "fmt"

type Bus struct {
	client      *Client
	baseAddress string
	Eq          *Eq
	Comp        *Comp
}

// newBus creates a new Bus instance
func newBus(c *Client) *Bus {
	return &Bus{
		client:      c,
		baseAddress: c.addressMap["bus"],
		Eq:          newEq(c, c.addressMap["bus"]),
		Comp:        newComp(c, c.addressMap["bus"]),
	}
}

// Mute requests the current mute status for a bus
func (b *Bus) Mute(bus int) (bool, error) {
	address := fmt.Sprintf(b.baseAddress, bus) + "/mix/on"
	err := b.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := b.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for bus mute value")
	}
	return val == 0, nil
}

// SetMute sets the mute status for a specific bus (1-based indexing)
func (b *Bus) SetMute(bus int, muted bool) error {
	address := fmt.Sprintf(b.baseAddress, bus) + "/mix/on"
	var value int32
	if !muted {
		value = 1
	}
	return b.client.SendMessage(address, value)
}

// Fader requests the current fader level for a bus
func (b *Bus) Fader(bus int) (float64, error) {
	address := fmt.Sprintf(b.baseAddress, bus) + "/mix/fader"
	err := b.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := b.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for bus fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetFader sets the fader level for a specific bus (1-based indexing)
func (b *Bus) SetFader(bus int, level float64) error {
	address := fmt.Sprintf(b.baseAddress, bus) + "/mix/fader"
	return b.client.SendMessage(address, float32(mustDbInto(level)))
}

// Name requests the name for a specific bus
func (b *Bus) Name(bus int) (string, error) {
	address := fmt.Sprintf(b.baseAddress, bus) + "/config/name"
	err := b.client.SendMessage(address)
	if err != nil {
		return "", fmt.Errorf("failed to send bus name request: %v", err)
	}

	msg, err := b.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	val, ok := msg.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for bus name value")
	}
	return val, nil
}

// SetName sets the name for a specific bus
func (b *Bus) SetName(bus int, name string) error {
	address := fmt.Sprintf(b.baseAddress, bus) + "/config/name"
	return b.client.SendMessage(address, name)
}
