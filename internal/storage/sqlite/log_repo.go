package sqlite

import "time"

type TrafficLog struct {
	Method     string
	URL        string
	StatusCode int
	Duration   time.Duration
	Action     string
	CreatedAt  time.Time
}

type LogRepository struct {
	store *Store
}

func NewLogRepository(store *Store) *LogRepository {
	return &LogRepository{store: store}
}

func (r *LogRepository) Insert(_ TrafficLog) error {
	return nil
}
