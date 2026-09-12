package adapters

import (
	"context"
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
		raw := req.GetRawArguments()
		args, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("toolware: tool %q recieved non-map[string]any object(%T)", name, raw)
		}

		resp := h(ctx, toolware.Request{
			ToolName: name,
			Args:     args,
		})
		if resp.Err != nil {
			return nil, resp.Err
		}
		return mcp.NewToolResultText(fmt.Sprint(resp.Result)), nil
	}
}
