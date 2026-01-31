package xair

import (
	"fmt"
	"net"

	"github.com/charmbracelet/log"

	"github.com/hypebeast/go-osc/osc"
)

type parser interface {
	Parse(data []byte) (*osc.Message, error)
}

type Client struct {
	engine
	Strip *Strip
	Bus   *Bus
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
	c.Bus = NewBus(*c)

	return c, nil
}

// Start begins listening for messages in a goroutine
func (c *Client) StartListening() {
	go c.engine.receiveLoop()
	log.Debugf("Started listening on %s...", c.engine.conn.LocalAddr().String())
}

// Stop stops the client and closes the connection
func (c *Client) Stop() {
	close(c.engine.done)
	if c.engine.conn != nil {
		c.engine.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (c *Client) SendMessage(address string, args ...any) error {
	return c.engine.sendToAddress(c.mixerAddr, address, args...)
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
