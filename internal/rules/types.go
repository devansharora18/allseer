package rules

import "time"

type ActionType string

const (
	ActionAllow          ActionType = "allow"
	ActionBlock          ActionType = "block"
	ActionRedirect       ActionType = "redirect"
	ActionModifyRequest  ActionType = "modify_request"
	ActionModifyResponse ActionType = "modify_response"
	ActionLog            ActionType = "log"
)

type Rule struct {
	Name     string  `json:"name" yaml:"name"`
	Enabled  bool    `json:"enabled" yaml:"enabled"`
	Priority int     `json:"priority" yaml:"priority"`
	Match    Matcher `json:"match" yaml:"match"`
	Action   Action  `json:"action" yaml:"action"`
}

type Matcher struct {
	HostPattern  string      `json:"host_pattern" yaml:"host_pattern"`
	HostPatterns []string    `json:"host_patterns,omitempty" yaml:"host_patterns,omitempty"`
	PathPattern  string      `json:"path_pattern" yaml:"path_pattern"`
	Methods      []string    `json:"methods" yaml:"methods"`
	TimeWindow   *TimeWindow `json:"time_window" yaml:"time_window"`
}

type TimeWindow struct {
	StartHour int            `json:"start_hour" yaml:"start_hour"`
	EndHour   int            `json:"end_hour" yaml:"end_hour"`
	Weekdays  []time.Weekday `json:"weekdays" yaml:"weekdays"`
}

type Action struct {
	Type     ActionType        `json:"type" yaml:"type"`
	Target   string            `json:"target,omitempty" yaml:"target,omitempty"`
	Headers  map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Metadata map[string]any    `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type RequestContext struct {
	Method string
	Host   string
	Path   string
	Time   time.Time
}

type Decision struct {
	Matched bool
	Rule    Rule
	Action  Action
}
