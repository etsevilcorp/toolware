package toolware

import (
	"context"
	"slices"
)

type Request struct {
	Toolname string
	Args     map[string]any // writing any feels like torture
}

type Response struct {
	Result any
	Err    error
}

type Handler func(ctx context.Context, req Request) Response

type Middleware func(Handler) Handler

func chain(mws []Middleware, final Handler) Handler {
	h := final
	for i, _ := range slices.Backward(mws) {
		h = mws[i](h)
	}
	return h
}
