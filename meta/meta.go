package meta

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type metaKey struct{}

func WithMeta(ctx context.Context, meta mcp.Meta) context.Context {
	return context.WithValue(ctx, metaKey{}, meta)
}

func MetaFromContext(ctx context.Context) (mcp.Meta, bool) {
	v, ok := ctx.Value(metaKey{}).(mcp.Meta)
	return v, ok
}
