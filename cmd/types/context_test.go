package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewContext(t *testing.T) {
	_ = NewContext()
}

func TestWrapAppContext(t *testing.T) {
	appCtx := NewContext()
	ctx := WrapAppContext(appCtx)
	require.NotNil(t, ctx)
	require.Equal(t, appCtx, ctx, "want the same context as original")
}

func TestUnwrapAppContext(t *testing.T) {
	appCtx := NewContext()
	ctx := WrapAppContext(appCtx)
	appCtx2 := UnwrapAppContext(ctx)
	require.Equal(t, appCtx, appCtx2, "want the same app context as original")
	require.Equal(t, appCtx.baseCtx, appCtx2.baseCtx, "want the same base context as original")
}
