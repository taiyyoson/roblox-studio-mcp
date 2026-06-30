// Package tools holds the concrete MCP tools. Each tool is exposed as an
// mcpkit.Registrar so it can be wired into the server from the entry point.
package tools

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/taiyyoson/roblox-studio-mcp/internal/mcpkit"
)

var start = time.Now()

// pingInput has no fields — ping takes no arguments.
type pingInput struct{}

// Ping is a health-check tool: returns "pong" plus the server uptime.
func Ping() mcpkit.Registrar {
	return func(s *mcp.Server) {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "ping",
			Description: "Health check. Returns 'pong' and the server uptime in seconds.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, in pingInput) (*mcp.CallToolResult, any, error) {
			uptime := int(time.Since(start).Seconds())
			return mcpkit.Textf("pong (uptime %ds)", uptime), nil, nil
		})
	}
}
