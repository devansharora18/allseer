package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type fileEnvelope struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}

func LoadFromFile(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var envelope fileEnvelope
	if err := unmarshalRules(path, data, &envelope); err != nil {
		return nil, err
	}

	if err := Validate(envelope.Rules); err != nil {
		return nil, err
	}

	return envelope.Rules, nil
}

func unmarshalRules(path string, data []byte, target *fileEnvelope) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, target); err != nil {
			return fmt.Errorf("parse YAML rules file %q: %w", path, err)
		}
	default:
		if err := json.Unmarshal(data, target); err != nil {
			return fmt.Errorf("parse JSON rules file %q: %w", path, err)
		}
	}

	return nil
}

func Validate(rules []Rule) error {
	for i, rule := range rules {
		if strings.TrimSpace(rule.Name) == "" {
			return fmt.Errorf("rules[%d]: name cannot be empty", i)
		}

		switch rule.Action.Type {
		case ActionAllow, ActionBlock, ActionRedirect, ActionModifyRequest, ActionModifyResponse, ActionLog:
		default:
			return fmt.Errorf("rules[%d] (%s): unknown action type %q", i, rule.Name, rule.Action.Type)
		}

		if rule.Action.Type == ActionRedirect && strings.TrimSpace(rule.Action.Target) == "" {
			return fmt.Errorf("rules[%d] (%s): redirect action requires target", i, rule.Name)
		}

		if rule.Match.TimeWindow != nil {
			if rule.Match.TimeWindow.StartHour < 0 || rule.Match.TimeWindow.StartHour > 23 {
				return fmt.Errorf("rules[%d] (%s): start_hour must be within [0,23]", i, rule.Name)
			}

			if rule.Match.TimeWindow.EndHour < 0 || rule.Match.TimeWindow.EndHour > 23 {
				return fmt.Errorf("rules[%d] (%s): end_hour must be within [0,23]", i, rule.Name)
			}
		}
	}

	return nil
}
