/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package xair

import (
	"fmt"
	"net"
	"time"

	"github.com/charmbracelet/log"

	"github.com/hypebeast/go-osc/osc"
)

type parser interface {
	Parse(data []byte) (*osc.Message, error)
}

type engine struct {
	Kind      string
	conn      *net.UDPConn
	mixerAddr *net.UDPAddr

	parser     parser
	addressMap map[string]string

	done     chan bool
	respChan chan *osc.Message
}

type Client struct {
	engine
}

// NewClient creates a new XAirClient instance
func NewClient(mixerIP string, mixerPort int, opts ...Option) (*Client, error) {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", 0))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local address: %v", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %v", err)
	}

	mixerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", mixerIP, mixerPort))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to resolve mixer address: %v", err)
	}

	log.Debugf("Local UDP connection: %s	", conn.LocalAddr().String())

	e := &engine{
		Kind:      "xair",
		conn:      conn,
		mixerAddr: mixerAddr,
		parser:    newParser(),
		done:      make(chan bool),
		respChan:  make(chan *osc.Message, 100),
	}

	for _, opt := range opts {
		opt(e)
	}

	return &Client{
		engine: *e,
	}, nil
}

// Start begins listening for messages in a goroutine
func (x *Client) StartListening() {
	go x.receiveLoop()
	log.Debugf("Started listening on %s...", x.conn.LocalAddr().String())
}

// receiveLoop handles incoming OSC messages
func (x *Client) receiveLoop() {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-x.done:
			return
		default:
			// Set read timeout to avoid blocking forever
			x.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, _, err := x.conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout is expected, continue loop
					continue
				}
				// Check if we're shutting down to avoid logging expected errors
				select {
				case <-x.done:
					return
				default:
					log.Errorf("Read error: %v", err)
					return
				}
			}

			msg, err := x.parseOSCMessage(buffer[:n])
			if err != nil {
				log.Errorf("Failed to parse OSC message: %v", err)
				continue
			}
			x.respChan <- msg
		}
	}
}

// parseOSCMessage parses raw bytes into an OSC message with improved error handling
func (x *Client) parseOSCMessage(data []byte) (*osc.Message, error) {
	msg, err := x.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Stop stops the client and closes the connection
func (x *Client) Stop() {
	close(x.done)
	if x.conn != nil {
		x.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (x *Client) SendMessage(address string, args ...any) error {
	return x.SendToAddress(x.mixerAddr, address, args...)
}

// SendToAddress sends an OSC message to a specific address (enables replying to different ports)
func (x *Client) SendToAddress(addr *net.UDPAddr, oscAddress string, args ...any) error {
	msg := osc.NewMessage(oscAddress)
	for _, arg := range args {
		msg.Append(arg)
	}

	log.Debugf("Sending to %v: %s", addr, msg.String())
	if len(args) > 0 {
		log.Debug(" - Arguments: ")
		for i, arg := range args {
			if i > 0 {
				log.Debug(", ")
			}
			log.Debugf("%v", arg)
		}
	}
	log.Debug("")

	data, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	_, err = x.conn.WriteToUDP(data, addr)
	return err
}

// RequestInfo requests mixer information
func (x *Client) RequestInfo() (error, InfoResponse) {
	err := x.SendMessage("/xinfo")
	if err != nil {
		return err, InfoResponse{}
	}

	val := <-x.respChan
	var info InfoResponse
	if len(val.Arguments) >= 3 {
		info.Host = val.Arguments[0].(string)
		info.Name = val.Arguments[1].(string)
		info.Model = val.Arguments[2].(string)
	}
	return nil, info
}

// KeepAlive sends keep-alive message (required for multi-client usage)
func (x *Client) KeepAlive() error {
	return x.SendMessage("/xremote")
}

// RequestStatus requests mixer status
func (x *Client) RequestStatus() error {
	return x.SendMessage("/status")
}

/* STRIP METHODS */

// StripMute gets mute state for a specific strip (1-based indexing)
func (x *Client) StripMute(strip int) (bool, error) {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	err := x.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for strip mute value")
	}
	return val == 0, nil
}

// SetStripMute sets mute state for a specific strip (1-based indexing)
func (x *Client) SetStripMute(strip int, muted bool) error {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	var value int32 = 0
	if !muted {
		value = 1
	}
	return x.SendMessage(address, value)
}

// StripFader requests the current fader level for a strip
func (x *Client) StripFader(strip int) (float64, error) {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	err := x.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetStripFader sets the fader level for a specific strip (1-based indexing)
func (x *Client) SetStripFader(strip int, level float64) error {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	return x.SendMessage(address, float32(mustDbInto(level)))
}

// StripMicGain requests the phantom gain for a specific strip
func (x *Client) StripMicGain(strip int) (float64, error) {
	address := fmt.Sprintf("/ch/%02d/mix/gain", strip)
	err := x.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip gain request: %v", err)
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip gain value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetStripMicGain sets the phantom gain for a specific strip (1-based indexing)
func (x *Client) SetStripMicGain(strip int, gain float32) error {
	address := fmt.Sprintf("/ch/%02d/mix/gain", strip)
	return x.SendMessage(address, gain)
}

// StripName requests the name for a specific strip
func (x *Client) StripName(strip int) (string, error) {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	err := x.SendMessage(address)
	if err != nil {
		return "", fmt.Errorf("failed to send strip name request: %v", err)
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected argument type for strip name value")
	}
	return val, nil
}

// SetStripName sets the name for a specific strip
func (x *Client) SetStripName(strip int, name string) error {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	return x.SendMessage(address, name)
}

// StripColor requests the color for a specific strip
func (x *Client) StripColor(strip int) (int32, error) {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	err := x.SendMessage(address)
	if err != nil {
		return 0, fmt.Errorf("failed to send strip color request: %v", err)
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for strip color value")
	}
	return val, nil
}

// SetStripColor sets the color for a specific strip (0-15)
func (x *Client) SetStripColor(strip int, color int32) error {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	return x.SendMessage(address, color)
}

/* BUS METHODS */

// BusMute requests the current mute status for a bus
func (x *Client) BusMute(bus int) (bool, error) {
	formatter := x.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	err := x.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for bus mute value")
	}
	return val == 0, nil
}

// SetBusMute sets the mute status for a specific bus (1-based indexing)
func (x *Client) SetBusMute(bus int, muted bool) error {
	formatter := x.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	var value int32 = 0
	if !muted {
		value = 1
	}
	return x.SendMessage(address, value)
}

// BusFader requests the current fader level for a bus
func (x *Client) BusFader(bus int) (float64, error) {
	address := fmt.Sprintf("/bus/%01d/mix/fader", bus)
	err := x.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for bus fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetBusFader sets the fader level for a specific bus (1-based indexing)
func (x *Client) SetBusFader(bus int, level float64) error {
	address := fmt.Sprintf("/bus/%01d/mix/fader", bus)
	return x.SendMessage(address, float32(mustDbInto(level)))
}

/* MAIN LR METHODS */

// MainLRFader requests the current main L/R fader level
func (x *Client) MainLRFader() (float64, error) {
	err := x.SendMessage("/lr/mix/fader")
	if err != nil {
		return 0, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for main LR fader value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetMainLRFader sets the main L/R fader level
func (x *Client) SetMainLRFader(level float64) error {
	return x.SendMessage("/lr/mix/fader", float32(mustDbInto(level)))
}

// MainLRMute requests the current main L/R mute status
func (x *Client) MainLRMute() (bool, error) {
	err := x.SendMessage("/lr/mix/on")
	if err != nil {
		return false, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for main LR mute value")
	}
	return val == 0, nil
}

// SetMainLRMute sets the main L/R mute status
func (x *Client) SetMainLRMute(muted bool) error {
	var value int32 = 0
	if !muted {
		value = 1
	}
	return x.SendMessage("/lr/mix/on", value)
}
