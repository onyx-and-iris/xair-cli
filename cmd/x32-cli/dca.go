package main

import (
	"fmt"
)

type DCACmdGroup struct {
	Index struct {
		Index int        `arg:"" help:"The index of the DCA group (1-8)."`
		Mute  DCAMuteCmd `help:"Get or set the mute status of the DCA group." cmd:""`
		Name  DCANameCmd `help:"Get or set the name of the DCA group." cmd:""`
	} `arg:"" help:"Control a specific DCA group by its index."`
}

// Validate checks if the provided index is within the valid range.
func (cmd *DCACmdGroup) Validate() error {
	if cmd.Index.Index < 1 || cmd.Index.Index > 8 {
		return fmt.Errorf("DCA group index must be between 1 and 8, got %d", cmd.Index.Index)
	}
	return nil
}

// DCAMuteCmd is the command to get or set the mute status of a DCA group.
type DCAMuteCmd struct {
	State *string `arg:"" help:"Set the mute status of the DCA group." optional:"" enum:"true,false"`
}

// Run executes the DCAMuteCmd command.
func (cmd *DCAMuteCmd) Run(ctx *context, dca *DCACmdGroup) error {
	if cmd.State == nil {
		resp, err := ctx.Client.DCA.Mute(dca.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get DCA mute status: %w", err)
		}
		fmt.Fprintf(ctx.Out, "DCA Group %d mute state: %t\n", dca.Index.Index, resp)
		return nil
	}

	if err := ctx.Client.DCA.SetMute(dca.Index.Index, *cmd.State == trueStr); err != nil {
		return fmt.Errorf("failed to set DCA mute status: %w", err)
	}
	return nil
}

// DCANameCmd is the command to get or set the name of a DCA group.
type DCANameCmd struct {
	Name *string `arg:"" help:"Set the name of the DCA group." optional:""`
}

// Run executes the DCANameCmd command.
func (cmd *DCANameCmd) Run(ctx *context, dca *DCACmdGroup) error {
	if cmd.Name == nil {
		resp, err := ctx.Client.DCA.Name(dca.Index.Index)
		if err != nil {
			return fmt.Errorf("failed to get DCA name: %w", err)
		}
		if resp == "" {
			resp = fmt.Sprintf("DCA %d", dca.Index.Index)
		}
		fmt.Fprintf(ctx.Out, "DCA Group %d is named '%s'\n", dca.Index.Index, resp)
		return nil
	}
	if err := ctx.Client.DCA.SetName(dca.Index.Index, *cmd.Name); err != nil {
		return fmt.Errorf("failed to set DCA name: %w", err)
	}
	return nil
}
