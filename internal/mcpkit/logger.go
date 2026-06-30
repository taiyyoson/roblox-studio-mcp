package mcpkit

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// serializes stderr writes so concurrent log lines don't interleave
var writeMu sync.Mutex

// Logger is a tiny leveled logger.
//
// IMPORTANT: it writes to stderr only. In an stdio MCP server, stdout is
// reserved exclusively for the JSON-RPC message stream — writing logs there
// corrupts the protocol and breaks the client connection.
type Logger struct {
	scope string
	level Level
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func parseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// NewLogger returns a logger whose threshold is read from the LOG_LEVEL env var
// (debug|info|warn|error), defaulting to info.
func NewLogger(scope string) *Logger {
	return &Logger{scope: scope, level: parseLevel(os.Getenv("LOG_LEVEL"))}
}

// Child returns a logger with a nested scope (e.g. "mcp" -> "mcp:ping").
func (l *Logger) Child(scope string) *Logger {
	return &Logger{scope: l.scope + ":" + scope, level: l.level}
}

func (l *Logger) Debug(msg string, kv ...any) { l.log(LevelDebug, "debug", msg, kv) }
func (l *Logger) Info(msg string, kv ...any)  { l.log(LevelInfo, "info", msg, kv) }
func (l *Logger) Warn(msg string, kv ...any)  { l.log(LevelWarn, "warn", msg, kv) }
func (l *Logger) Error(msg string, kv ...any) { l.log(LevelError, "error", msg, kv) }

func (l *Logger) log(lvl Level, name, msg string, kv []any) {
	if lvl < l.level {
		return
	}
	ts := time.Now().UTC().Format(time.RFC3339)
	extra := ""
	if len(kv) > 0 {
		if b, err := json.Marshal(kvToMap(kv)); err == nil {
			extra = " " + string(b)
		}
	}
	writeMu.Lock()
	fmt.Fprintf(os.Stderr, "%s [%s] [%s] %s%s\n", ts, name, l.scope, msg, extra)
	writeMu.Unlock()
}

// kvToMap turns alternating key/value args into a map for structured logging.
func kvToMap(kv []any) map[string]any {
	m := make(map[string]any, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[fmt.Sprint(kv[i])] = kv[i+1]
	}
	return m
}
