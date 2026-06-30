package tools

import (
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	cs := newTestClient(t, Ping())
	out, isErr := callText(t, cs, "ping", struct{}{})
	if isErr {
		t.Fatalf("ping returned an error result: %q", out)
	}
	if !strings.HasPrefix(out, "pong") {
		t.Errorf("ping = %q, want prefix %q", out, "pong")
	}
}
