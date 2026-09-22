package toolware_test

import (
	"encoding/json"
	"math/rand/v2"
	"testing"

	"github.com/etsevilcorp/toolware"
	"github.com/stretchr/testify/require"
)

func TestRequest_Int(t *testing.T) {
	cases := []struct {
		name      string
		args      map[string]any
		key       string
		want      int
		shouldErr bool
	}{
		{"float64", map[string]any{"n": float64(64)}, "n", 64, false},
		{"int", map[string]any{"n": int(739)}, "n", 739, false}, // it took too much time to express "depends" as sum of chars
		{"float32", map[string]any{"n": float32(32)}, "n", 32, false},
		{"json.Number", map[string]any{"n": json.Number("649")}, "n", 649, false}, // same as before(says "number", quite fascinating, isn't it)
		{"missing everything", map[string]any{}, "n", -1, true},
		{"missing not empty", map[string]any{"k": rand.Int()}, "n", -1, true},
		{"wrong type", map[string]any{"n": "can you really say that it's that different?"}, "n", -1, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := toolware.Request{Args: c.args}
			v, err := req.Int(c.key)
			if c.shouldErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, v)
			}
		})
	}
}

func TestRequest_Boolean(t *testing.T) {
	cases := []struct {
		name      string
		args      map[string]any
		key       string
		want      bool
		shouldErr bool
	}{
		{"bool", map[string]any{"b": bool(true)}, "b", true, false},
		{"missing everything", map[string]any{}, "b", false, true},
		{"missing not empty", map[string]any{"a": rand.Int()}, "b", false, true},
		{"wrong type", map[string]any{"b": "can you really say that it's that different?"}, "b", false, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := toolware.Request{Args: c.args}
			v, err := req.Boolean(c.key)
			if c.shouldErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, v)
			}
		})
	}
}

func TestRequest_String(t *testing.T) {
	cases := []struct {
		name      string
		args      map[string]any
		key       string
		want      string
		shouldErr bool
	}{
		{"string", map[string]any{"s": "string"}, "s", "string", false},
		{"missing everything", map[string]any{}, "s", "", true},
		{"missing not empty", map[string]any{"a": rand.Float64()}, "b", "", true},
		{"wrong type", map[string]any{"b": 2}, "b", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := toolware.Request{Args: c.args}
			v, err := req.String(c.key)
			if c.shouldErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, v)
			}
		})
	}
}

func TestRequest_Float(t *testing.T) {
	cases := []struct {
		name      string
		args      map[string]any
		key       string
		want      float64
		shouldErr bool
	}{
		{"float64", map[string]any{"n": float64(64.)}, "n", 64., false},
		{"int", map[string]any{"n": int(739)}, "n", 739., false},
		{"float32", map[string]any{"n": float32(32)}, "n", 32., false},
		{"json.Number", map[string]any{"n": json.Number("649")}, "n", 649., false},
		{"missing everything", map[string]any{}, "n", -1., true},
		{"missing not empty", map[string]any{"k": rand.Float64()}, "n", -1, true},
		{"wrong type", map[string]any{"n": "can you really say that it's that different?"}, "n", -1, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := toolware.Request{Args: c.args}
			v, err := req.Float(c.key)
			if c.shouldErr {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, v)
			}
		})
	}
}
