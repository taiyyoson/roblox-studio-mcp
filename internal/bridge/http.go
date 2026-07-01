package bridge

import (
	"encoding/json"
	"net/http"
	"time"
)

func (b *Bridge) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/poll", b.handlePoll)
	mux.HandleFunc("/result", b.handleResult)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func (b *Bridge) handlePoll(w http.ResponseWriter, r *http.Request) {
	select {
	case cmd := <-b.pending:
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cmd)
	case <-time.After(b.pollWait):
		w.WriteHeader(http.StatusNoContent)
	case <-r.Context().Done():
	}
}

func (b *Bridge) handleResult(w http.ResponseWriter, r *http.Request) {
	var res Result
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if !b.deliver(res) {
		b.log.Warn("result for unknown or expired command", "id", res.ID)
	}
	w.WriteHeader(http.StatusOK)
}
