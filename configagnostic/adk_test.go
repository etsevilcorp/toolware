package configagnostic_test

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/etsevilcorp/toolware/configagnostic"
	"github.com/stretchr/testify/require"
)

type TestInput struct {
	Offset int `json:"offset"`
}
type TestOutput struct {
	Text string `json:"text"`
}

func TestToADKConfig(t *testing.T) {

	cfg, err := configagnostic.ToADKConfig(configagnostic.Config[TestInput, TestOutput]{
		Name:        rand.Text(),
		Description: rand.Text(),
	})
	require.NoError(t, err)

	b, _ := json.Marshal(cfg.InputSchema)
	t.Logf("%s", b)
}
