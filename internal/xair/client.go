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

type XAirClient struct {
	conn      *net.UDPConn
	mixerAddr *net.UDPAddr

	parser parser

	done     chan bool
	respChan chan *osc.Message
}

// NewClient creates a new XAirClient instance
func NewClient(mixerIP string, mixerPort int) (*XAirClient, error) {
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

	return &XAirClient{
		conn:      conn,
		mixerAddr: mixerAddr,
		parser:    newParser(),
		done:      make(chan bool),
		respChan:  make(chan *osc.Message),
	}, nil
}

// Start begins listening for messages in a goroutine
func (x *XAirClient) StartListening() {
	go x.receiveLoop()
	log.Debugf("Started listening on %s...", x.conn.LocalAddr().String())
}

// receiveLoop handles incoming OSC messages
func (x *XAirClient) receiveLoop() {
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
func (x *XAirClient) parseOSCMessage(data []byte) (*osc.Message, error) {
	msg, err := x.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Stop stops the client and closes the connection
func (x *XAirClient) Stop() {
	close(x.done)
	if x.conn != nil {
		x.conn.Close()
	}
}

// SendMessage sends an OSC message to the mixer using the unified connection
func (x *XAirClient) SendMessage(address string, args ...any) error {
	return x.SendToAddress(x.mixerAddr, address, args...)
}

// SendToAddress sends an OSC message to a specific address (enables replying to different ports)
func (x *XAirClient) SendToAddress(addr *net.UDPAddr, oscAddress string, args ...any) error {
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
func (x *XAirClient) RequestInfo() (error, InfoResponse) {
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
func (x *XAirClient) KeepAlive() error {
	return x.SendMessage("/xremote")
}

// RequestStatus requests mixer status
func (x *XAirClient) RequestStatus() error {
	return x.SendMessage("/status")
}

/* STRIP METHODS */

// StripFader requests the current fader level for a strip
func (x *XAirClient) StripFader(strip int) (float64, error) {
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
func (x *XAirClient) SetStripFader(strip int, level float64) error {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	return x.SendMessage(address, float32(mustDbInto(level)))
}

// StripGain requests gain for a specific strip (1-based indexing)
func (x *XAirClient) StripGain(strip int) error {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	return x.SendMessage(address)
}

// SetStripGain sets gain for a specific strip (1-based indexing)
func (x *XAirClient) SetStripGain(strip int, gain float32) error {
	address := fmt.Sprintf("/ch/%02d/mix/fader", strip)
	return x.SendMessage(address, gain)
}

// StripMute gets mute state for a specific strip (1-based indexing)
func (x *XAirClient) StripMute(strip int) error {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	return x.SendMessage(address)
}

// SetStripMute sets mute state for a specific strip (1-based indexing)
func (x *XAirClient) SetStripMute(strip int, muted bool) error {
	address := fmt.Sprintf("/ch/%02d/mix/on", strip)
	var value int32 = 0
	if !muted {
		value = 1
	}
	return x.SendMessage(address, value)
}

// StripName requests the name for a specific strip
func (x *XAirClient) StripName(strip int) error {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	return x.SendMessage(address)
}

// SetStripName sets the name for a specific strip
func (x *XAirClient) SetStripName(strip int, name string) error {
	address := fmt.Sprintf("/ch/%02d/config/name", strip)
	return x.SendMessage(address, name)
}

// StripColor requests the color for a specific strip
func (x *XAirClient) StripColor(strip int) error {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	return x.SendMessage(address)
}

// SetStripColor sets the color for a specific strip (0-15)
func (x *XAirClient) SetStripColor(strip int, color int32) error {
	address := fmt.Sprintf("/ch/%02d/config/color", strip)
	return x.SendMessage(address, color)
}

// MainLRFader requests the current main L/R fader level
func (x *XAirClient) MainLRFader() (float64, error) {
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
func (x *XAirClient) SetMainLRFader(level float64) error {
	return x.SendMessage("/lr/mix/fader", float32(mustDbInto(level)))
}

// MainLRMute requests the current main L/R mute status
func (x *XAirClient) MainLRMute() (bool, error) {
	err := x.SendMessage("/lr/mix/on")
	if err != nil {
		return false, err
	}

	resp := <-x.respChan
	val, ok := resp.Arguments[0].(int32)
	if !ok {
		return false, fmt.Errorf("unexpected argument type for main LR mute value")
	}
	return val == 0, nil // 0 = muted, 1 = unmuted
}

// SetMainLRMute sets the main L/R mute status
func (x *XAirClient) SetMainLRMute(muted bool) error {
	var value int32 = 0
	if !muted {
		value = 1
	}
	return x.SendMessage("/lr/mix/on", value)
}
