package xair

import "fmt"

type Bus struct {
	client Client
}

func NewBus(c Client) *Bus {
	return &Bus{
		client: c,
	}
}

// Mute requests the current mute status for a bus
func (b *Bus) Mute(bus int) (bool, error) {
	formatter := b.client.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	err := b.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-b.client.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for bus mute value")
	}
	return val == 0, nil
}

// SetMute sets the mute status for a specific bus (1-based indexing)
func (b *Bus) SetMute(bus int, muted bool) error {
	formatter := b.client.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	var value int32
	if !muted {
		value = 1
	}
	return b.client.SendMessage(address, value)
}

// Fader requests the current fader level for a bus
func (b *Bus) Fader(bus int) (float64, error) {
	formatter := b.client.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/fader"
	err := b.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-b.client.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for bus fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetFader sets the fader level for a specific bus (1-based indexing)
func (b *Bus) SetFader(bus int, level float64) error {
	formatter := b.client.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/fader"
	return b.client.SendMessage(address, float32(mustDbInto(level)))
}
