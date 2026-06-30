package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/taiyyoson/roblox-studio-mcp/internal/mcpkit"
)

// newTestClient wires the given tools into a server and returns a connected
// in-memory client session for calling them
func newTestClient(t *testing.T, regs ...mcpkit.Registrar) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	srv := mcpkit.Build("test", "0.0.0", regs...)

	ct, st := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

// callText calls a tool and returns the first content block's text plus IsError
func callText(t *testing.T, cs *mcp.ClientSession, name string, args any) (string, bool) {
	t.Helper()
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	if len(res.Content) == 0 {
		t.Fatalf("call %s: empty content", name)
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("call %s: content[0] is %T, want *mcp.TextContent", name, res.Content[0])
	}
	return tc.Text, res.IsError
}
