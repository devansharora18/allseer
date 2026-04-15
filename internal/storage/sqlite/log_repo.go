package sqlite

import (
	"fmt"
	"strings"
	"time"
)

type TrafficLog struct {
	Method     string
	URL        string
	Host       string
	StatusCode int
	Duration   time.Duration
	Action     string
	RuleName   string
	Error      string
	CreatedAt  time.Time
}

type LogRepository struct {
	store *Store
}

func NewLogRepository(store *Store) *LogRepository {
	return &LogRepository{store: store}
}

func (r *LogRepository) Insert(entry TrafficLog) error {
	if r == nil || r.store == nil || r.store.DB() == nil {
		return fmt.Errorf("sqlite log repository is not initialized")
	}

	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}

	entry.Method = strings.TrimSpace(entry.Method)
	if entry.Method == "" {
		entry.Method = "UNKNOWN"
	}

	entry.URL = strings.TrimSpace(entry.URL)
	if entry.URL == "" {
		entry.URL = "unknown"
	}

	entry.Host = strings.TrimSpace(entry.Host)
	if entry.Host == "" {
		entry.Host = "unknown"
	}

	entry.Action = strings.TrimSpace(entry.Action)
	if entry.Action == "" {
		entry.Action = "allow"
	}

	_, err := r.store.DB().Exec(
		`INSERT INTO traffic_logs(created_at, method, url, host, status_code, duration_ms, action, rule_name, error)
         VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.CreatedAt.UTC().Format(time.RFC3339Nano),
		entry.Method,
		entry.URL,
		entry.Host,
		entry.StatusCode,
		entry.Duration.Milliseconds(),
		entry.Action,
		entry.RuleName,
		entry.Error,
	)
	if err != nil {
		return fmt.Errorf("insert traffic log: %w", err)
	}

	return nil
}
