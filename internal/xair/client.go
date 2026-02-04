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

type Client struct {
	engine
	Main     *Main
	Strip    *Strip
	Bus      *Bus
	HeadAmp  *HeadAmp
	Snapshot *Snapshot
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
	c.Main = newMain(c)
	c.Strip = NewStrip(c)
	c.Bus = NewBus(c)
	c.HeadAmp = NewHeadAmp(c)

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

// ReceiveMessage receives an OSC message from the mixer
func (c *Client) ReceiveMessage(timeout time.Duration) (*osc.Message, error) {
	t := time.Tick(timeout)
	select {
	case <-t:
		return nil, nil
	case val := <-c.respChan:
		if val == nil {
			return nil, fmt.Errorf("no message received")
		}
		return val, nil
	}
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
