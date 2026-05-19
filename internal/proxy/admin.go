package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"allseer/internal/storage/sqlite"
)

const (
	defaultAdminLogsLimit = 100
	maxAdminLogsLimit     = 1000
)

type adminLogsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Logs  []adminLogEntry `json:"logs"`
}

type adminLogEntry struct {
	CreatedAt  string `json:"created_at"`
	Method     string `json:"method"`
	URL        string `json:"url"`
	Host       string `json:"host"`
	StatusCode int    `json:"status_code"`
	DurationMS int64  `json:"duration_ms"`
	Action     string `json:"action"`
	RuleName   string `json:"rule_name"`
	Error      string `json:"error"`
}

func isLocalAdminLogsRequest(r *http.Request) bool {
	if r == nil || r.URL == nil {
		return false
	}

	if r.URL.IsAbs() {
		return false
	}

	path := strings.TrimSuffix(strings.TrimSpace(r.URL.Path), "/")
	return path == "/admin/logs"
}

func (s *Server) handleAdminLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.logRepo == nil {
		http.Error(w, "log repository is not configured", http.StatusServiceUnavailable)
		return
	}

	limit, err := parseAdminLogsLimit(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logs, err := s.logRepo.ListRecent(limit)
	if err != nil {
		s.logger.Error("failed to list recent traffic logs", "error", err)
		http.Error(w, "failed to list logs", http.StatusInternalServerError)
		return
	}

	response := adminLogsResponse{
		Count: len(logs),
		Limit: limit,
		Logs:  make([]adminLogEntry, 0, len(logs)),
	}

	for _, record := range logs {
		response.Logs = append(response.Logs, toAdminLogEntry(record))
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("failed to encode admin logs response", "error", err)
	}
}

func parseAdminLogsLimit(r *http.Request) (int, error) {
	if r == nil || r.URL == nil {
		return defaultAdminLogsLimit, nil
	}

	raw := strings.TrimSpace(r.URL.Query().Get("limit"))
	if raw == "" {
		return defaultAdminLogsLimit, nil
	}

	limit, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid limit value")
	}

	if limit <= 0 {
		return 0, fmt.Errorf("limit must be positive")
	}

	if limit > maxAdminLogsLimit {
		limit = maxAdminLogsLimit
	}

	return limit, nil
}

func toAdminLogEntry(record sqlite.TrafficLog) adminLogEntry {
	createdAt := ""
	if !record.CreatedAt.IsZero() {
		createdAt = record.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000000Z07:00")
	}

	return adminLogEntry{
		CreatedAt:  createdAt,
		Method:     record.Method,
		URL:        record.URL,
		Host:       record.Host,
		StatusCode: record.StatusCode,
		DurationMS: record.Duration.Milliseconds(),
		Action:     record.Action,
		RuleName:   record.RuleName,
		Error:      record.Error,
	}
}
