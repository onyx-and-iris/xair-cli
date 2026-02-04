package xair

import "fmt"

type Snapshot struct {
	baseAddress string
	client      *Client
}

func NewSnapshot(c *Client) *Snapshot {
	return &Snapshot{
		baseAddress: c.addressMap["snapshot"],
		client:      c,
	}
}

// Name sets the name of the snapshot at the given index.
func (s *Snapshot) Name(index int, name string) error {
	address := s.baseAddress + fmt.Sprintf("/name/%d", index)
	return s.client.SendMessage(address, name)
}

// Load loads the snapshot at the given index.
func (s *Snapshot) Load(index int) error {
	address := s.baseAddress + fmt.Sprintf("/load/%d", index)
	return s.client.SendMessage(address)
}

// Save saves the current state to the snapshot at the given index.
func (s *Snapshot) Save(index int) error {
	address := s.baseAddress + fmt.Sprintf("/save/%d", index)
	return s.client.SendMessage(address)
}

// Delete deletes the snapshot at the given index.
func (s *Snapshot) Delete(index int) error {
	address := s.baseAddress + fmt.Sprintf("/delete/%d", index)
	return s.client.SendMessage(address)
}
