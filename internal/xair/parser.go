package xair

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/charmbracelet/log"
	"github.com/hypebeast/go-osc/osc"
)

type xairParser struct {
}

func newParser() *xairParser {
	return &xairParser{}
}

// parseOSCMessage parses raw bytes into an OSC message with improved error handling
func (p *xairParser) Parse(data []byte) (*osc.Message, error) {
	log.Debug("=== PARSING OSC MESSAGE BEGIN ===")
	defer log.Debug("=== PARSING OSC MESSAGE END ===")

	if err := p.validateOSCData(data); err != nil {
		return nil, err
	}

	address, addressEnd, err := p.extractOSCAddress(data)
	if err != nil {
		return nil, err
	}

	msg := osc.NewMessage(address)

	typeTags, typeTagsEnd, err := p.extractOSCTypeTags(data, addressEnd)
	if err != nil || typeTags == "" {
		log.Debug("No valid type tags, returning address-only message")
		return msg, nil
	}

	if err := p.parseOSCArguments(data, typeTagsEnd, typeTags, msg); err != nil {
		return nil, err
	}

	log.Debugf("Successfully parsed message with %d arguments", len(msg.Arguments))
	return msg, nil
}

// validateOSCData performs basic validation on OSC message data
func (p *xairParser) validateOSCData(data []byte) error {
	if len(data) < 4 {
		return fmt.Errorf("data too short for OSC message")
	}
	if data[0] != '/' {
		return fmt.Errorf("invalid OSC message: does not start with '/'")
	}
	return nil
}

// extractOSCAddress extracts the OSC address from the message data
func (p *xairParser) extractOSCAddress(data []byte) (address string, nextPos int, err error) {
	nullPos := bytes.IndexByte(data, 0)
	if nullPos <= 0 {
		return "", 0, fmt.Errorf("no null terminator found for address")
	}

	address = string(data[:nullPos])
	log.Debugf("Parsed OSC address: %s", address)

	// Calculate next 4-byte aligned position
	nextPos = ((nullPos + 4) / 4) * 4
	return address, nextPos, nil
}

// extractOSCTypeTags extracts and validates OSC type tags
func (p *xairParser) extractOSCTypeTags(data []byte, start int) (typeTags string, nextPos int, err error) {
	if start >= len(data) {
		return "", start, nil // No type tags available
	}

	typeTagsEnd := bytes.IndexByte(data[start:], 0)
	if typeTagsEnd <= 0 {
		return "", start, nil // No type tags found
	}

	typeTags = string(data[start : start+typeTagsEnd])
	log.Debugf("Parsed type tags: %s", typeTags)

	if len(typeTags) == 0 || typeTags[0] != ',' {
		log.Debug("Invalid type tags format")
		return "", start, nil
	}

	// Calculate arguments start position (4-byte aligned)
	nextPos = ((start + typeTagsEnd + 4) / 4) * 4
	return typeTags, nextPos, nil
}

// parseOSCArguments parses OSC arguments based on type tags
func (p *xairParser) parseOSCArguments(data []byte, argsStart int, typeTags string, msg *osc.Message) error {
	argData := data[argsStart:]
	argNum := 0

	for i := 1; i < len(typeTags) && len(argData) > 0; i++ {
		var consumed int
		var err error

		switch typeTags[i] {
		case 's':
			consumed, err = p.parseStringArgument(argData, msg, argNum)
		case 'i':
			consumed, err = p.parseInt32Argument(argData, msg, argNum)
		case 'f':
			consumed, err = p.parseFloat32Argument(argData, msg, argNum)
		case 'b':
			consumed, err = p.parseBlobArgument(argData, msg, argNum)
		default:
			log.Debugf("Unknown type tag: %c (skipping)", typeTags[i])
			consumed = p.skipUnknownArgument(argData)
		}

		if err != nil {
			log.Debugf("Error parsing argument %d: %v", argNum+1, err)
			break
		}

		if consumed == 0 {
			break // No more data to consume
		}

		argData = argData[consumed:]
		if typeTags[i] != '?' { // Don't count skipped arguments
			argNum++
		}
	}

	return nil
}

// parseStringArgument parses a string argument from OSC data
func (p *xairParser) parseStringArgument(data []byte, msg *osc.Message, argNum int) (int, error) {
	nullPos := bytes.IndexByte(data, 0)
	if nullPos < 0 {
		return 0, fmt.Errorf("no null terminator found for string")
	}

	argStr := string(data[:nullPos])
	log.Debugf("Parsed string argument %d: %s", argNum+1, argStr)
	msg.Append(argStr)

	// Return next 4-byte aligned position
	return ((nullPos + 4) / 4) * 4, nil
}

// parseInt32Argument parses an int32 argument from OSC data
func (p *xairParser) parseInt32Argument(data []byte, msg *osc.Message, argNum int) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data for int32")
	}

	val := int32(binary.BigEndian.Uint32(data[:4]))
	log.Debugf("Parsed int32 argument %d: %d", argNum+1, val)
	msg.Append(val)

	return 4, nil
}

// parseFloat32Argument parses a float32 argument from OSC data
func (p *xairParser) parseFloat32Argument(data []byte, msg *osc.Message, argNum int) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data for float32")
	}

	val := math.Float32frombits(binary.BigEndian.Uint32(data[:4]))
	log.Debugf("Parsed float32 argument %d: %f", argNum+1, val)
	msg.Append(val)

	return 4, nil
}

// parseBlobArgument parses a blob argument from OSC data
func (p *xairParser) parseBlobArgument(data []byte, msg *osc.Message, argNum int) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data for blob size")
	}

	size := int32(binary.BigEndian.Uint32(data[:4]))
	if size < 0 || size >= 10000 {
		return 0, fmt.Errorf("invalid blob size: %d", size)
	}

	if len(data) < int(4+size) {
		return 0, fmt.Errorf("insufficient data for blob content")
	}

	blob := make([]byte, size)
	copy(blob, data[4:4+size])
	log.Debugf("Parsed blob argument %d (%d bytes)", argNum+1, size)
	msg.Append(blob)

	// Return next 4-byte aligned position
	return ((4 + int(size) + 3) / 4) * 4, nil
}

// skipUnknownArgument skips an unknown argument type
func (p *xairParser) skipUnknownArgument(data []byte) int {
	// Skip unknown types by moving 4 bytes if available
	if len(data) >= 4 {
		return 4
	}
	return 0
}
