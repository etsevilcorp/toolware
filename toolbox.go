package toolware

import (
	"context"
	"encoding/json"
	"fmt"
)

type Toolbox struct {
	middlewares []Middleware
	registry    map[string]Handler
}

func New() *Toolbox {
	return &Toolbox{registry: make(map[string]Handler)}
}

func (t *Toolbox) Use(mw Middleware) {
	t.middlewares = append(t.middlewares, mw)
}

func (t *Toolbox) With(mws ...Middleware) *Toolbox {
	combined := make([]Middleware, 0, len(t.middlewares)+len(mws))

	combined = append(combined, t.middlewares...)
	combined = append(combined, mws...)

	return &Toolbox{middlewares: combined, registry: t.registry}
}

func (t *Toolbox) Group(fn func(*Toolbox)) *Toolbox {
	child := t.With()
	fn(child)
	return child
}

func (t *(Toolbox)) Register[I, O any](name string, fn func(context.Context, I) (O, error)) {
	t.registry[name] = chain(t.middlewares, func(ctx context.Context, req Request) Response {
		var in I
		if len(req.Args) > 0 {
			if err := json.Unmarshal(req.Args, &in); err != nil {
				return Response{Err: fmt.Errorf("toolware: tool %q: bad args: %w", name, err)}
			}
		}
		out, err := fn(ctx, in)
		return Response{Result: out, Err: err}
	})
}

func (t *Toolbox) Get(name string) (Handler, bool) {
	h, ok := t.registry[name]
	return h, ok
}
