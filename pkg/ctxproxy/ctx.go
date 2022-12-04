package ctxproxy

import (
	"context"
)

var ctxEP0Key int

func NewCtxEP0Data(ctx context.Context, v *InfoCtxEP0Data) context.Context {
	return context.WithValue(ctx, ctxEP0Key, v)
}

func CtxEP0Data(ctx context.Context) *InfoCtxEP0Data {
	u, ok := ctx.Value(ctxEP0Key).(*InfoCtxEP0Data)
	if !ok {
		return nil
	}
	return u
}
