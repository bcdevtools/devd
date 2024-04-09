package types

import (
	"context"
	"github.com/bcdevtools/devd/constants"
	"time"
)

type Context struct {
	baseCtx           context.Context
	operationUserInfo *OperationUserInfo
	workingUserInfo   *UserInfo
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

func NewContext(operationUserInfo *OperationUserInfo) Context {
	var workingUser *UserInfo
	if operationUserInfo != nil {
		workingUser = operationUserInfo.GetDefaultWorkingUser()
	}
	return Context{
		baseCtx:           context.Background(),
		operationUserInfo: operationUserInfo,
		workingUserInfo:   workingUser,
	}
}

func (c Context) GetOperationUserInfo() *OperationUserInfo {
	return c.operationUserInfo
}

func (c Context) GetWorkingUserInfo() *UserInfo {
	return c.workingUserInfo
}

func (c Context) WithContext(ctx context.Context) Context {
	c.baseCtx = ctx
	return c
}

func (c Context) WithWorkingUserInfo(workingUserInfo *UserInfo) Context {
	c.workingUserInfo = workingUserInfo
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
