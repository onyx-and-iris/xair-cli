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
	timeout   time.Duration
	conn      *net.UDPConn
	mixerAddr *net.UDPAddr

	parser     parser
	addressMap map[string]string

	done     chan bool
	respChan chan *osc.Message
}

// receiveLoop handles incoming OSC messages
func (e *engine) receiveLoop() {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-e.done:
			return
		default:
			// Set a short read deadline to prevent blocking indefinitely
			e.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, _, err := e.conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Timeout is expected, continue loop
					continue
				}
				// Check if we're shutting down to avoid logging expected errors
				select {
				case <-e.done:
					return
				default:
					log.Errorf("Read error: %v", err)
					return
				}
			}

			msg, err := e.parseOSCMessage(buffer[:n])
			if err != nil {
				log.Errorf("Failed to parse OSC message: %v", err)
				continue
			}
			e.respChan <- msg
		}
	}
}

// parseOSCMessage parses raw bytes into an OSC message with improved error handling
func (e *engine) parseOSCMessage(data []byte) (*osc.Message, error) {
	msg, err := e.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// sendToAddress sends an OSC message to a specific address (enables replying to different ports)
func (e *engine) sendToAddress(addr *net.UDPAddr, oscAddress string, args ...any) error {
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

	_, err = e.conn.WriteToUDP(data, addr)
	return err
}
