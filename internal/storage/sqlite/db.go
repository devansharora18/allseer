package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Config struct {
	Path string
}

type Store struct {
	path string
	db   *sql.DB
}

func New(cfg Config) (*Store, error) {
	path := cfg.Path
	if path == "" {
		path = "allseer.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database %q: %w", path, err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database %q: %w", path, err)
	}

	store := &Store{path: path, db: db}
	if err := store.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *Store) migrate() error {
	const schema = `
CREATE TABLE IF NOT EXISTS traffic_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL,
  method TEXT NOT NULL,
  url TEXT NOT NULL,
  host TEXT NOT NULL,
  status_code INTEGER NOT NULL,
  duration_ms INTEGER NOT NULL,
  action TEXT NOT NULL,
  rule_name TEXT NOT NULL DEFAULT '',
  error TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_traffic_logs_created_at ON traffic_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_traffic_logs_host ON traffic_logs(host);
CREATE INDEX IF NOT EXISTS idx_traffic_logs_action ON traffic_logs(action);
`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("run sqlite migrations for %q: %w", s.path, err)
	}

	return nil
}
