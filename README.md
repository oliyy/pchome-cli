# PChome 24h 命令列工具 (CLI)

<p align="center">
  <img src="./docs/banner.png" alt="pchome-cli" width="600" />
</p>

快速且易於整合腳本的 PChome 24h 商品查詢 CLI 工具。支援商品搜尋、檢視商品詳情、商品比較與推薦功能。內建 JSON 優先輸出、人類可讀的表格格式，以及標準化的資料結構（Schema）。

## 功能特色

- **搜尋 (Search)** - 透過關鍵字搜尋商品，可依照品牌、價格區間、評價、庫存狀態、24h 到貨、排序方式進行篩選，同時支援自訂輸出欄位
- **檢視 (View)** - 顯示商品資訊，包含規格、圖片以及相關警告
- **推薦 (Recommend)** - 根據商品取得相關商品推薦，並可選擇顯示推薦原因
- **比較 (Compare)** - 比較多項商品，支援自訂顯示欄位
- **建議 (Suggest)** - 關鍵字自動補全與建議
- **多種輸出格式** - 包含適合閱讀的文字表格、供腳本使用的 JSON 格式，以及供串流或 AI agent 整合的 NDJSON 格式
- **標準化資料結構 (Schema)** - 無論語系皆維持英文欄位鍵值（`v1`），確保 agent 或腳本整合不受語系切換影響
- **多國語系 (i18n)** - 支援繁體中文（預設）與英文介面
- **商品輸入格式** - 支援直接輸入商品編號（例如 `DRAA5K-A900JOK9O`）、含後綴的編號（例如 `DRAA5K-A900JOK9O-000`），或是完整的 PChome 商品網址
- **設定** - 可透過 `~/.pchome/config.toml` 設定各指令預設值、欄位排序以及輸出偏好

## 安裝方式

### Homebrew

```bash
brew install oliyy/tap/pchome-cli
```

### Scoop (Windows)

```powershell
scoop bucket add oliyy https://github.com/oliyy/scoop-bucket.git
scoop install oliyy/pchome-cli
```

### 執行檔 (Prebuilt Binaries)

從 [GitHub Releases](https://github.com/oliy/pchome-cli/releases) 下載壓縮檔，解壓縮後將 `pchome` 放置於系統 `PATH` 路徑下。

### 從原始碼編譯

```bash
git clone https://github.com/oliy/pchome-cli.git
cd pchome-cli
make build
```

執行：

```bash
./bin/pchome --help
```

或全域安裝：

```bash
go install github.com/oliy/pchome-cli/cmd/pchome@latest
```

取得協助：

- `pchome --help` 顯示最上層的指令群組。
- 使用 `pchome <command> --help` 查看特定指令說明。

## 快速開始

```bash
# 搜尋商品
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# 檢視商品詳情
pchome view DRAA5K-A900JOK9O

# 取得商品推薦
pchome recommend DMBL53-A900JDNJS --top 8

# 比較多項商品
pchome compare DRAA5K-A900JOK9O DMBL53-A900JDNJS

# 取得搜尋關鍵字建議
pchome suggest "掃地機"
```

## 指令說明

### 搜尋 (Search)

```bash
# 基本搜尋
pchome search "掃地機器人"

# 價格區間與庫存篩選
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# 品牌與評價篩選
pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 僅限 24h 到貨，並依價格由低至高排序
pchome search "掃地機器人" --arrival-24h --sort price-asc

# 自訂文字輸出欄位
pchome search "掃地機器人" \
  --columns "#,name,price,rating,reviews,24h,brand,qty,url"
```

### 檢視 (View)

```bash
pchome view DRAA5K-A900JOK9O
pchome view DRAA5K-A900JOK9O --format json
pchome view https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O
```

### 推薦 (Recommend)

```bash
# 基本推薦
pchome recommend DMBL53-A900JDNJS --top 8

# 顯示商品的推薦原因
pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

### 比較 (Compare)

```bash
pchome compare DRAA5K-A900JOK9O DMBL53-A900JDNJS
```

### 建議 (Suggest)

```bash
pchome suggest "掃地機"
```

### 商品輸入格式

`PRODUCT` 參數支援以下格式：

- 原始商品編號：`DRAA5K-A900JOK9O`
- 有後綴的商品編號：`DRAA5K-A900JOK9O-000`
- PChome 商品網址：`https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O`

## 輸出格式

### 文字 (Text)

以簡潔的表格呈現，方便人類閱讀 (預設)：

```bash
pchome search "掃地機器人" --limit 3
```

### JSON

輸出機器可讀格式，方便腳本與自動化流程使用：

```bash
pchome search "掃地機器人" --limit 3 --format json
pchome view DRAA5K-A900JOK9O --format json
```

### NDJSON

每行一筆資料，適用於串流處理與 agent 整合：

```bash
pchome search "掃地機器人" --limit 5 --format ndjson
pchome recommend DMBL53-A900JDNJS --top 5 --format ndjson
```

`ndjson` 格式支援 `search`、`recommend`、`compare` 與 `suggest` 指令。

將資料輸出至 stdout，錯誤訊息與進度則輸出至 stderr，方便進行管線（Piping）處理：

```bash
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000)'
```

## 設定與組態

在啟動時，`pchome`會確認設定檔是否存在以下路徑：

```bash
~/.pchome/config.toml
```

若檔案不存在，系統會自動建立並填入所有預設值。

設定優先順序：

1. CLI 指令選項（Flags）
2. `~/.pchome/config.toml` 設定檔
3. 內建預設值

範例：

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

備註：

- `columns = []` 代表「使用該指令的預設欄位排序」。
- 更改 `columns` 後，該列表即為預設欄位順序。
- `i18n.language` 目前支援 `zh-TW` 與 `en`。
- 為避免拼寫錯誤被靜默忽略，載入器會拒絕未知鍵值。

## 實際應用範例

### 搜尋並篩選商品

```bash
# 價格區間與庫存篩選的搜尋
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# 依照品牌與最低評價進行篩選
pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 僅限 24h 到貨，並依價格由低至高排序
pchome search "掃地機器人" --arrival-24h --sort price-asc
```

### 取得推薦商品的原因

```bash
pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

### 將 JSON 輸出給 jq 處理

```bash
# 取得低於特定價格的商品名稱
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000) | .name'
```

### 使用 NDJSON 進行串流處理

```bash
pchome search "掃地機器人" --limit 20 --format ndjson | while read -r line; do
  echo "$line" | jq -r '.name'
done
```
