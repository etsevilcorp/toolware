package toolware

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
)

type Request struct {
	ToolName string
	Args     map[string]any // writing any feels like torture
}

func (r Request) Int(key string) (int, error) {
	v, ok := r.Args[key]
	if !ok {
		return 0, fmt.Errorf("toolware: missing arg %q", key)
	}
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	case json.Number:
		i, err := n.Int64()
		return int(i), err
	default:
		return 0, fmt.Errorf("toolware: arg %q is not numeric(%T)", key, v)
	}
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
