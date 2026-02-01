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
