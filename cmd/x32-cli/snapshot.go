package main

import "fmt"

type SnapshotCmdGroup struct {
	List  ListCmd `cmd:"" help:"List all snapshots."`
	Index struct {
		Index  *int      `arg:"" help:"The index of the snapshot." optional:""`
		Name   NameCmd   `cmd:"" help:"Get or set the name of a snapshot."`
		Save   SaveCmd   `cmd:"" help:"Save the current mixer state to a snapshot."`
		Load   LoadCmd   `cmd:"" help:"Load a mixer state from a snapshot."`
		Delete DeleteCmd `cmd:"" help:"Delete a snapshot."`
	} `       help:"The index of the snapshot." arg:""`
}

// Validate checks if the provided snapshot index is within the valid range (1-64) when any of the subcommands that require an index are used.
func (cmd *SnapshotCmdGroup) Validate() error {
	if cmd.Index.Index == nil {
		return nil
	}

	if *cmd.Index.Index < 1 || *cmd.Index.Index > 64 {
		return fmt.Errorf("snapshot index must be between 1 and 64")
	}

	return nil
}

type ListCmd struct{}

func (cmd *ListCmd) Run(ctx *context) error {
	for i := range 64 {
		name, err := ctx.Client.Snapshot.Name(i + 1)
		if err != nil {
			return fmt.Errorf("failed to get name for snapshot %d: %w", i+1, err)
		}
		if name == "" {
			continue
		}
		fmt.Fprintf(ctx.Out, "%d: %s\n", i+1, name)
	}
	return nil
}

type NameCmd struct {
	Name *string `arg:"" help:"The name of the snapshot." optional:""`
}

func (cmd *NameCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	if cmd.Name == nil {
		name, err := ctx.Client.Snapshot.Name(*snapshot.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintln(ctx.Out, name)
		return nil
	}

	return ctx.Client.Snapshot.SetName(*snapshot.Index.Index, *cmd.Name)
}

type SaveCmd struct {
	Name string `arg:"" help:"The name of the snapshot."`
}

func (cmd *SaveCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	err := ctx.Client.Snapshot.CurrentName(cmd.Name)
	if err != nil {
		return err
	}

	return ctx.Client.Snapshot.CurrentSave(*snapshot.Index.Index)
}

type LoadCmd struct{}

func (cmd *LoadCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	return ctx.Client.Snapshot.CurrentLoad(*snapshot.Index.Index)
}

type DeleteCmd struct{}

func (cmd *DeleteCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	return ctx.Client.Snapshot.CurrentDelete(*snapshot.Index.Index)
}
