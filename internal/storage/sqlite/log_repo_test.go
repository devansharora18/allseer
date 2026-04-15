package sqlite

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLogRepositoryInsert(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "traffic.db")
	store, err := New(Config{Path: path})
	if err != nil {
		t.Fatalf("failed to initialize sqlite store: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	repo := NewLogRepository(store)
	err = repo.Insert(TrafficLog{
		Method:     "GET",
		URL:        "http://example.com/page",
		Host:       "example.com",
		StatusCode: 200,
		Duration:   1530 * time.Millisecond,
		Action:     "allow",
		RuleName:   "allow-default",
	})
	if err != nil {
		t.Fatalf("failed to insert traffic log: %v", err)
	}

	var (
		method     string
		url        string
		host       string
		statusCode int
		durationMs int64
		action     string
		ruleName   string
	)

	err = store.DB().QueryRow(`
		SELECT method, url, host, status_code, duration_ms, action, rule_name
		FROM traffic_logs
		LIMIT 1
	`).Scan(&method, &url, &host, &statusCode, &durationMs, &action, &ruleName)
	if err != nil {
		t.Fatalf("failed to query inserted traffic log: %v", err)
	}

	if method != "GET" {
		t.Fatalf("unexpected method: got %q", method)
	}

	if url != "http://example.com/page" {
		t.Fatalf("unexpected url: got %q", url)
	}

	if host != "example.com" {
		t.Fatalf("unexpected host: got %q", host)
	}

	if statusCode != 200 {
		t.Fatalf("unexpected status code: got %d", statusCode)
	}

	if durationMs != 1530 {
		t.Fatalf("unexpected duration ms: got %d", durationMs)
	}

	if action != "allow" {
		t.Fatalf("unexpected action: got %q", action)
	}

	if ruleName != "allow-default" {
		t.Fatalf("unexpected rule_name: got %q", ruleName)
	}
}
