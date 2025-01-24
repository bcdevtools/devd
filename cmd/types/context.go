package types

import (
	"context"
	"time"

	"github.com/bcdevtools/devd/v3/constants"
)

type Context struct {
	baseCtx context.Context
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.baseCtx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.baseCtx.Done()
}

func (c Context) Err() error {
	return c.baseCtx.Err()
}

func (c Context) Value(key any) any {
	if key == AppContextKey {
		return c
	}

	return c.baseCtx.Value(key)
}

func NewContext() Context {
	return Context{
		baseCtx: context.Background(),
	}
}

func (c Context) WithContext(ctx context.Context) Context {
	c.baseCtx = ctx
	return c
}

var _ context.Context = Context{}

// ContextKey defines a type alias for a stdlib Context key.
type ContextKey string

// AppContextKey is the key in the context.Context which holds the application context.
const AppContextKey = constants.BINARY_NAME + "-app-context"

func WrapAppContext(ctx Context) context.Context {
	return ctx
}

func UnwrapAppContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	return ctx.Value(AppContextKey).(Context)
}
