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

// Name gets the name of the snapshot at the given index.
func (s *Snapshot) Name(index int) (string, error) {
	address := s.baseAddress + fmt.Sprintf("/%02d/name", index)
	err := s.client.SendMessage(address)
	if err != nil {
		return "", err
	}

	msg, err := s.client.ReceiveMessage()
	if err != nil {
		return "", err
	}
	name, ok := msg.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for snapshot name")
	}
	return name, nil
}

// SetName sets the name of the snapshot at the given index.
func (s *Snapshot) SetName(index int, name string) error {
	address := s.baseAddress + fmt.Sprintf("/%02d/name", index)
	return s.client.SendMessage(address, name)
}

// CurrentName sets the name of the current snapshot.
func (s *Snapshot) CurrentName(name string) error {
	address := s.baseAddress + "/name"
	return s.client.SendMessage(address, name)
}

// CurrentLoad loads the snapshot at the given index.
func (s *Snapshot) CurrentLoad(index int) error {
	address := s.baseAddress + "/load"
	return s.client.SendMessage(address, int32(index))
}

// CurrentSave saves the current state to the snapshot at the given index.
func (s *Snapshot) CurrentSave(index int) error {
	address := s.baseAddress + "/save"
	return s.client.SendMessage(address, int32(index))
}

// CurrentDelete deletes the snapshot at the given index.
func (s *Snapshot) CurrentDelete(index int) error {
	address := s.baseAddress + "/delete"
	return s.client.SendMessage(address, int32(index))
}
