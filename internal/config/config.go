package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Duration time.Duration

func (d Duration) Value() time.Duration {
	return time.Duration(d)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		parsed, err := time.ParseDuration(text)
		if err != nil {
			return fmt.Errorf("invalid duration %q: %w", text, err)
		}

		*d = Duration(parsed)
		return nil
	}

	var ns int64
	if err := json.Unmarshal(data, &ns); err == nil {
		*d = Duration(time.Duration(ns))
		return nil
	}

	return fmt.Errorf("duration must be string (example: \"10s\") or integer nanoseconds")
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

type Config struct {
	Proxy    ProxyConfig    `json:"proxy"`
	Rules    RulesConfig    `json:"rules"`
	Database DatabaseConfig `json:"database"`
}

type ProxyConfig struct {
	ListenAddr        string   `json:"listen_addr"`
	ReadHeaderTimeout Duration `json:"read_header_timeout"`
	IdleTimeout       Duration `json:"idle_timeout"`
	ShutdownTimeout   Duration `json:"shutdown_timeout"`
}

type RulesConfig struct {
	File        string `json:"file"`
	AdBlockFile string `json:"ad_block_file"`
}

type DatabaseConfig struct {
	Path string `json:"path"`
}

func Default() Config {
	return Config{
		Proxy: ProxyConfig{
			ListenAddr:        "127.0.0.1:8080",
			ReadHeaderTimeout: Duration(10 * time.Second),
			IdleTimeout:       Duration(60 * time.Second),
			ShutdownTimeout:   Duration(10 * time.Second),
		},
		Rules: RulesConfig{
			File:        "config/rules.example.yaml",
			AdBlockFile: "config/ad_domains.txt",
		},
		Database: DatabaseConfig{
			Path: "allseer.db",
		},
	}
}

func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c Config) Validate() error {
	if c.Proxy.ListenAddr == "" {
		return fmt.Errorf("proxy.listen_addr cannot be empty")
	}

	if c.Proxy.ReadHeaderTimeout.Value() <= 0 {
		return fmt.Errorf("proxy.read_header_timeout must be positive")
	}

	if c.Proxy.IdleTimeout.Value() <= 0 {
		return fmt.Errorf("proxy.idle_timeout must be positive")
	}

	if c.Proxy.ShutdownTimeout.Value() <= 0 {
		return fmt.Errorf("proxy.shutdown_timeout must be positive")
	}

	if c.Database.Path == "" {
		return fmt.Errorf("database.path cannot be empty")
	}

	return nil
}
