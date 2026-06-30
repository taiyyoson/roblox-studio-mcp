// Package mcpkit is a thin, reusable layer over the official MCP Go SDK that
// makes bootstrapping a new server fast: a stderr logger, result helpers, and a
// Registrar pattern for grouping tools.
package mcpkit

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Registrar registers one or more tools onto a server. Define each tool (or
// group of related tools) as a Registrar so the entry point can wire them all
// up with a single call to Build.
type Registrar func(*mcp.Server)

// Build creates a server with the given identity and applies every Registrar.
func Build(name, version string, regs ...Registrar) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil)
	for _, r := range regs {
		r(s)
	}
	return s
}

// Serve runs the server over stdio until the context is cancelled or the client
// disconnects. This is the standard transport for local MCP clients (Claude
// Desktop, Cursor, Claude Code).
func Serve(ctx context.Context, s *mcp.Server) error {
	return s.Run(ctx, &mcp.StdioTransport{})
}
