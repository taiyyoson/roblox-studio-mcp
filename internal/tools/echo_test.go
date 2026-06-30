package tools

import "testing"

func TestEcho(t *testing.T) {
	cs := newTestClient(t, Echo())

	t.Run("plain", func(t *testing.T) {
		out, isErr := callText(t, cs, "echo", map[string]any{"message": "hi"})
		if isErr {
			t.Fatalf("error result: %q", out)
		}
		if out != "hi" {
			t.Errorf("echo = %q, want %q", out, "hi")
		}
	})

	t.Run("shout", func(t *testing.T) {
		out, _ := callText(t, cs, "echo", map[string]any{"message": "hi", "shout": true})
		if out != "HI" {
			t.Errorf("echo shout = %q, want %q", out, "HI")
		}
	})

	t.Run("blank is error", func(t *testing.T) {
		_, isErr := callText(t, cs, "echo", map[string]any{"message": "   "})
		if !isErr {
			t.Error("blank message should return an error result")
		}
	})
}
