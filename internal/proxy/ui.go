package proxy

import (
	"encoding/json"
	"net/http"
	"time"
)

// Stats represents firewall statistics
type Stats struct {
	ActiveConnections    int64  `json:"activeConnections"`
	BlockedRequests      int64  `json:"blockedRequests"`
	AllowedRequests      int64  `json:"allowedRequests"`
	TotalBytesForwarded  int64  `json:"totalBytesForwarded"`
	UptimeSeconds        int64  `json:"uptimeSeconds"`
}

// Rule represents a firewall rule
type Rule struct {
	ID        string    `json:"id"`
	Domain    string    `json:"domain,omitempty"`
	Pattern   string    `json:"pattern,omitempty"`
	Type      string    `json:"type"` // DOMAIN, IP, CIDR, etc.
	Action    string    `json:"action"` // ALLOW, BLOCK, REDIRECT
	CreatedAt time.Time `json:"createdAt"`
	Enabled   bool      `json:"enabled"`
}

// handleStatistics returns current firewall stats
func (s *Server) handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := &Stats{
		ActiveConnections:   0, // TODO: track from proxy
		BlockedRequests:     0, // TODO: query from log repo
		AllowedRequests:     0, // TODO: query from log repo
		TotalBytesForwarded: 0, // TODO: track cumulative bytes
		UptimeSeconds:       int64(time.Since(s.startTime).Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleRulesList returns all firewall rules
func (s *Server) handleRulesList(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.getRules(w, r)
	} else if r.Method == http.MethodPost {
		s.createRule(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleRuleDetail handles GET, PUT, DELETE for a specific rule
func (s *Server) handleRuleDetail(w http.ResponseWriter, r *http.Request) {
	ruleID := r.URL.Query().Get("id")
	if ruleID == "" {
		http.Error(w, "Rule ID required", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodDelete {
		s.deleteRule(w, r, ruleID)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getRules(w http.ResponseWriter, r *http.Request) {
	// TODO: Load from rule engine or persistence
	rules := []Rule{}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var rule Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Validate and persist rule
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request, ruleID string) {
	// TODO: Delete rule by ID
	w.WriteHeader(http.StatusNoContent)
}

// handleSPA serves the React SPA
func (s *Server) handleSPA(w http.ResponseWriter, r *http.Request) {
	// For now, serve a simple HTML page
	// TODO: Embed React build files and serve them
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Allseer Firewall</title>
		<script crossorigin src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
		<script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
	</head>
	<body>
		<div id="root">
			<div style="text-align: center; padding: 50px; font-family: Arial, sans-serif;">
				<h1>🔥 Allseer Firewall Dashboard</h1>
				<p>React UI is being loaded... (build with <code>npm run build</code> in web/ folder)</p>
				<p>For now, see API endpoints:</p>
				<ul style="text-align: left; display: inline-block;">
					<li><a href="/api/stats">/api/stats</a> - Firewall statistics</li>
					<li><a href="/api/rules">/api/rules</a> - Firewall rules</li>
					<li><a href="/admin/logs">/admin/logs</a> - Traffic logs</li>
				</ul>
			</div>
		</div>
	</body>
	</html>
	`))
}
