package mcpkit

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Text wraps a plain string in a tool result.
func Text(s string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: s}},
	}
}

// Textf is Text with fmt-style formatting.
func Textf(format string, args ...any) *mcp.CallToolResult {
	return Text(fmt.Sprintf(format, args...))
}

// JSON marshals v to indented JSON and returns it as a text result, for
// returning structured data (maps, slices, structs) to the model
func JSON(v any) *mcp.CallToolResult {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return Errorf("failed to encode result: %v", err)
	}
	return Text(string(b))
}

// Errorf returns a tool-level error result (IsError set). Use this for expected
// failures the model should see and react to; reserve returning a Go error for
// unexpected internal faults.
func Errorf(format string, args ...any) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(format, args...)}},
	}
}
