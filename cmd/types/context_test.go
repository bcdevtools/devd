package types

import (
	"testing"
)

func TestNewContext(t *testing.T) {
	appCtx := NewContext(nil)
	if appCtx.operationUserInfo != nil {
		t.Errorf("NewContext() = want operationUserInfo == nil")
		return
	}
	if appCtx.workingUserInfo != nil {
		t.Errorf("NewContext() = want workingUserInfo == nil")
		return
	}

	appCtx = NewContext(&OperationUserInfo{
		EffectiveUserInfo: &UserInfo{},
		RealUserInfo:      &UserInfo{},
	})
	if appCtx.operationUserInfo == nil {
		t.Errorf("NewContext() = want operationUserInfo != nil")
		return
	}
	if appCtx.workingUserInfo == nil {
		t.Errorf("NewContext() = want workingUserInfo != nil")
		return
	}
}

func TestWrapAppContext(t *testing.T) {
	appCtx := NewContext(nil)
	ctx := WrapAppContext(appCtx)
	if ctx == nil {
		t.Errorf("WrapAppContext() = nil, want not nil")
		return
	}
	if ctx != appCtx {
		t.Errorf("WrapAppContext() = want the same context as original")
		return
	}
}

func TestUnwrapAppContext(t *testing.T) {
	appCtx := NewContext(nil)
	if appCtx.operationUserInfo != nil {
		panic("appCtx.operationUserInfo != nil")
	}
	ctx := WrapAppContext(appCtx)
	appCtx2 := UnwrapAppContext(ctx)
	if appCtx != appCtx2 {
		t.Errorf("UnwrapAppContext() = want the same app context as original")
		return
	}
	if appCtx.baseCtx != appCtx2.baseCtx {
		t.Errorf("UnwrapAppContext() = want the same base context as original")
		return
	}
	if appCtx.operationUserInfo != appCtx2.operationUserInfo {
		t.Errorf("UnwrapAppContext() = want the same operationUserInfo as original")
		return
	}
	if appCtx.operationUserInfo != appCtx2.operationUserInfo {
		t.Errorf("UnwrapAppContext() = want the same operationUserInfo as original")
		return
	}

	//

	appCtx = NewContext(&OperationUserInfo{})
	if appCtx.operationUserInfo == nil {
		panic("appCtx.operationUserInfo == nil")
	}
	ctx = WrapAppContext(appCtx)
	appCtx2 = UnwrapAppContext(ctx)
	if appCtx.operationUserInfo != appCtx2.operationUserInfo {
		t.Errorf("UnwrapAppContext() = want the same operationUserInfo as original")
		return
	}
}
