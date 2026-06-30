# roblox-studio-mcp

mcp server scaffolded from [go-mcp-bootstrap.](https://github.com/taiyyoson/mcp-server-bootstrap) 

(this is my bootstrap)

```sh
go run ./cmd/server   # serve over stdio
```

add tools under internal/tools as mcpkit.Registrars and wire them into cmd/server/main.go.
