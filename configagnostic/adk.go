package configagnostic

import (
	"github.com/google/jsonschema-go/jsonschema"
	"google.golang.org/adk/v2/tool/functiontool"
)

func ToADKConfig[I, O any](cfg Config[I, O]) (functiontool.Config, error) {
	iSchema, err := jsonschema.For[I](&jsonschema.ForOptions{IgnoreInvalidTypes: true})
	if err != nil {
		return functiontool.Config{}, err
	}

	oSchema, err := jsonschema.For[O](&jsonschema.ForOptions{IgnoreInvalidTypes: true})
	if err != nil {
		return functiontool.Config{}, err
	}

	return functiontool.Config{
		Name:         cfg.Name,
		Description:  cfg.Description,
		InputSchema:  iSchema,
		OutputSchema: oSchema,
	}, nil
}
