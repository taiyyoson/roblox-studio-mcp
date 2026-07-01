package bridge

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/taiyyoson/roblox-studio-mcp/internal/mcpkit"
)

type Command struct {
	ID   string         `json:"id"`
	Type string         `json:"type"`
	Args map[string]any `json:"args,omitempty"`
}

type Result struct {
	ID     string `json:"id"`
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

type Bridge struct {
	log      *mcpkit.Logger
	timeout  time.Duration
	pollWait time.Duration

	pending chan Command

	mu      sync.Mutex
	waiters map[string]chan Result
	seq     uint64
}

func New(log *mcpkit.Logger) *Bridge {
	return &Bridge{
		log:      log,
		timeout:  30 * time.Second,
		pollWait: 25 * time.Second,
		pending:  make(chan Command),
		waiters:  make(map[string]chan Result),
	}
}

func (b *Bridge) nextID() string {
	return fmt.Sprintf("cmd-%d", atomic.AddUint64(&b.seq, 1))
}

func (b *Bridge) Dispatch(ctx context.Context, typ string, args map[string]any) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	cmd := Command{ID: b.nextID(), Type: typ, Args: args}

	resCh := make(chan Result, 1)
	b.mu.Lock()
	b.waiters[cmd.ID] = resCh
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		delete(b.waiters, cmd.ID)
		b.mu.Unlock()
	}()

	select {
	case b.pending <- cmd:
		b.log.Debug("dispatched", "id", cmd.ID, "type", typ)
	case <-ctx.Done():
		return Result{}, fmt.Errorf("no Studio plugin picked up %q (is the plugin running and connected?): %w", typ, ctx.Err())
	}

	select {
	case res := <-resCh:
		return res, nil
	case <-ctx.Done():
		return Result{}, fmt.Errorf("Studio did not return a result for %q in time: %w", typ, ctx.Err())
	}
}

func (b *Bridge) deliver(res Result) bool {
	b.mu.Lock()
	ch := b.waiters[res.ID]
	b.mu.Unlock()
	if ch == nil {
		return false
	}
	ch <- res
	return true
}
