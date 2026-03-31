package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/i18n"
)

const (
	dirName       = ".pchome"
	fileName      = "config.toml"
	currentFormat = 1
)

type Config struct {
	Version   int             `toml:"version"`
	HTTP      HTTPConfig      `toml:"http"`
	Output    OutputConfig    `toml:"output"`
	I18N      I18NConfig      `toml:"i18n"`
	Search    SearchConfig    `toml:"search"`
	Recommend RecommendConfig `toml:"recommend"`
	Compare   CompareConfig   `toml:"compare"`
	Suggest   SuggestConfig   `toml:"suggest"`
	Hermes    HermesConfig    `toml:"hermes"`
}

type HTTPConfig struct {
	Timeout string `toml:"timeout"`
}

type OutputConfig struct {
	Format        string `toml:"format"`
	SchemaVersion string `toml:"schema_version"`
	NameWidth     int    `toml:"name_width"`
}

type I18NConfig struct {
	Language string `toml:"language"`
}

type SearchConfig struct {
	Category   string   `toml:"category"`
	Brand      string   `toml:"brand"`
	Sort       string   `toml:"sort"`
	PageSize   int      `toml:"page_size"`
	Limit      int      `toml:"limit"`
	MinPrice   int      `toml:"min_price"`
	MaxPrice   int      `toml:"max_price"`
	MinRating  float64  `toml:"min_rating"`
	InStock    bool     `toml:"in_stock"`
	Arrival24h bool     `toml:"arrival_24h"`
	ShowURL    bool     `toml:"show_url"`
	Compact    bool     `toml:"compact"`
	Wide       bool     `toml:"wide"`
	Columns    []string `toml:"columns"`
}

type RecommendConfig struct {
	Top     int      `toml:"top"`
	ShowURL bool     `toml:"show_url"`
	ShowWhy bool     `toml:"show_why"`
	Compact bool     `toml:"compact"`
	Wide    bool     `toml:"wide"`
	Columns []string `toml:"columns"`
}

type CompareConfig struct {
	ShowURL bool     `toml:"show_url"`
	Compact bool     `toml:"compact"`
	Wide    bool     `toml:"wide"`
	Columns []string `toml:"columns"`
}

type SuggestConfig struct {
	Limit int `toml:"limit"`
}

type HermesConfig struct {
	Token string `toml:"token"`
}

func Default() Config {
	return Config{
		Version: currentFormat,
		HTTP: HTTPConfig{
			Timeout: "20s",
		},
		Output: OutputConfig{
			Format:        "text",
			SchemaVersion: catalog.SchemaVersion,
			NameWidth:     30,
		},
		I18N: I18NConfig{
			Language: string(i18n.LangZHTW),
		},
		Search: SearchConfig{
			Category:   "",
			Brand:      "",
			Sort:       "relevance",
			PageSize:   10,
			Limit:      10,
			MinPrice:   0,
			MaxPrice:   0,
			MinRating:  0,
			InStock:    false,
			Arrival24h: false,
			ShowURL:    true,
			Compact:    false,
			Wide:       false,
			Columns:    nil,
		},
		Recommend: RecommendConfig{
			Top:     10,
			ShowURL: true,
			ShowWhy: false,
			Compact: false,
			Wide:    false,
			Columns: nil,
		},
		Compare: CompareConfig{
			ShowURL: true,
			Compact: false,
			Wide:    false,
			Columns: nil,
		},
		Suggest: SuggestConfig{
			Limit: 10,
		},
		Hermes: HermesConfig{
			Token: "",
		},
	}
}

func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, dirName, fileName), nil
}

func EnsureExists() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	if err := ensureExistsAt(path); err != nil {
		return "", err
	}
	return path, nil
}

func Load() (Config, string, error) {
	cfg := Default()

	path, err := EnsureExists()
	if err != nil {
		return Config{}, "", err
	}

	meta, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return Config{}, "", fmt.Errorf("load config %s: %w", path, err)
	}
	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, 0, len(undecoded))
		for _, key := range undecoded {
			keys = append(keys, key.String())
		}
		return Config{}, "", fmt.Errorf("unknown config keys: %s", strings.Join(keys, ", "))
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, "", err
	}

	return cfg, path, nil
}

func (c Config) Validate() error {
	if c.Version != currentFormat {
		return fmt.Errorf("unsupported config version %d (expected %d)", c.Version, currentFormat)
	}

	if _, err := c.TimeoutDuration(); err != nil {
		return fmt.Errorf("invalid http.timeout %q: %w", c.HTTP.Timeout, err)
	}

	switch c.Output.Format {
	case "text", "json", "ndjson":
	default:
		return fmt.Errorf("invalid output.format %q", c.Output.Format)
	}

	if c.Output.SchemaVersion != catalog.SchemaVersion {
		return fmt.Errorf("invalid output.schema_version %q", c.Output.SchemaVersion)
	}
	if c.Output.NameWidth <= 0 {
		return fmt.Errorf("output.name_width must be greater than 0")
	}
	if _, err := i18n.ParseLanguage(c.I18N.Language); err != nil {
		return fmt.Errorf("invalid i18n.language %q: %w", c.I18N.Language, err)
	}

	if err := validateSort(c.Search.Sort); err != nil {
		return err
	}
	if c.Search.PageSize <= 0 {
		return fmt.Errorf("search.page_size must be greater than 0")
	}
	if c.Search.Limit <= 0 {
		return fmt.Errorf("search.limit must be greater than 0")
	}
	if c.Search.MinPrice < 0 {
		return fmt.Errorf("search.min_price must be >= 0")
	}
	if c.Search.MaxPrice < 0 {
		return fmt.Errorf("search.max_price must be >= 0")
	}
	if c.Search.MaxPrice > 0 && c.Search.MinPrice > c.Search.MaxPrice {
		return fmt.Errorf("search.min_price cannot exceed search.max_price")
	}
	if c.Search.MinRating < 0 {
		return fmt.Errorf("search.min_rating must be >= 0")
	}
	if c.Search.Compact && c.Search.Wide {
		return fmt.Errorf("search.compact and search.wide cannot both be true")
	}
	if err := validateColumns("search.columns", c.Search.Columns, []string{"#", "price", "list", "discount", "rating", "reviews", "24h", "qty", "stock", "brand", "name", "id", "url", "desc"}); err != nil {
		return err
	}

	if c.Recommend.Top <= 0 {
		return fmt.Errorf("recommend.top must be greater than 0")
	}
	if c.Recommend.Compact && c.Recommend.Wide {
		return fmt.Errorf("recommend.compact and recommend.wide cannot both be true")
	}
	if err := validateColumns("recommend.columns", c.Recommend.Columns, []string{"#", "score", "price", "list", "discount", "rating", "reviews", "24h", "qty", "stock", "brand", "name", "id", "url", "why"}); err != nil {
		return err
	}

	if c.Compare.Compact && c.Compare.Wide {
		return fmt.Errorf("compare.compact and compare.wide cannot both be true")
	}
	if err := validateColumns("compare.columns", c.Compare.Columns, []string{"#", "price", "list", "discount", "rating", "reviews", "24h", "qty", "stock", "brand", "name", "id", "url"}); err != nil {
		return err
	}

	if c.Suggest.Limit <= 0 {
		return fmt.Errorf("suggest.limit must be greater than 0")
	}

	return nil
}

func (c Config) TimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(strings.TrimSpace(c.HTTP.Timeout))
}

func DefaultTOML() string {
	return `version = 1

[http]
timeout = "20s"

[output]
format = "text"
schema_version = "v1"
name_width = 30

[i18n]
language = "zh-TW"

[search]
category = ""
brand = ""
sort = "relevance"
page_size = 10
limit = 10
min_price = 0
max_price = 0
min_rating = 0
in_stock = false
arrival_24h = false
show_url = true
compact = false
wide = false
columns = []

[recommend]
top = 10
show_url = true
show_why = false
compact = false
wide = false
columns = []

[compare]
show_url = true
compact = false
wide = false
columns = []

[suggest]
limit = 10

[hermes]
token = ""
`
}

func ensureExistsAt(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config directory %s: %w", dir, err)
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config file %s: %w", path, err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("create config file %s: %w", path, err)
	}
	defer file.Close()

	if _, err := file.WriteString(DefaultTOML()); err != nil {
		return fmt.Errorf("write default config %s: %w", path, err)
	}
	return nil
}

func validateSort(sort string) error {
	switch strings.TrimSpace(sort) {
	case "relevance", "price-asc", "price-desc", "newest", "best-selling":
		return nil
	default:
		return fmt.Errorf("invalid search.sort %q", sort)
	}
}

func validateColumns(name string, columns []string, allowed []string) error {
	if len(columns) == 0 {
		return nil
	}

	allowedSet := map[string]bool{}
	for _, key := range allowed {
		allowedSet[key] = true
	}
	for _, key := range columns {
		if !allowedSet[strings.ToLower(strings.TrimSpace(key))] {
			return fmt.Errorf("invalid %s entry %q", name, key)
		}
	}
	return nil
}
