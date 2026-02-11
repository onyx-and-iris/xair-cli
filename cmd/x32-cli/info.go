package main

import "fmt"

type InfoCmd struct {
}

func (c *InfoCmd) Run(ctx *context) error {
	fmt.Fprintf(
		ctx.Out,
		"Host: %s | Name: %s | Model: %s | Firmware: %s\n",
		ctx.Client.Info.Host,
		ctx.Client.Info.Name,
		ctx.Client.Info.Model,
		ctx.Client.Info.Firmware,
	)
	return nil
}
