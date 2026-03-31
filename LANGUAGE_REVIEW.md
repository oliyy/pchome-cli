# 語言校對清單

這份文件是給協助校對繁體中文（台灣）用語的朋友使用。

除了 `README.md` 之外，請優先檢查下面這些地方，確認：

- 用語是否自然、符合台灣常見說法
- 語氣是否一致
- CLI 介面是否清楚好懂
- 標點、空格、大小寫是否一致
- 英文是否有不小心漏在使用者會看到的地方

---

## 1. 最優先：翻譯字典

請先看：

- `pkg/i18n/i18n.go`

這是目前所有主要使用者介面文字的集中位置。請特別檢查：

- 根指令說明
- `--help` 頁面標題
- 各指令說明與範例
- 各個 flags 的描述
- 表格欄位名稱
- 商品詳情頁欄位名稱
- 搜尋／推薦／比較／建議的摘要句子
- 錯誤訊息與警告訊息

建議特別留意這些詞是否自然：

- 「檢視」是否要改成「查看」
- 「供貨資訊」是否合適
- 「標記」是否夠清楚
- 「銷售名稱」/「別名」/「標語」是否符合實際語境
- 「數量」是否應保留或改成更貼近購物語境的說法
- 「建議」是否足以表達搜尋建議

---

## 2. 第二優先：實際 CLI 顯示結果

請看這些 golden fixtures，它們代表目前實際輸出的樣子：

- `cmd/testdata/help/root_zh_tw.golden`
- `cmd/testdata/render/search_zh_tw.golden`
- `cmd/testdata/render/product_detail_zh_tw.golden`

這一層最重要，因為字串單看可能沒問題，但放進 CLI 介面後可能會顯得：

- 太長
- 太硬
- 不像台灣常見產品介面
- 前後文不自然

請特別檢查：

- 欄位標題是否清楚
- 搜尋摘要句子是否通順
- 商品詳情頁每一行的名稱是否自然
- `24h`、`Y/N`、英文欄位名是否需要保留或改寫

---

## 3. 第三優先：Help / 輸出格式的上下文

如果想確認某些字串是怎麼被組裝出來的，可以看：

- `cmd/root.go`
- `cmd/localize.go`
- `cmd/render.go`

這些檔案不一定要逐字校對，但如果你看到某些句子覺得卡卡的，通常可以回頭對照這幾個檔案，理解：

- 它是單獨的一段文字
- 還是被格式化後跟別的欄位拼在一起

---

## 4. 也請順手看：設定檔相關的使用者文字

雖然設定檔 key 會維持英文，但使用者仍會接觸到下列內容：

- `pkg/config/config.go`

請檢查：

- `DefaultTOML()` 中的註解
- 如果有需要顯示給一般使用者的設定說明，是否夠清楚

補充：

- `config.toml` 的 key 名稱像 `language`、`limit`、`columns`、`show_url` 都是刻意保留英文，不需要改成中文

---

## 5. 非集中翻譯、但使用者可能看得到的錯誤訊息

下面這些檔案裡還有一些不是放在 `pkg/i18n/i18n.go` 裡的字串。

它們不一定每次都會出現，但如果出現，使用者會看得到，所以也值得確認：

- `pkg/config/config.go`
  - 設定檔讀取與驗證錯誤
- `pkg/catalog/refs.go`
  - 商品 ID / URL 格式錯誤
- `pkg/catalog/service.go`
  - 搜尋、比較、推薦、建議相關錯誤與 warning
- `pkg/httpx/client.go`
  - HTTP 層錯誤
- `pkg/pchome/search/client.go`
  - 搜尋 API 參數或回應錯誤
- `pkg/pchome/prodapi/client.go`
  - 商品 API / JSONP 解析錯誤

這一區請特別幫忙判斷兩件事：

1. 這些訊息是否應該也改成繁體中文
2. 如果保留英文，是否合理，還是會讓一般使用者看不懂

---

## 6. 建議實際操作一次

如果方便，請不要只看程式碼，最好實際跑幾個指令：

```bash
go run ./cmd/pchome --help
go run ./cmd/pchome search "掃地機器人" --limit 3
go run ./cmd/pchome view DRAA5K-A900JOK9O
go run ./cmd/pchome suggest "掃地機"
```

如果你有可用的推薦 token，也可以再看：

```bash
go run ./cmd/pchome recommend DMBL53-A900JDNJS --top 3 --why
```

實際執行時，請特別留意：

- 有沒有字太長導致表格不好看
- 哪些詞看起來像中國用語、香港用語、或工程內部用語
- 哪些地方雖然正確，但不像台灣電商或台灣 CLI 使用者會習慣的說法

---

## 7. 回饋時請這樣標註

為了方便修改，回饋時最好附上：

- 檔案路徑
- 原文
- 建議改成什麼
- 為什麼

例如：

```text
pkg/i18n/i18n.go
原文：供貨資訊
建議：庫存與配送
原因：一般使用者比較容易理解，也更像購物情境
```

或：

```text
cmd/testdata/render/product_detail_zh_tw.golden
原文：標記
建議：特殊條件
原因：目前太抽象，看不出是在說 Prime / 訂單折扣之類的欄位
```

---

## 8. 本次不用檢查的地方

下面這些通常不用花時間：

- `README.en.md`
- `cmd/testdata/help/root_en.golden`
- Go test 裡的 `t.Fatalf(...)`
- struct field 名稱、JSON key、config key

這些地方主要是給英文文件、測試或程式介面使用，不是本次繁體中文文案校對重點。
