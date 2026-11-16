package meta_test

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/etsevilcorp/toolware/meta"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

func TesMeta(t *testing.T) {
	k := rand.Text()
	v := rand.Text()
	m := mcp.NewMetaFromMap(map[string]any{k: v})
	ctx := meta.WithMeta(t.Context(), m)
	context.WithoutCancel(ctx)

	val, ok := meta.MetaFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, v, val)
}
