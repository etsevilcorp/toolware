package configagnostic

import "github.com/mark3labs/mcp-go/mcp"

func ToMCPConfig[I, O any](cfg Config[I, O]) (mcp.Tool, error) {
	return mcp.NewTool(
		cfg.Name,
		mcp.WithDescription(cfg.Description),
		mcp.WithInputSchema[I](),
		mcp.WithOutputSchema[O](),
	), nil
}
