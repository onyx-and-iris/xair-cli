package xair

// InfoResponse represents the response from the /info OSC address, containing details about the X-Air device.
type InfoResponse struct {
	Host     string
	Name     string
	Model    string
	Firmware string
}
