package configagnostic_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/etsevilcorp/toolware/configagnostic"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestToMCPTool(t *testing.T) {
	t.Run("all good", func(t *testing.T) {
		testcfg := configagnostic.Config[testInput, testOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		}

		cfg, err := configagnostic.ToMCPConfig(testcfg)
		require.NoError(t, err)
		require.Equal(t, testcfg.Name, cfg.Name)
		require.Equal(t, testcfg.Description, cfg.Description)

		bi, err := json.Marshal(cfg.InputSchema)
		require.NoError(t, err)
		bo, err := json.Marshal(cfg.OutputSchema)
		require.NoError(t, err)

		ti, err := jsonschema.For[testInput](nil)
		require.NoError(t, err)
		bti, err := json.Marshal(ti)
		require.JSONEq(t, string(bti), string(bi))

		to, err := jsonschema.For[testOutput](nil)
		require.NoError(t, err)
		bto, err := json.Marshal(to)
		require.JSONEq(t, string(bto), string(bo))
	})
	t.Run("bad input", func(t *testing.T) {
		_, err := configagnostic.ToMCPConfig(configagnostic.Config[EviltestInput, testOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		})
		require.Error(t, err)
	})
	t.Run("bad output", func(t *testing.T) {
		_, err := configagnostic.ToMCPConfig(configagnostic.Config[testInput, EviltestOutput]{
			Name:        rand.Text(),
			Description: rand.Text(),
		})
		require.Error(t, err)
	})
}
