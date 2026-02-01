package xair

type Gate struct {
	client *Client
}

func newGate(c *Client) *Gate {
	return &Gate{client: c}
}
