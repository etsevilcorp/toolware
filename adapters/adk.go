package adapters

import (
	"fmt"

	"github.com/etsevilcorp/toolware"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
)

func RegisterADKTools(tb *toolware.Toolbox, defs []functiontool.Config) ([]tool.Tool, error) {
	tools := make([]tool.Tool, 0, len(defs))
	for _, def := range defs {
		h, ok := tb.Get(def.Name)
		if !ok {
			return nil, fmt.Errorf("toolware: tool %q not registered - cannot expose to ADK", def.Name)
		}
		tool, err := functiontool.New(def, toADKHandlerFunc(h, def.Name))
		if err != nil {
			return nil, fmt.Errorf("toolware: tool %q couldn't be registered in adk: %w", def.Name, err)
		}
		tools = append(tools, tool)
	}

	return tools, nil
}

func toADKHandlerFunc(h toolware.Handler, name string) func(ctx agent.Context, req toolware.Request) (toolware.Response, error) {
	return func(ctx agent.Context, req toolware.Request) (toolware.Response, error) {
		resp := h(ctx, req)
		if resp.Err != nil {
			return toolware.Response{}, resp.Err
		}

		return resp, nil
	}
}
