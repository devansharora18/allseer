package proxy

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"allseer/internal/rules"
	"allseer/internal/storage/sqlite"
)

func TestHandleAdminLogsSuccess(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "logs.db")
	store, err := sqlite.New(sqlite.Config{Path: path})
	if err != nil {
		t.Fatalf("failed to initialize sqlite store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	repo := sqlite.NewLogRepository(store)
	if err := repo.Insert(sqlite.TrafficLog{
		Method:     http.MethodGet,
		URL:        "http://example.com/page",
		Host:       "example.com",
		StatusCode: http.StatusOK,
		Duration:   42 * time.Millisecond,
		Action:     "allow",
		RuleName:   "allow-default",
		CreatedAt:  time.Now().UTC(),
	}); err != nil {
		t.Fatalf("failed to insert seed log: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(Options{
		ListenAddr:    "127.0.0.1:0",
		RuleEngine:    rules.NewEngine(nil),
		LogRepository: repo,
	}, logger)

	req := httptest.NewRequest(http.MethodGet, "/admin/logs?limit=1", nil)
	rr := httptest.NewRecorder()

	server.handleProxyRequest(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected json content type, got %q", got)
	}

	var payload struct {
		Count int `json:"count"`
		Limit int `json:"limit"`
		Logs  []struct {
			URL    string `json:"url"`
			Host   string `json:"host"`
			Action string `json:"action"`
		} `json:"logs"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response payload: %v", err)
	}

	if payload.Count != 1 {
		t.Fatalf("expected count=1, got %d", payload.Count)
	}

	if payload.Limit != 1 {
		t.Fatalf("expected limit=1, got %d", payload.Limit)
	}

	if len(payload.Logs) != 1 {
		t.Fatalf("expected one log entry, got %d", len(payload.Logs))
	}

	if payload.Logs[0].URL != "http://example.com/page" {
		t.Fatalf("unexpected log URL %q", payload.Logs[0].URL)
	}
}

func TestHandleAdminLogsInvalidLimit(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "logs.db")
	store, err := sqlite.New(sqlite.Config{Path: path})
	if err != nil {
		t.Fatalf("failed to initialize sqlite store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(Options{
		ListenAddr:    "127.0.0.1:0",
		RuleEngine:    rules.NewEngine(nil),
		LogRepository: sqlite.NewLogRepository(store),
	}, logger)

	req := httptest.NewRequest(http.MethodGet, "/admin/logs?limit=abc", nil)
	rr := httptest.NewRecorder()

	server.handleProxyRequest(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleAdminLogsMethodNotAllowed(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "logs.db")
	store, err := sqlite.New(sqlite.Config{Path: path})
	if err != nil {
		t.Fatalf("failed to initialize sqlite store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(Options{
		ListenAddr:    "127.0.0.1:0",
		RuleEngine:    rules.NewEngine(nil),
		LogRepository: sqlite.NewLogRepository(store),
	}, logger)

	req := httptest.NewRequest(http.MethodPost, "/admin/logs", nil)
	rr := httptest.NewRecorder()

	server.handleProxyRequest(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	if allow := rr.Header().Get("Allow"); allow != http.MethodGet {
		t.Fatalf("expected Allow header to be %q, got %q", http.MethodGet, allow)
	}
}

func TestHandleAdminLogsRepositoryUnavailable(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := NewServer(Options{ListenAddr: "127.0.0.1:0", RuleEngine: rules.NewEngine(nil)}, logger)

	req := httptest.NewRequest(http.MethodGet, "/admin/logs", nil)
	rr := httptest.NewRecorder()

	server.handleProxyRequest(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rr.Code)
	}
}

func TestIsLocalAdminLogsRequestAbsoluteURL(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://remote.example/admin/logs", nil)
	if isLocalAdminLogsRequest(req) {
		t.Fatalf("expected absolute URL request not to be treated as local admin request")
	}
}
