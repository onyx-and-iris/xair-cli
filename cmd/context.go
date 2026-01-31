package cmd

import (
	"context"

	"github.com/onyx-and-iris/xair-cli/internal/xair"
)

type clientKey string

// WithContext returns a new context with the provided xair.Client.
func WithContext(ctx context.Context, client *xair.Client) context.Context {
	return context.WithValue(ctx, clientKey("oscClient"), client)
}

// ClientFromContext retrieves the xair.Client from the context.
func ClientFromContext(ctx context.Context) *xair.Client {
	if client, ok := ctx.Value(clientKey("oscClient")).(*xair.Client); ok {
		return client
	}
	return nil
}
