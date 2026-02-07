package xair

import "fmt"

type HeadAmp struct {
	client      *Client
	baseAddress string
}

// newHeadAmp creates a new HeadAmp instance with the provided client.
func newHeadAmp(c *Client) *HeadAmp {
	return &HeadAmp{
		client:      c,
		baseAddress: c.addressMap["headamp"],
	}
}

// Gain gets the gain level for the specified headamp index.
func (h *HeadAmp) Gain(index int) (float64, error) {
	address := fmt.Sprintf(h.baseAddress, index) + "/gain"
	err := h.client.SendMessage(address)
	if err != nil {
		return 0, err
	}

	msg, err := h.client.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	val, ok := msg.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for headamp gain value")
	}

	return linGet(-12, 60, float64(val)), nil
}

// SetGain sets the gain level for the specified headamp index.
func (h *HeadAmp) SetGain(index int, level float64) error {
	address := fmt.Sprintf(h.baseAddress, index) + "/gain"
	return h.client.SendMessage(address, float32(linSet(-12, 60, level)))
}

// PhantomPower gets the phantom power status for the specified headamp index.
func (h *HeadAmp) PhantomPower(index int) (bool, error) {
	address := fmt.Sprintf(h.baseAddress, index) + "/phantom"
	err := h.client.SendMessage(address)
	if err != nil {
		return false, err
	}

	msg, err := h.client.ReceiveMessage()
	if err != nil {
		return false, err
	}
	val, ok := msg.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for phantom power value")
	}

	return val != 0, nil
}

// SetPhantomPower sets the phantom power status for the specified headamp index.
func (h *HeadAmp) SetPhantomPower(index int, enabled bool) error {
	address := fmt.Sprintf(h.baseAddress, index) + "/phantom"
	var val int32
	if enabled {
		val = 1
	} else {
		val = 0
	}
	return h.client.SendMessage(address, val)
}
