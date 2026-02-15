package xair

import "fmt"

type DCA struct {
	client      *client
	baseAddress string
}

// newDCA creates a new DCA instance
func newDCA(c *client) *DCA {
	return &DCA{
		client:      c,
		baseAddress: c.addressMap["dca"],
	}
}

// Mute requests the current mute status for a DCA group
func (d *DCA) Mute(group int) (bool, error) {
	address := fmt.Sprintf(d.baseAddress, group) + "/on"
	err := d.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := d.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for DCA mute value")
	}
	return val == 0, nil
}

// SetMute sets the mute status for a specific DCA group (1-based indexing)
func (d *DCA) SetMute(group int, muted bool) error {
	address := fmt.Sprintf(d.baseAddress, group) + "/on"
	var value int32
	if !muted {
		value = 1
	}
	return d.client.SendMessage(address, value)
}

// Name requests the current name for a DCA group
func (d *DCA) Name(group int) (string, error) {
	address := fmt.Sprintf(d.baseAddress, group) + "/config/name"
	err := d.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	msg, err := d.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	name, ok := msg.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for DCA name value")
	}
	return name, nil
}

// SetName sets the name for a specific DCA group (1-based indexing)
func (d *DCA) SetName(group int, name string) error {
	address := fmt.Sprintf(d.baseAddress, group) + "/config/name"
	return d.client.SendMessage(address, name)
}

// Colour requests the current colour for a DCA group
func (d *DCA) Colour(group int) (int32, error) {
	address := fmt.Sprintf(d.baseAddress, group) + "/config/colour"
	err := d.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := d.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	colour, ok := msg.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for DCA colour value")
	}
	return colour, nil
}

// SetColor sets the colour for a specific DCA group (1-based indexing)
func (d *DCA) SetColor(group int, colour int32) error {
	address := fmt.Sprintf(d.baseAddress, group) + "/config/colour"
	return d.client.SendMessage(address, colour)
}
