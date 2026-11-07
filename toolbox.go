package toolware

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

func (t *(Toolbox)) Register(name string, h Handler) {
	t.registry[name] = chain(t.middlewares, h)
}

func (t *Toolbox) Get(name string) (Handler, bool) {
	h, ok := t.registry[name]
	return h, ok
}
