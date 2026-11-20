package configagnostic

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

func ToMCPConfig[I, O any](cfg Config[I, O]) (mcp.Tool, error) {
	iSchema, err := jsonschema.For[I](&jsonschema.ForOptions{IgnoreInvalidTypes: true})
	if err != nil {
		return mcp.Tool{}, err
	}
	iBytes, err := json.Marshal(iSchema)
	if err != nil {
		return mcp.Tool{}, err
	}

	oSchema, err := jsonschema.For[O](&jsonschema.ForOptions{IgnoreInvalidTypes: true})
	if err != nil {
		return mcp.Tool{}, err
	}
	oBytes, err := json.Marshal(oSchema)
	if err != nil {
		return mcp.Tool{}, err
	}

	return mcp.NewTool(
		cfg.Name,
		mcp.WithDescription(cfg.Description),
		mcp.WithRawInputSchema(iBytes),
		mcp.WithRawOutputSchema(oBytes),
	), nil
}
