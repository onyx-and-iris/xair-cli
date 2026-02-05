package main

import (
	"fmt"
	"time"
)

// RawCmd represents the command to send raw OSC messages to the mixer.
type RawCmd struct {
	Timeout time.Duration `help:"Timeout for the OSC message send operation."  default:"200ms" short:"t"`
	Address string        `help:"The OSC address to send the message to."                                arg:""`
	Args    []string      `help:"The arguments to include in the OSC message."                           arg:"" optional:""`
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

	msg, err := ctx.Client.ReceiveMessage(cmd.Timeout)
	if err != nil {
		return fmt.Errorf("failed to receive response for raw OSC message: %w", err)
	}
	if msg != nil {
		fmt.Fprintf(ctx.Out, "Received response: %s with args: %v\n", msg.Address, msg.Arguments)
	}

	return nil
}
