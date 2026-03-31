<p align="center">
  <img src="./banner.png" alt="pchome-cli" />
</p>

# pchome-cli

`pchome` is a Go CLI for searching PChome 24h products, viewing product details, comparing products, and getting recommendations.

The CLI is optimized for two modes:

- Human-friendly text output for browsing and decision-making.
- Agent-friendly `json` and `ndjson` output with a normalized schema (`v1`).

The default human-facing language is Traditional Chinese (Taiwan). To switch the app to English, set `i18n.language = "en"` in `~/.pchome/config.toml`.

## Installation

```bash
go install github.com/oliy/pchome-cli/cmd/pchome@latest
```

## Quickstart

```bash
# Search
go run ./cmd/pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# View a product
go run ./cmd/pchome view DRAA5K-A900JOK9O

# Recommendations
go run ./cmd/pchome recommend DMBL53-A900JDNJS --top 8 --why

# Compare products
go run ./cmd/pchome compare DRAA5K-A900JOK9O DMBL53-A900JDNJS

# Suggestions
go run ./cmd/pchome suggest "掃地機"
```

## Command Model

The command surface is intentionally task-oriented:

- `search QUERY`
- `view PRODUCT`
- `recommend PRODUCT`
- `compare PRODUCT [PRODUCT...]`
- `suggest QUERY`

`PRODUCT` can be:

- A raw product ID like `DRAA5K-A900JOK9O`
- A suffixed product ID like `DRAA5K-A900JOK9O-000`
- A full PChome product URL like `https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O`

## Repository Layout

The source tree now follows a more typical open-source Go CLI structure:

- `cmd/`: Cobra command layer and text rendering
- `cmd/pchome/`: binary entrypoint
- `pkg/catalog/`: normalized product models and aggregate service
- `pkg/config/`: config loading and validation
- `pkg/i18n/`: language handling and translation catalog
- `pkg/output/`: shared table rendering
- `pkg/pchome/`: upstream PChome API clients
- `cmd/testdata/`: golden fixtures for CLI help and text output
- `pkg/*/testdata/`: package-level test fixtures

## Config

On startup, `pchome` ensures a config file exists at:

```bash
~/.pchome/config.toml
```

If the file is missing, it is created automatically with the full default configuration.

Config precedence is:

- CLI flags
- `~/.pchome/config.toml`
- Built-in defaults

Example:

```toml
version = 1

[http]
timeout = "20s"

[output]
format = "text"
schema_version = "v1"
name_width = 30

[i18n]
language = "zh-TW"

[search]
sort = "relevance"
page_size = 10
limit = 10
show_url = true
columns = []

[recommend]
top = 10
show_url = true
show_why = false
columns = []

[compare]
show_url = true
columns = []

[suggest]
limit = 10

[hermes]
token = ""
```

Notes:

- `columns = []` means “use the built-in default column order for that command”.
- If you set `columns`, that list becomes the default column order for the command.
- `i18n.language` currently supports `zh-TW` and `en`.
- The config loader rejects unknown keys so typos do not silently get ignored.

## Output Modes

### Text

Human-oriented output with compact tables for list commands and a structured detail view for `view`.

### JSON

Normalized machine-readable output:

```bash
go run ./cmd/pchome search "掃地機器人" --limit 3 --format json
go run ./cmd/pchome view DRAA5K-A900JOK9O --format json
```

### NDJSON

Line-delimited items for list-oriented commands:

```bash
go run ./cmd/pchome search "掃地機器人" --limit 5 --format ndjson
go run ./cmd/pchome recommend DMBL53-A900JDNJS --top 5 --format ndjson
```

`ndjson` is supported on `search`, `recommend`, `compare`, and `suggest`.

## Search Examples

```bash
# Brand and rating filter
go run ./cmd/pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 24h only, sorted by price
go run ./cmd/pchome search "掃地機器人" --arrival-24h --sort price-asc

# Custom text columns
go run ./cmd/pchome search "掃地機器人" \
  --columns "#,name,price,rating,reviews,24h,brand,qty,url"
```

## Recommendation Example

```bash
# Explain why each item was recommended
go run ./cmd/pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

## Notes

- Recommendation token precedence is `hermes.token` -> `PCHOME_HERMES_TOKEN` -> bundled fallback token.
- Machine-readable output keeps stable English schema keys regardless of locale, so agent integrations do not break.
- The current schema version is `--schema-version v1`.
- Run `go test ./...` to execute the current unit tests.
