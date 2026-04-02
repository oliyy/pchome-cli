# PChome 24h - Command Line Interface

<p align="center">
  <img src="./docs/banner.png" alt="pchome-cli" width="600" />
</p>

Fast, script-friendly CLI for searching PChome 24h products, viewing product details, comparing products, and getting recommendations. JSON-first output, human-friendly tables, and a normalized schema built in.

## Features

- **Search** - search products with filters for brand, price range, rating, stock status, 24h arrival, sorting, and custom column output
- **View** - display detailed product information including specs, images, and warnings
- **Recommend** - get product recommendations based on a reference product, with optional reasoning
- **Compare** - compare multiple products side-by-side with customizable columns
- **Suggest** - get autocomplete suggestions for search queries
- **Multi-format output** - human-friendly text tables, JSON, and NDJSON for scripting and agent integrations
- **Normalized schema** - stable English schema keys (`v1`) regardless of locale, so agent integrations do not break
- **i18n** - Traditional Chinese (default) and English interface
- **Flexible product input** - accepts raw IDs (`DRAA5K-A900JOK9O`), suffixed IDs (`DRAA5K-A900JOK9O-000`), or full PChome URLs
- **Configurable** - per-command defaults, column ordering, and output preferences via `~/.pchome/config.toml`

## Installation

### Homebrew

```bash
brew install oliyy/tap/pchome-cli
```

### Scoop (Windows)

```powershell
scoop bucket add oliyy https://github.com/oliyy/scoop-bucket.git
scoop install oliyy/pchome-cli
```

### Prebuilt Binaries

Download the archive for your platform from [GitHub Releases](https://github.com/oliy/pchome-cli/releases), extract it, and place `pchome` in your `PATH`.

### Build from Source

```bash
git clone https://github.com/oliy/pchome-cli.git
cd pchome-cli
make build
```

Run:

```bash
./bin/pchome --help
```

Or install globally:

```bash
go install github.com/oliy/pchome-cli/cmd/pchome@latest
```

Help:

- `pchome --help` shows top-level command groups.
- You can get help for a specific command with `pchome <command> --help`.

## Quick Start

```bash
# Search for products
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# View a product
pchome view DRAA5K-A900JOK9O

# Get recommendations
pchome recommend DMBL53-A900JDNJS --top 8

# Compare products
pchome compare DRAA5K-A900JOK9O DMBL53-A900JDNJS

# Autocomplete suggestions
pchome suggest "掃地機"
```

## Commands

### Search

```bash
# Basic search
pchome search "掃地機器人"

# Price range and stock filter
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# Brand and rating filter
pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 24h arrival only, sorted by price
pchome search "掃地機器人" --arrival-24h --sort price-asc

# Custom text columns
pchome search "掃地機器人" \
  --columns "#,name,price,rating,reviews,24h,brand,qty,url"
```

### View

```bash
pchome view DRAA5K-A900JOK9O
pchome view DRAA5K-A900JOK9O --format json
pchome view https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O
```

### Recommend

```bash
# Basic recommendations
pchome recommend DMBL53-A900JDNJS --top 8

# Explain why each item was recommended
pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

### Compare

```bash
pchome compare DRAA5K-A900JOK9O DMBL53-A900JDNJS
```

### Suggest

```bash
pchome suggest "掃地機"
```

### Product Input

`PRODUCT` can be:

- A raw product ID like `DRAA5K-A900JOK9O`
- A suffixed product ID like `DRAA5K-A900JOK9O-000`
- A full PChome product URL like `https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O`

## Output Formats

### Text

Human-readable output with compact tables (default):

```bash
pchome search "掃地機器人" --limit 3
```

### JSON

Machine-readable output for scripting and automation:

```bash
pchome search "掃地機器人" --limit 3 --format json
pchome view DRAA5K-A900JOK9O --format json
```

### NDJSON

Line-delimited items for streaming and agent integrations:

```bash
pchome search "掃地機器人" --limit 5 --format ndjson
pchome recommend DMBL53-A900JDNJS --top 5 --format ndjson
```

`ndjson` is supported on `search`, `recommend`, `compare`, and `suggest`.

Data goes to stdout, errors and progress to stderr for clean piping:

```bash
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000)'
```

## Configuration

On startup, `pchome` ensures a config file exists at:

```bash
~/.pchome/config.toml
```

If the file is missing, it is created automatically with the full default configuration.

Config precedence:

1. CLI flags
2. `~/.pchome/config.toml`
3. Built-in defaults

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

- `columns = []` means "use the built-in default column order for that command".
- If you set `columns`, that list becomes the default column order for the command.
- `i18n.language` currently supports `zh-TW` and `en`.
- The config loader rejects unknown keys so typos do not silently get ignored.

## Examples

### Search and filter products

```bash
# Search with price range
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# Filter by brand and rating
pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 24h delivery only, sorted by price
pchome search "掃地機器人" --arrival-24h --sort price-asc
```

### Get recommendations with reasoning

```bash
pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

### Pipe JSON output to jq

```bash
# Extract product names under a price threshold
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000) | .name'
```

### Stream results with NDJSON

```bash
pchome search "掃地機器人" --limit 20 --format ndjson | while read -r line; do
  echo "$line" | jq -r '.name'
done
```

## 

