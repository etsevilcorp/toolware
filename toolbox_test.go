package toolware_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/etsevilcorp/toolware"
	"github.com/stretchr/testify/require"
)

func TestEnchainer_Chainer(t *testing.T) {
	t.Run("chain", func(t *testing.T) {
		var order []string
		track := func(name string) toolware.Middleware {
			return func(next toolware.Handler) toolware.Handler {
				return func(ctx context.Context, req toolware.Request) toolware.Response {
					order = append(order, name+":in")
					resp := next(ctx, req)
					order = append(order, name+":out")
					return resp
				}
			}
		}
		final := func(ctx context.Context, req toolware.Request) toolware.Response {
			order = append(order, "handler")
			return toolware.Response{Result: "ok"}
		}

		h := toolware.Chain([]toolware.Middleware{track("a"), track("b")}, final)
		resp := h(context.Background(), toolware.Request{})
		require.NoError(t, resp.Err)

		want := []string{"a:in", "b:in", "handler", "b:out", "a:out"}
		if !reflect.DeepEqual(order, want) {
			t.Fatalf("got %v, want %v", order, want)
		}

		require.Equal(t, want, order)
	})
	t.Run("premature bubbling", func(t *testing.T) {
		called := false
		deny := func(next toolware.Handler) toolware.Handler {
			return func(ctx context.Context, req toolware.Request) toolware.Response {
				return toolware.Response{Err: errors.New("stop drilling!")}
			}
		}
		final := func(ctx context.Context, req toolware.Request) toolware.Response {
			called = true
			return toolware.Response{}
		}
		resp := toolware.Chain([]toolware.Middleware{deny}, final)(context.Background(), toolware.Request{})
		require.Error(t, resp.Err)
		require.False(t, called)
	})
}

func TestWith(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var order []string

		tw := toolware.New()
		tw.With(func(name string) toolware.Middleware {
			return func(next toolware.Handler) toolware.Handler {
				return func(ctx context.Context, req toolware.Request) toolware.Response {
					order = append(order, name+":in")
					resp := next(ctx, req)
					order = append(order, name+":out")
					return resp
				}
			}
		}("a")).Register("final", func(ctx context.Context, _ int) (string, error) {
			order = append(order, "handler")
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)

		h(context.Background(), toolware.Request{})

		want := []string{"a:in", "handler", "a:out"}

		require.Equal(t, want, order)
	})
	t.Run("use", func(t *testing.T) {
		var order []string

		track := func(name string) toolware.Middleware {
			return func(next toolware.Handler) toolware.Handler {
				return func(ctx context.Context, req toolware.Request) toolware.Response {
					order = append(order, name+":in")
					resp := next(ctx, req)
					order = append(order, name+":out")
					return resp
				}
			}
		}

		tw := toolware.New()

		tw.Use(track("a"))

		tw.With(track("b")).Register("final", func(ctx context.Context, _ int) (string, error) {
			order = append(order, "handler")
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)

		h(context.Background(), toolware.Request{})

		want := []string{"a:in", "b:in", "handler", "b:out", "a:out"}

		require.Equal(t, want, order)
	})
	t.Run("empty", func(t *testing.T) {
		var order []string

		tw := toolware.New()
		tw.With().Register("final", func(ctx context.Context, _ int) (string, error) {
			order = append(order, "handler")
			return "ok", nil
		}) // Give me a pizza with extra cheese... and hold the tomato sauce... and hold the crust!

		h, ok := tw.Get("final") // Hey Jimmy, give me a cheese with nothing!
		require.True(t, ok)

		h(context.Background(), toolware.Request{})

		want := []string{"handler"}

		require.Equal(t, want, order)
	})

	t.Run("premature bubbling", func(t *testing.T) {
		called := false
		deny := func(next toolware.Handler) toolware.Handler {
			return func(ctx context.Context, req toolware.Request) toolware.Response {
				return toolware.Response{Err: errors.New("stop drilling!")}
			}
		}
		final := func(ctx context.Context, _ toolware.Request) (int, error) {
			called = true
			return -1, nil
		}
		tw := toolware.New()
		tw.With(deny).Register("denied", final)

		h, ok := tw.Get("denied")
		require.True(t, ok)

		resp := h(t.Context(), toolware.Request{})
		require.Error(t, resp.Err)
		require.False(t, called)
	})
}

func TestNamespace(t *testing.T) {
	t.Run("namespacing", func(t *testing.T) {
		tw := toolware.New()
		tw.Namespace("namespace_", func(t *toolware.Toolbox) {
			t.Register("final", func(ctx context.Context, _ int) (string, error) {
				return "ok", nil
			})
		})

		h, ok := tw.Get("namespace_final")
		require.True(t, ok)
		resp := h(t.Context(), toolware.Request{})

		require.NoError(t, resp.Err)
		require.Equal(t, "ok", resp.Result)
	})
	t.Run("outside naming", func(t *testing.T) {
		tw := toolware.New()
		tw.Namespace("namespace_", func(t *toolware.Toolbox) {
			t.Register("final", func(ctx context.Context, _ int) (string, error) {
				return "not ok", nil
			})
		})

		tw.Register("final", func(ctx context.Context, _ int) (string, error) {
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)
		resp := h(t.Context(), toolware.Request{})

		require.NoError(t, resp.Err)
		require.Equal(t, "ok", resp.Result)
	})
	t.Run("nothing", func(t *testing.T) {
		tw := toolware.New()
		tw.Namespace("namespace_", func(t *toolware.Toolbox) {
		})

		tw.Register("final", func(ctx context.Context, _ int) (string, error) {
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)
		resp := h(t.Context(), toolware.Request{})

		require.NoError(t, resp.Err)
		require.Equal(t, "ok", resp.Result)
	})
	t.Run("with mws", func(t *testing.T) {
		var order []string

		track := func(name string) toolware.Middleware {
			return func(next toolware.Handler) toolware.Handler {
				return func(ctx context.Context, req toolware.Request) toolware.Response {
					order = append(order, name+":in")
					resp := next(ctx, req)
					order = append(order, name+":out")
					return resp
				}
			}
		}

		tw := toolware.New()

		tw.Namespace("ns_", func(t *toolware.Toolbox) {
			t.With(track("b")).Register("final", func(ctx context.Context, _ int) (string, error) {
				order = append(order, "handler")
				return "ok", nil
			})
		})

		h, ok := tw.Get("ns_final")
		require.True(t, ok)

		h(context.Background(), toolware.Request{})

		want := []string{"b:in", "handler", "b:out"}

		require.Equal(t, want, order)
	})
}

func TestGroup(t *testing.T) {
	t.Run("grouping", func(t *testing.T) {
		tw := toolware.New()
		tw.Group(func(t *toolware.Toolbox) {
			t.Register("final", func(ctx context.Context, _ int) (string, error) {
				return "ok", nil
			})
		})

		h, ok := tw.Get("final")
		require.True(t, ok)
		resp := h(t.Context(), toolware.Request{})

		require.NoError(t, resp.Err)
		require.Equal(t, "ok", resp.Result)
	})

	t.Run("nothing", func(t *testing.T) {
		tw := toolware.New()
		tw.Group(func(t *toolware.Toolbox) {
		})

		tw.Register("final", func(ctx context.Context, _ int) (string, error) {
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)
		resp := h(t.Context(), toolware.Request{})

		require.NoError(t, resp.Err)
		require.Equal(t, "ok", resp.Result)
	})
	t.Run("with mws", func(t *testing.T) {
		var order []string

		track := func(name string) toolware.Middleware {
			return func(next toolware.Handler) toolware.Handler {
				return func(ctx context.Context, req toolware.Request) toolware.Response {
					order = append(order, name+":in")
					resp := next(ctx, req)
					order = append(order, name+":out")
					return resp
				}
			}
		}

		tw := toolware.New()

		tw.Group(func(t *toolware.Toolbox) {
			t.With(track("b")).Register("final", func(ctx context.Context, _ int) (string, error) {
				order = append(order, "handler")
				return "ok", nil
			})
		})

		h, ok := tw.Get("final")
		require.True(t, ok)

		h(context.Background(), toolware.Request{})

		want := []string{"b:in", "handler", "b:out"}

		require.Equal(t, want, order)
	})
}

func TestRegister(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tw := toolware.New()
		tw.Register("final", func(ctx context.Context, _ int) (string, error) {
			return "ok", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)

		resp := h(context.Background(), toolware.Request{})
		require.NoError(t, resp.Err)

		require.Equal(t, "ok", resp.Result)
	})
	t.Run("bad args", func(t *testing.T) {
		req := toolware.Request{Args: []byte("{is bad JSON")}

		tw := toolware.New()
		tw.Register("final", func(ctx context.Context, _ int) (string, error) {
			return "huh?", nil
		})

		h, ok := tw.Get("final")
		require.True(t, ok)

		resp := h(context.Background(), req)
		require.Error(t, resp.Err)

		require.NotEqual(t, "huh?", resp.Result)
	})
}
