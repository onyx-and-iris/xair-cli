package xair

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/hypebeast/go-osc/osc"
)

// XAirClient is a client for controlling XAir mixers
type XAirClient struct {
	client
	Main     *Main
	Strip    *Strip
	Bus      *Bus
	HeadAmp  *HeadAmp
	Snapshot *Snapshot
	DCA      *DCA
}

// NewXAirClient creates a new XAirClient instance with optional engine configuration
func NewXAirClient(mixerIP string, mixerPort int, opts ...EngineOption) (*XAirClient, error) {
	e, err := newEngine(mixerIP, mixerPort, kindXAir, opts...)
	if err != nil {
		return nil, err
	}

	c := &XAirClient{
		client: client{e, InfoResponse{}},
	}
	c.Main = newMainStereo(&c.client)
	c.Strip = newStrip(&c.client)
	c.Bus = newBus(&c.client)
	c.HeadAmp = newHeadAmp(&c.client)
	c.Snapshot = newSnapshot(&c.client)
	c.DCA = newDCA(&c.client)

	return c, nil
}

// X32Client is a client for controlling X32 mixers
type X32Client struct {
	client
	Main     *Main
	MainMono *Main
	Matrix   *Matrix
	Strip    *Strip
	Bus      *Bus
	HeadAmp  *HeadAmp
	Snapshot *Snapshot
	DCA      *DCA
}

// NewX32Client creates a new X32Client instance with optional engine configuration
func NewX32Client(mixerIP string, mixerPort int, opts ...EngineOption) (*X32Client, error) {
	e, err := newEngine(mixerIP, mixerPort, kindX32, opts...)
	if err != nil {
		return nil, err
	}

	c := &X32Client{
		client: client{e, InfoResponse{}},
	}
	c.Main = newMainStereo(&c.client)
	c.MainMono = newMainMono(&c.client)
	c.Matrix = newMatrix(&c.client)
	c.Strip = newStrip(&c.client)
	c.Bus = newBus(&c.client)
	c.HeadAmp = newHeadAmp(&c.client)
	c.Snapshot = newSnapshot(&c.client)
	c.DCA = newDCA(&c.client)

	return c, nil
}

type client struct {
	*engine
	Info InfoResponse
}

// Start begins listening for messages in a goroutine
func (c *client) StartListening() {
	go c.engine.receiveLoop()
	log.Debugf("Started listening on %s...", c.engine.conn.LocalAddr().String())
}

// Close stops the client and closes the connection
func (c *client) Close() {
	close(c.engine.done)
	if c.engine.conn != nil {
		c.engine.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (c *client) SendMessage(address string, args ...any) error {
	return c.engine.sendToAddress(c.mixerAddr, address, args...)
}

// ReceiveMessage receives an OSC message from the mixer
func (c *client) ReceiveMessage() (*osc.Message, error) {
	t := time.Tick(c.engine.timeout)
	select {
	case <-t:
		return nil, fmt.Errorf("timeout waiting for response")
	case msg := <-c.respChan:
		if msg == nil {
			return nil, fmt.Errorf("no message received")
		}
		return msg, nil
	}
}

// RequestInfo requests mixer information
func (c *client) RequestInfo() (InfoResponse, error) {
	var info InfoResponse
	err := c.SendMessage("/xinfo")
	if err != nil {
		return info, err
	}

	msg, err := c.ReceiveMessage()
	if err != nil {
		return info, err
	}
	if len(msg.Arguments) == 4 {
		if host, ok := msg.Arguments[0].(string); ok {
			info.Host = host
		}
		if name, ok := msg.Arguments[1].(string); ok {
			info.Name = name
		}
		if model, ok := msg.Arguments[2].(string); ok {
			info.Model = model
		}
		if firmware, ok := msg.Arguments[3].(string); ok {
			info.Firmware = firmware
		}
	}
	c.Info = info

	return info, nil
}

// KeepAlive sends keep-alive message (required for multi-client usage)
func (c *client) KeepAlive() error {
	return c.SendMessage("/xremote")
}

// RequestStatus requests mixer status
func (c *client) RequestStatus() error {
	return c.SendMessage("/status")
}
