package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCreatesDefaultConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	wantPath := filepath.Join(home, ".pchome", "config.toml")
	if path != wantPath {
		t.Fatalf("expected path %q, got %q", wantPath, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config file: %v", err)
	}
	if string(data) != DefaultTOML() {
		t.Fatalf("default config file content mismatch")
	}

	if cfg.Output.Format != "text" {
		t.Fatalf("expected default format text, got %q", cfg.Output.Format)
	}
	if cfg.Output.NameWidth != 30 {
		t.Fatalf("expected default name width 30, got %d", cfg.Output.NameWidth)
	}
	if cfg.I18N.Language != "zh-TW" {
		t.Fatalf("expected default i18n.language zh-TW, got %q", cfg.I18N.Language)
	}
	if cfg.Search.Limit != 10 {
		t.Fatalf("expected default search limit 10, got %d", cfg.Search.Limit)
	}
	if !cfg.Search.ShowURL {
		t.Fatalf("expected default search.show_url to be true")
	}
}

func TestLoadHonorsOverridesAndDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".pchome", "config.toml")
	if err := ensureExistsAt(path); err != nil {
		t.Fatalf("ensureExistsAt: %v", err)
	}

	if err := writeFixtureConfig(path, "override.toml"); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, _, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.HTTP.Timeout != "45s" {
		t.Fatalf("expected timeout override, got %q", cfg.HTTP.Timeout)
	}
	if cfg.Output.Format != "json" {
		t.Fatalf("expected format override, got %q", cfg.Output.Format)
	}
	if cfg.I18N.Language != "en" {
		t.Fatalf("expected i18n.language override, got %q", cfg.I18N.Language)
	}
	if cfg.Search.Limit != 25 {
		t.Fatalf("expected search.limit override, got %d", cfg.Search.Limit)
	}
	if strings.Join(cfg.Search.Columns, ",") != "#,name,price,url" {
		t.Fatalf("unexpected search.columns: %v", cfg.Search.Columns)
	}
	if cfg.Search.ShowURL != true {
		t.Fatalf("expected search.show_url to preserve default true")
	}
	if cfg.Recommend.Top != 7 {
		t.Fatalf("expected recommend.top override, got %d", cfg.Recommend.Top)
	}
	if !cfg.Recommend.ShowWhy {
		t.Fatalf("expected recommend.show_why override true")
	}
	if cfg.Compare.ShowURL != true {
		t.Fatalf("expected compare.show_url to preserve default true")
	}
	if cfg.Suggest.Limit != 4 {
		t.Fatalf("expected suggest.limit override, got %d", cfg.Suggest.Limit)
	}
	if cfg.Hermes.Token != "abc123" {
		t.Fatalf("expected hermes token override, got %q", cfg.Hermes.Token)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".pchome", "config.toml")
	if err := ensureExistsAt(path); err != nil {
		t.Fatalf("ensureExistsAt: %v", err)
	}

	if err := writeFixtureConfig(path, "invalid_unknown_key.toml"); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, _, err := Load()
	if err == nil {
		t.Fatalf("expected Load to reject unknown keys")
	}
	if !strings.Contains(err.Error(), "unknown config keys") {
		t.Fatalf("expected unknown key error, got %v", err)
	}
}

func TestLoadRejectsInvalidLanguage(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".pchome", "config.toml")
	if err := ensureExistsAt(path); err != nil {
		t.Fatalf("ensureExistsAt: %v", err)
	}

	if err := writeFixtureConfig(path, "invalid_language.toml"); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, _, err := Load()
	if err == nil {
		t.Fatalf("expected Load to reject invalid i18n.language")
	}
	if !strings.Contains(err.Error(), "invalid i18n.language") {
		t.Fatalf("expected invalid i18n.language error, got %v", err)
	}
}
