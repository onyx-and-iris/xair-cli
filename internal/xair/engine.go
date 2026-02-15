package xair

import (
	"errors"
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
	Kind      mixerKind
	timeout   time.Duration
	conn      *net.UDPConn
	mixerAddr *net.UDPAddr

	parser     parser
	addressMap map[string]string

	done     chan bool
	respChan chan *osc.Message
}

func newEngine(
	mixerIP string,
	mixerPort int,
	kind mixerKind,
	opts ...EngineOption,
) (*engine, error) {
	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", 0))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local address: %w", err)
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	mixerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", mixerIP, mixerPort))
	if err != nil {
		if err := conn.Close(); err != nil {
			log.Errorf("failed to close UDP connection: %v", err)
		}
		return nil, fmt.Errorf("failed to resolve mixer address: %w", err)
	}

	log.Debugf("Local UDP connection: %s	", conn.LocalAddr().String())

	e := &engine{
		timeout:    100 * time.Millisecond,
		conn:       conn,
		mixerAddr:  mixerAddr,
		parser:     newParser(),
		addressMap: addressMapFromMixerKind(kind),
		done:       make(chan bool),
		respChan:   make(chan *osc.Message, 100),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

// receiveLoop handles incoming OSC messages.
func (e *engine) receiveLoop() {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-e.done:
			return
		default:
			// Set a short read deadline to prevent blocking indefinitely
			if err := e.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond)); err != nil {
				log.Errorf("Failed to set read deadline: %v", err)
				continue
			}
			n, _, err := e.conn.ReadFromUDP(buffer)
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) {
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

// parseOSCMessage parses raw bytes into an OSC message with improved error handling.
func (e *engine) parseOSCMessage(data []byte) (*osc.Message, error) {
	msg, err := e.parser.Parse(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// sendToAddress sends an OSC message to a specific address (enables replying to different ports).
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
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = e.conn.WriteToUDP(data, addr)
	return err
}
