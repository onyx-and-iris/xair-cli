package main

import "fmt"

type SnapshotCmdGroup struct {
	List  ListCmd `help:"List all snapshots."        cmd:"list"`
	Index struct {
		Index  int       `arg:"" help:"The index of the snapshot."`
		Name   NameCmd   `help:"Get or set the name of a snapshot."      cmd:"name"`
		Save   SaveCmd   `help:"Save the current mixer state." cmd:"save"`
		Load   LoadCmd   `help:"Load a mixer state."         cmd:"load"`
		Delete DeleteCmd `help:"Delete a snapshot."                      cmd:"delete"`
	} `help:"The index of the snapshot."            arg:""`
}

type ListCmd struct {
}

func (c *ListCmd) Run(ctx *context) error {
	for i := range 64 {
		name, err := ctx.Client.Snapshot.Name(i + 1)
		if err != nil {
			break
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

func (c *NameCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	if c.Name == nil {
		name, err := ctx.Client.Snapshot.Name(snapshot.Index.Index)
		if err != nil {
			return err
		}
		fmt.Fprintln(ctx.Out, name)
		return nil
	}

	return ctx.Client.Snapshot.SetName(snapshot.Index.Index, *c.Name)
}

type SaveCmd struct {
	Name string `arg:"" help:"The name of the snapshot."`
}

func (c *SaveCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	err := ctx.Client.Snapshot.CurrentName(c.Name)
	if err != nil {
		return err
	}

	return ctx.Client.Snapshot.CurrentSave(snapshot.Index.Index)
}

type LoadCmd struct {
}

func (c *LoadCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	return ctx.Client.Snapshot.CurrentLoad(snapshot.Index.Index)
}

type DeleteCmd struct {
}

func (c *DeleteCmd) Run(ctx *context, snapshot *SnapshotCmdGroup) error {
	return ctx.Client.Snapshot.CurrentDelete(snapshot.Index.Index)
}
