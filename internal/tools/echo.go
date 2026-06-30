package tools

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/taiyyoson/roblox-studio-mcp/internal/mcpkit"
)

// echoInput demonstrates a typed, validated argument. The `jsonschema` tag
// becomes the parameter description advertised to the model.
type echoInput struct {
	Message string `json:"message" jsonschema:"the text to echo back"`
	Shout   bool   `json:"shout,omitempty" jsonschema:"if true, return the message in uppercase"`
}

// Echo returns the supplied message, optionally shouted. A minimal template for
// writing your own argument-taking tools.
func Echo() mcpkit.Registrar {
	return func(s *mcp.Server) {
		mcp.AddTool(s, &mcp.Tool{
			Name:        "echo",
			Description: "Echo a message back, optionally in uppercase.",
		}, func(ctx context.Context, req *mcp.CallToolRequest, in echoInput) (*mcp.CallToolResult, any, error) {
			if strings.TrimSpace(in.Message) == "" {
				return mcpkit.Errorf("message must not be empty"), nil, nil
			}
			msg := in.Message
			if in.Shout {
				msg = strings.ToUpper(msg)
			}
			return mcpkit.Text(msg), nil, nil
		})
	}
}
