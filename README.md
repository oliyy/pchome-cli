# PChome 24h 命令列工具 (CLI)

<p align="center">
  <img src="./docs/banner.png" alt="pchome-cli" width="600" />
</p>

快速、易於整合腳本的 PChome 24h 購物命令列工具。支援商品搜尋、檢視商品詳情、比較多項商品以及取得商品推薦。內建 JSON 優先輸出、適合人類閱讀的表格格式，以及標準化的資料結構（Schema）。

## 功能特色

- **搜尋 (Search)** - 支援透過關鍵字搜尋商品，並可依照品牌、價格區間、評價、庫存狀態、24h 到貨、排序方式進行篩選，同時支援自訂輸出欄位
- **檢視 (View)** - 顯示詳細的商品資訊，包含規格、圖片以及相關警告提示
- **推薦 (Recommend)** - 根據指定商品取得相關的推薦商品，並可選擇顯示推薦原因
- **比較 (Compare)** - 並排比較多項商品，支援自訂顯示欄位
- **建議 (Suggest)** - 提供搜尋關鍵字的自動補全與建議
- **多種輸出格式** - 支援適合閱讀的文字表格格式、供腳本使用的 JSON 格式，以及適合串流或 AI 代理程式（Agent）整合的 NDJSON 格式
- **標準化資料結構 (Schema)** - 無論語系為何，皆維持穩定的英文欄位鍵值（`v1`），確保 AI 代理程式或腳本整合不會因語系切換而失效
- **多國語系 (i18n)** - 支援繁體中文（預設）與英文介面
- **彈性的商品輸入格式** - 支援直接輸入商品編號（例如 `DRAA5K-A900JOK9O`）、帶有後綴的編號（例如 `DRAA5K-A900JOK9O-000`），或是完整的 PChome 商品網址
- **高度可設定** - 可透過 `~/.pchome/config.toml` 設定各指令的預設值、欄位排序以及輸出偏好

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

### 預先編譯的二進位檔 (Prebuilt Binaries)

您可以從 [GitHub Releases](https://github.com/oliy/pchome-cli/releases) 下載符合您作業系統的壓縮檔，解壓縮後將 `pchome` 放置於系統的 `PATH` 路徑下。

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

- `pchome --help` 會顯示最上層的指令群組。
- 若要查看特定指令的說明，可使用 `pchome <command> --help`。

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

# 顯示每項商品的推薦原因
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
- 帶後綴的商品編號：`DRAA5K-A900JOK9O-000`
- 完整的 PChome 商品網址：`https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O`

## 輸出格式

### 文字 (Text)

預設選項，以簡潔的表格呈現適合人類閱讀的格式：

```bash
pchome search "掃地機器人" --limit 3
```

### JSON

適合腳本與自動化流程的機器可讀格式：

```bash
pchome search "掃地機器人" --limit 3 --format json
pchome view DRAA5K-A900JOK9O --format json
```

### NDJSON

每行一個 JSON 物件，非常適合串流處理與 AI 代理程式整合：

```bash
pchome search "掃地機器人" --limit 5 --format ndjson
pchome recommend DMBL53-A900JDNJS --top 5 --format ndjson
```

`ndjson` 格式支援 `search`、`recommend`、`compare` 與 `suggest` 指令。

資料輸出至 stdout，錯誤訊息與進度則輸出至 stderr，方便您進行管線（Piping）處理：

```bash
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000)'
```

## 設定與組態

在啟動時，`pchome` 會確認以下路徑的設定檔是否存在：

```bash
~/.pchome/config.toml
```

若檔案不存在，系統會自動建立並填入所有預設設定。

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

- `columns = []` 代表「使用該指令的內建預設欄位排序」。
- 若您設定了 `columns`，該列表將會成為該指令的預設顯示欄位。
- `i18n.language` 目前支援 `zh-TW` 與 `en`。
- 設定載入器會拒絕未知的鍵值，以防止拼寫錯誤被靜默忽略。

## 實際應用範例

### 搜尋並篩選商品

```bash
# 帶有價格區間與庫存篩選的搜尋
pchome search "掃地機器人" --min-price 5000 --max-price 15000 --in-stock

# 依照品牌與最低評價進行篩選
pchome search "掃地機器人" --brand Roborock --min-rating 4.8

# 僅限 24h 到貨，並依價格由低至高排序
pchome search "掃地機器人" --arrival-24h --sort price-asc
```

### 取得附帶原因的推薦商品

```bash
pchome recommend DMAB3X-A900EVNNM --top 10 --why
```

### 將 JSON 輸出管線連接至 jq

```bash
# 擷取低於特定價格門檻的商品名稱
pchome search "掃地機器人" --format json | jq '.products[] | select(.price < 10000) | .name'
```

### 使用 NDJSON 進行串流處理

```bash
pchome search "掃地機器人" --limit 20 --format ndjson | while read -r line; do
  echo "$line" | jq -r '.name'
done
```
