package xair

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"

	"github.com/hypebeast/go-osc/osc"
)

type Client struct {
	*engine
}

type XAirClient struct {
	Client
	Main     *Main
	Strip    *Strip
	Bus      *Bus
	HeadAmp  *HeadAmp
	Snapshot *Snapshot
}

type X32Client struct {
	Client
	Main     *Main
	MainMono *Main
	Matrix   *Matrix
	Strip    *Strip
	Bus      *Bus
	HeadAmp  *HeadAmp
	Snapshot *Snapshot
}

// NewX32Client creates a new X32Client instance
func NewX32Client(mixerIP string, mixerPort int, opts ...Option) (*X32Client, error) {
	e, err := newEngine(mixerIP, mixerPort, kindX32, opts...)
	if err != nil {
		return nil, err
	}

	c := &X32Client{
		Client: Client{e},
	}
	c.Main = newMainStereo(&c.Client)
	c.MainMono = newMainMono(&c.Client)
	c.Matrix = newMatrix(&c.Client)
	c.Strip = newStrip(&c.Client)
	c.Bus = newBus(&c.Client)
	c.HeadAmp = newHeadAmp(&c.Client)
	c.Snapshot = newSnapshot(&c.Client)

	return c, nil
}

// NewXAirClient creates a new XAirClient instance
func NewXAirClient(mixerIP string, mixerPort int, opts ...Option) (*XAirClient, error) {
	e, err := newEngine(mixerIP, mixerPort, kindXAir, opts...)
	if err != nil {
		return nil, err
	}

	c := &XAirClient{
		Client: Client{e},
	}
	c.Main = newMainStereo(&c.Client)
	c.Strip = newStrip(&c.Client)
	c.Bus = newBus(&c.Client)
	c.HeadAmp = newHeadAmp(&c.Client)
	c.Snapshot = newSnapshot(&c.Client)

	return c, nil
}

// Start begins listening for messages in a goroutine
func (c *Client) StartListening() {
	go c.engine.receiveLoop()
	log.Debugf("Started listening on %s...", c.engine.conn.LocalAddr().String())
}

// Close stops the client and closes the connection
func (c *Client) Close() {
	close(c.engine.done)
	if c.engine.conn != nil {
		c.engine.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (c *Client) SendMessage(address string, args ...any) error {
	return c.engine.sendToAddress(c.mixerAddr, address, args...)
}

// ReceiveMessage receives an OSC message from the mixer
func (c *Client) ReceiveMessage() (*osc.Message, error) {
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
func (c *Client) RequestInfo() (InfoResponse, error) {
	var info InfoResponse
	err := c.SendMessage("/xinfo")
	if err != nil {
		return info, err
	}

	msg, err := c.ReceiveMessage()
	if err != nil {
		return info, err
	}
	if len(msg.Arguments) >= 3 {
		info.Host = msg.Arguments[0].(string)
		info.Name = msg.Arguments[1].(string)
		info.Model = msg.Arguments[2].(string)
	}
	return info, nil
}

// KeepAlive sends keep-alive message (required for multi-client usage)
func (c *Client) KeepAlive() error {
	return c.SendMessage("/xremote")
}

// RequestStatus requests mixer status
func (c *Client) RequestStatus() error {
	return c.SendMessage("/status")
}
