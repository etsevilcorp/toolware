package configagnostic_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/etsevilcorp/toolware/configagnostic"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/require"
)

type TestInput struct {
	Offset int `json:"offset"`
}
type TestOutput struct {
	Text string `json:"text"`
}

type Node struct {
	Value    string  `json:"value"`
	Children []*Node `json:"children,omitempty"`
}
type EvilTestInput struct {
	Node Node `json:"evilNode"`
}
type EvilTestOutput struct {
	Node Node `json:"evilNode"`
}

func TestToADKConfig(t *testing.T) {
	t.Run("all good", func(t *testing.T) {
		testcfg := configagnostic.Config[TestInput, TestOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		}

		cfg, err := configagnostic.ToADKConfig(testcfg)
		require.NoError(t, err)
		bi, err := json.Marshal(cfg.InputSchema)
		require.NoError(t, err)
		bo, err := json.Marshal(cfg.OutputSchema)
		require.NoError(t, err)

		ti, err := jsonschema.For[TestInput](nil)
		require.NoError(t, err)
		bti, err := json.Marshal(ti)
		require.Equal(t, bti, bi)

		to, err := jsonschema.For[TestOutput](nil)
		require.NoError(t, err)
		bto, err := json.Marshal(to)
		require.Equal(t, bto, bo)
	})
	t.Run("bad input", func(t *testing.T) {
		_, err := configagnostic.ToADKConfig(configagnostic.Config[EvilTestInput, TestOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		})
		require.Error(t, err)
	})
	t.Run("bad output", func(t *testing.T) {
		_, err := configagnostic.ToADKConfig(configagnostic.Config[TestInput, EvilTestOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		})
		require.Error(t, err)
	})
}
