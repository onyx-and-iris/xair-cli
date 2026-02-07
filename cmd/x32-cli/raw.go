package main

import (
	"fmt"

	"github.com/charmbracelet/log"
)

// RawCmd represents the command to send raw OSC messages to the mixer.
type RawCmd struct {
	Address string   `help:"The OSC address to send the message to."      arg:""`
	Args    []string `help:"The arguments to include in the OSC message." arg:"" optional:""`
}

// Run executes the RawCmd by sending the specified OSC message to the mixer and optionally waiting for a response.
func (cmd *RawCmd) Run(ctx *context) error {
	params := make([]any, len(cmd.Args))
	for i, arg := range cmd.Args {
		params[i] = arg
	}
	if err := ctx.Client.SendMessage(cmd.Address, params...); err != nil {
		return fmt.Errorf("failed to send raw OSC message: %w", err)
	}

	if len(params) > 0 {
		log.Debugf("Sent OSC message: %s with args: %v\n", cmd.Address, cmd.Args)
		return nil
	}

	msg, err := ctx.Client.ReceiveMessage()
	if err != nil {
		return fmt.Errorf("failed to receive response for raw OSC message: %w", err)
	}
	if msg != nil {
		fmt.Fprintf(ctx.Out, "Received response: %s with args: %v\n", msg.Address, msg.Arguments)
	}

	return nil
}
