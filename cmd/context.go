package cmd

import (
	"context"

	"github.com/onyx-and-iris/xair-cli/internal/xair"
)

type clientKey string

func WithContext(ctx context.Context, client *xair.XAirClient) context.Context {
	return context.WithValue(ctx, clientKey("oscClient"), client)
}

func ClientFromContext(ctx context.Context) *xair.XAirClient {
	if client, ok := ctx.Value(clientKey("oscClient")).(*xair.XAirClient); ok {
		return client
	}
	return nil
}
