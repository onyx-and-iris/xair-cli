/*
LICENSE: https://github.com/onyx-and-iris/xair-cli/blob/main/LICENSE
*/
package cmd

import (
	"context"

	"github.com/onyx-and-iris/xair-cli/internal/xair"
)

type clientKey string

func WithContext(ctx context.Context, client *xair.Client) context.Context {
	return context.WithValue(ctx, clientKey("oscClient"), client)
}

func ClientFromContext(ctx context.Context) *xair.Client {
	if client, ok := ctx.Value(clientKey("oscClient")).(*xair.Client); ok {
		return client
	}
	return nil
}
