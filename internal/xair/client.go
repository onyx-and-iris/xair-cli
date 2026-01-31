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
	Kind      MixerKind
	conn      *net.UDPConn
	mixerAddr *net.UDPAddr

	parser     parser
	addressMap map[string]string

	done     chan bool
	respChan chan *osc.Message
}

type Client struct {
	engine
	Strip *Strip
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
		Kind:       KindXAir,
		conn:       conn,
		mixerAddr:  mixerAddr,
		parser:     newParser(),
		addressMap: addressMapForMixerKind(KindXAir),
		done:       make(chan bool),
		respChan:   make(chan *osc.Message, 100),
	}

	for _, opt := range opts {
		opt(e)
	}

	c := &Client{
		engine: *e,
	}
	c.Strip = NewStrip(*c)

	return c, nil
}

// Start begins listening for messages in a goroutine
func (c *Client) StartListening() {
	go c.receiveLoop()
	log.Debugf("Started listening on %s...", c.conn.LocalAddr().String())
}

// receiveLoop handles incoming OSC messages
func (c *Client) receiveLoop() {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-c.done:
			return
		default:
			// Set read timeout to avoid blocking forever
			c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			n, _, err := c.conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout is expected, continue loop
					continue
				}
				// Check if we're shutting down to avoid logging expected errors
				select {
				case <-c.done:
					return
				default:
					log.Errorf("Read error: %v", err)
					return
				}
			}

			msg, err := c.parseOSCMessage(buffer[:n])
			if err != nil {
				log.Errorf("Failed to parse OSC message: %v", err)
				continue
			}
			c.respChan <- msg
		}
	}
}

// parseOSCMessage parses raw bytes into an OSC message with improved error handling
func (c *Client) parseOSCMessage(data []byte) (*osc.Message, error) {
	msg, err := c.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Stop stops the client and closes the connection
func (c *Client) Stop() {
	close(c.done)
	if c.conn != nil {
		c.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (c *Client) SendMessage(address string, args ...any) error {
	return c.SendToAddress(c.mixerAddr, address, args...)
}

// SendToAddress sends an OSC message to a specific address (enables replying to different ports)
func (c *Client) SendToAddress(addr *net.UDPAddr, oscAddress string, args ...any) error {
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

	_, err = c.conn.WriteToUDP(data, addr)
	return err
}

// RequestInfo requests mixer information
func (c *Client) RequestInfo() (error, InfoResponse) {
	err := c.SendMessage("/xinfo")
	if err != nil {
		return err, InfoResponse{}
	}

	val := <-c.respChan
	var info InfoResponse
	if len(val.Arguments) >= 3 {
		info.Host = val.Arguments[0].(string)
		info.Name = val.Arguments[1].(string)
		info.Model = val.Arguments[2].(string)
	}
	return nil, info
}

// KeepAlive sends keep-alive message (required for multi-client usage)
func (c *Client) KeepAlive() error {
	return c.SendMessage("/xremote")
}

// RequestStatus requests mixer status
func (c *Client) RequestStatus() error {
	return c.SendMessage("/status")
}

/* BUS METHODS */

// BusMute requests the current mute status for a bus
func (c *Client) BusMute(bus int) (bool, error) {
	formatter := c.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	err := c.SendMessage(address)
	if err != nil {
		return false, err
	}

	resp := <-c.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for bus mute value")
	}
	return val == 0, nil
}

// SetBusMute sets the mute status for a specific bus (1-based indexing)
func (c *Client) SetBusMute(bus int, muted bool) error {
	formatter := c.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/on"
	var value int32
	if !muted {
		value = 1
	}
	return c.SendMessage(address, value)
}

// BusFader requests the current fader level for a bus
func (c *Client) BusFader(bus int) (float64, error) {
	formatter := c.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/fader"
	err := c.SendMessage(address)
	if err != nil {
		return 0, err
	}

	resp := <-c.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for bus fader value")
	}

	return mustDbFrom(float64(val)), nil
}

// SetBusFader sets the fader level for a specific bus (1-based indexing)
func (c *Client) SetBusFader(bus int, level float64) error {
	formatter := c.addressMap["bus"]
	address := fmt.Sprintf(formatter, bus) + "/mix/fader"
	return c.SendMessage(address, float32(mustDbInto(level)))
}

/* MAIN LR METHODS */

// MainLRFader requests the current main L/R fader level
func (c *Client) MainLRFader() (float64, error) {
	err := c.SendMessage("/lr/mix/fader")
	if err != nil {
		return 0, err
	}

	resp := <-c.respChan
	val, ok := resp.Arguments[0].(float32)
	if !ok {
		return 0, fmt.Errorf("unexpected argument type for main LR fader value")
	}
	return mustDbFrom(float64(val)), nil
}

// SetMainLRFader sets the main L/R fader level
func (c *Client) SetMainLRFader(level float64) error {
	return c.SendMessage("/lr/mix/fader", float32(mustDbInto(level)))
}

// MainLRMute requests the current main L/R mute status
func (c *Client) MainLRMute() (bool, error) {
	err := c.SendMessage("/lr/mix/on")
	if err != nil {
		return false, err
	}

	resp := <-c.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for main LR mute value")
	}
	return val == 0, nil
}

// SetMainLRMute sets the main L/R mute status
func (c *Client) SetMainLRMute(muted bool) error {
	var value int32
	if !muted {
		value = 1
	}
	return c.SendMessage("/lr/mix/on", value)
}
