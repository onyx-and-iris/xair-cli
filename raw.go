package main

import (
	"fmt"
	"time"
)

type RawCmdGroup struct {
	Send RawSendCmd `help:"Send a raw OSC message to the mixer." cmd:""`
}

type RawSendCmd struct {
	Timeout time.Duration `help:"Timeout for the OSC message send operation."  default:"200ms" short:"t"`
	Address string        `help:"The OSC address to send the message to."                                arg:""`
	Args    []string      `help:"The arguments to include in the OSC message."                           arg:"" optional:""`
}

func (cmd *RawSendCmd) Run(ctx *context) error {
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
	fmt.Fprintf(ctx.Out, "Received response: %s with args: %v\n", msg.Address, msg.Arguments)

	return nil
}
