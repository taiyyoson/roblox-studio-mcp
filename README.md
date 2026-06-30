# roblox-studio-mcp

mcp server scaffolded from go-mcp-bootstrap.

```sh
go run ./cmd/server   # serve over stdio
```

add tools under internal/tools as mcpkit.Registrars and wire them into cmd/server/main.go.
