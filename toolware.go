package toolware

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
)

type Request struct {
	ToolName string
	Args     json.RawMessage
}

func (r *Request) Get[T any](key string) (result T, err error) {
	var m map[string]any
	if err := json.Unmarshal(r.Args, &m); err != nil {
		return result, fmt.Errorf("toolware: failed to unmarshal args: %w", err)
	}

	v, ok := m[key]
	if !ok {
		return result, fmt.Errorf("toolware: argument with key %q was not found", key)
	}

	result, ok = v.(T)
	if !ok {
		return result, fmt.Errorf("toolware: given argument is wrongly typed, want %T, got %T", result, v)
	}

	return result, nil
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
