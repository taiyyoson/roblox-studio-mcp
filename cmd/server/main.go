// Command server is the entry point for the MCP server. It wires up the tool
// registrars and serves them over stdio.
//
// Add a new tool: implement it under internal/tools as an mcpkit.Registrar,
// then add it to the Build call below.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/taiyyoson/roblox-studio-mcp/internal/mcpkit"
	"github.com/taiyyoson/roblox-studio-mcp/internal/tools"
)

const (
	serverName    = "github.com/taiyyoson/roblox-studio-mcp"
	serverVersion = "0.1.0"
)

func main() {
	log := mcpkit.NewLogger("mcp")

	// cancel cleanly on SIGINT/SIGTERM so the stdio loop unwinds
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := mcpkit.Build(serverName, serverVersion,
		tools.Ping(),
		tools.Echo(),
	)

	log.Info("starting server over stdio", "name", serverName, "version", serverVersion)
	if err := mcpkit.Serve(ctx, srv); err != nil {
		log.Error("server exited with error", "err", err.Error())
		os.Exit(1)
	}
	log.Info("server stopped")
}
