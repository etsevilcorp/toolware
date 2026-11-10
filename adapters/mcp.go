package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/etsevilcorp/toolware"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterMCPTools(tb *toolware.Toolbox, server *server.MCPServer, defs map[string]mcp.Tool) error {
	for name, def := range defs {
		h, ok := tb.Get(name)
		if !ok {
			return fmt.Errorf("toolware: tool %q not registered - cannot expose to MCP", name)
		}
		server.AddTool(def, toMCPHandlerFunc(h, name))
	}
	return nil
}

func toMCPHandlerFunc(h toolware.Handler, name string) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := json.Marshal(req.GetArguments())
		if err != nil {
			return nil, fmt.Errorf("toolware: tool %q recieved non-json object", name)
		}

		resp := h(ctx, toolware.Request{
			ToolName: name,
			Args:     args,
		})
		if resp.Err != nil {
			return mcp.NewToolResultError(resp.Err.Error()), nil
		}
		return mcp.NewToolResultStructuredOnly(resp.Result), nil
	}
}
