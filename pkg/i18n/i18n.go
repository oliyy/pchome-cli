package i18n

import (
	"fmt"
	"strings"
)

type Language string

const (
	LangZHTW Language = "zh-TW"
	LangEN   Language = "en"
)

type Key string

const (
	RootShort Key = "root.short"
	RootLong  Key = "root.long"

	GroupShopping  Key = "group.shopping"
	GroupDiscovery Key = "group.discovery"

	HelpHeadingUsage              Key = "help.heading.usage"
	HelpHeadingAliases            Key = "help.heading.aliases"
	HelpHeadingExamples           Key = "help.heading.examples"
	HelpHeadingAvailableCommands  Key = "help.heading.available_commands"
	HelpHeadingAdditionalCommands Key = "help.heading.additional_commands"
	HelpHeadingFlags              Key = "help.heading.flags"
	HelpHeadingGlobalFlags        Key = "help.heading.global_flags"
	HelpHeadingAdditionalTopics   Key = "help.heading.additional_topics"
	HelpMoreInfo                  Key = "help.more_info"
	HelpFlagForCommand            Key = "help.flag.for_command"
	HelpFlagForThisCommand        Key = "help.flag.for_this_command"
	VersionFlagForCommand         Key = "version.flag.for_command"
	VersionFlagForThisCommand     Key = "version.flag.for_this_command"
	HelpCommandShort              Key = "help.command.short"
	HelpCommandLong               Key = "help.command.long"
	CompletionCommandShort        Key = "completion.command.short"

	RootFlagFormat        Key = "root.flag.format"
	RootFlagSchemaVersion Key = "root.flag.schema_version"
	RootFlagTimeout       Key = "root.flag.timeout"
	RootFlagNameWidth     Key = "root.flag.name_width"

	SearchShort          Key = "search.short"
	SearchExample        Key = "search.example"
	SearchFlagCategory   Key = "search.flag.category"
	SearchFlagBrand      Key = "search.flag.brand"
	SearchFlagSort       Key = "search.flag.sort"
	SearchFlagPage       Key = "search.flag.page"
	SearchFlagPageSize   Key = "search.flag.page_size"
	SearchFlagLimit      Key = "search.flag.limit"
	SearchFlagMinPrice   Key = "search.flag.min_price"
	SearchFlagMaxPrice   Key = "search.flag.max_price"
	SearchFlagMinRating  Key = "search.flag.min_rating"
	SearchFlagInStock    Key = "search.flag.in_stock"
	SearchFlagArrival24h Key = "search.flag.arrival_24h"
	SearchFlagColumns    Key = "search.flag.columns"
	SearchFlagShowURL    Key = "search.flag.show_url"
	SearchFlagCompact    Key = "search.flag.compact"
	SearchFlagWide       Key = "search.flag.wide"

	ViewShort                Key = "view.short"
	ViewExample              Key = "view.example"
	ViewErrNDJSONUnsupported Key = "view.err.ndjson_unsupported"

	RecommendShort       Key = "recommend.short"
	RecommendExample     Key = "recommend.example"
	RecommendFlagTop     Key = "recommend.flag.top"
	RecommendFlagColumns Key = "recommend.flag.columns"
	RecommendFlagShowURL Key = "recommend.flag.show_url"
	RecommendFlagShowWhy Key = "recommend.flag.show_why"
	RecommendFlagCompact Key = "recommend.flag.compact"
	RecommendFlagWide    Key = "recommend.flag.wide"

	CompareShort       Key = "compare.short"
	CompareExample     Key = "compare.example"
	CompareFlagColumns Key = "compare.flag.columns"
	CompareFlagShowURL Key = "compare.flag.show_url"
	CompareFlagCompact Key = "compare.flag.compact"
	CompareFlagWide    Key = "compare.flag.wide"

	SuggestShort     Key = "suggest.short"
	SuggestExample   Key = "suggest.example"
	SuggestFlagLimit Key = "suggest.flag.limit"

	RenderSearchSummary    Key = "render.search.summary"
	RenderFilters          Key = "render.filters"
	RenderRecommendSummary Key = "render.recommend.summary"
	RenderCompareSummary   Key = "render.compare.summary"
	RenderSuggestSummary   Key = "render.suggest.summary"
	RenderWarning          Key = "render.warning"

	FieldID             Key = "field.id"
	FieldURL            Key = "field.url"
	FieldBrand          Key = "field.brand"
	FieldSalesName      Key = "field.sales_name"
	FieldNick           Key = "field.nick"
	FieldTagline        Key = "field.tagline"
	FieldDescription    Key = "field.description"
	FieldPrice          Key = "field.price"
	FieldListPrice      Key = "field.list_price"
	FieldLowestObserved Key = "field.lowest_observed"
	FieldDiscount       Key = "field.discount"
	FieldRating         Key = "field.rating"
	FieldAvailability   Key = "field.availability"
	FieldFlags          Key = "field.flags"
	FieldCategories     Key = "field.categories"
	FieldBrandAliases   Key = "field.brand_aliases"
	FieldPrimaryImage   Key = "field.primary_image"
	FieldImages         Key = "field.images"
	FieldMoreImages     Key = "field.more_images"
	FieldNotAvailable   Key = "field.not_available"

	AvailabilitySummary Key = "availability.summary"
	AvailabilityStock   Key = "availability.stock"
	AvailabilityShip    Key = "availability.ship"
	FlagsSummary        Key = "flags.summary"
	FlagPrimeOnly       Key = "flag.prime_only"
	FlagOrderDiscount   Key = "flag.order_discount"
	RatingWithReviews   Key = "rating.with_reviews"

	FilterCategory   Key = "filter.category"
	FilterBrand      Key = "filter.brand"
	FilterMinPrice   Key = "filter.min_price"
	FilterMaxPrice   Key = "filter.max_price"
	FilterMinRating  Key = "filter.min_rating"
	FilterInStock    Key = "filter.in_stock"
	FilterArrival24h Key = "filter.arrival_24h"

	HeaderIndex       Key = "header.index"
	HeaderPrice       Key = "header.price"
	HeaderList        Key = "header.list"
	HeaderDiscount    Key = "header.discount"
	HeaderRating      Key = "header.rating"
	HeaderReviews     Key = "header.reviews"
	Header24h         Key = "header.24h"
	HeaderQty         Key = "header.qty"
	HeaderBrand       Key = "header.brand"
	HeaderName        Key = "header.name"
	HeaderProductID   Key = "header.product_id"
	HeaderURL         Key = "header.url"
	HeaderDescription Key = "header.description"
	HeaderScore       Key = "header.score"
	HeaderWhy         Key = "header.why"
	HeaderSuggestion  Key = "header.suggestion"

	ErrInvalidFormat        Key = "err.invalid_format"
	ErrInvalidSchemaVersion Key = "err.invalid_schema_version"
	ErrUnknownColumn        Key = "err.unknown_column"
)

type Catalog struct {
	lang Language
}

func New(lang Language) *Catalog {
	return &Catalog{lang: lang}
}

func (c *Catalog) Language() Language {
	if c == nil {
		return LangZHTW
	}
	return c.lang
}

func (c *Catalog) Text(key Key) string {
	lang := c.Language()
	if value, ok := translations[lang][key]; ok {
		return value
	}
	if value, ok := translations[LangEN][key]; ok {
		return value
	}
	return string(key)
}

func (c *Catalog) Sprintf(key Key, args ...any) string {
	return fmt.Sprintf(c.Text(key), args...)
}

func ParseLanguage(s string) (Language, error) {
	switch normalized := strings.TrimSpace(s); normalized {
	case "", string(LangZHTW):
		return LangZHTW, nil
	case string(LangEN):
		return LangEN, nil
	default:
		return "", fmt.Errorf("unsupported language %q (use %s or %s)", s, LangZHTW, LangEN)
	}
}

func SupportedLanguages() []Language {
	return []Language{LangZHTW, LangEN}
}

var translations = map[Language]map[Key]string{
	LangZHTW: {
	        RootShort: "搜尋、檢視、比較與推薦 PChome 24h 商品",
	        RootLong: `搜尋、檢視、比較與推薦 PChome 24h 商品。

範例：
  pchome search "掃地機器人" --min-price 5000 --max-price 15000
  pchome view DMBL53-A900JDNJS
  pchome recommend https://24h.pchome.com.tw/prod/DMBL53-A900JDNJS --why
  pchome compare DMBL53-A900JDNJS DMBL1C-A900JA04J`,

	        GroupShopping:  "購物指令",
	        GroupDiscovery: "探索指令",

	        HelpHeadingUsage:              "使用方式：",
	        HelpHeadingAliases:            "別名：",
	        HelpHeadingExamples:           "範例：",
	        HelpHeadingAvailableCommands:  "可用指令：",
	        HelpHeadingAdditionalCommands: "其他指令：",
	        HelpHeadingFlags:              "選項：",
	        HelpHeadingGlobalFlags:        "全域選項：",
	        HelpHeadingAdditionalTopics:   "其他說明主題：",
	        HelpMoreInfo:                  "使用 \"%s [command] --help\" 取得指令的詳細資訊。",
	        HelpFlagForCommand:            "顯示 %s 的說明",
	        HelpFlagForThisCommand:        "顯示此指令的說明",
	        VersionFlagForCommand:         "顯示 %s 的版本資訊",
	        VersionFlagForThisCommand:     "顯示此指令的版本資訊",
	        HelpCommandShort:              "顯示指令的說明",
	        HelpCommandLong:               "顯示應用程式中任何指令的說明。\n輸入 `%s help [指令名稱]` 即可查看完整內容。",
	        CompletionCommandShort:        "產生指定 shell 的自動補全腳本",

	        RootFlagFormat:        "輸出格式：text|json|ndjson",
	        RootFlagSchemaVersion: "機器可讀輸出的 schema 版本",
	        RootFlagTimeout:       "API 請求逾時時間",
	        RootFlagNameWidth:     "商品名稱欄位寬度",

	        SearchShort:          "搜尋商品",
	        SearchExample:        "  pchome search \"掃地機器人\" --min-price 5000 --max-price 15000 --in-stock",
	        SearchFlagCategory:   "分類識別碼",
	        SearchFlagBrand:      "品牌篩選（部分符合即可）",
	        SearchFlagSort:       "排序方式：relevance（相關度）|price-asc（價格低至高）|price-desc（價格高至低）|newest（最新上架）|best-selling（最暢銷）",
	        SearchFlagPage:       "起始頁碼",
	        SearchFlagPageSize:   "每次 API 請求的資料筆數",
	        SearchFlagLimit:      "最多顯示幾筆商品",
	        SearchFlagMinPrice:   "最低價格",
	        SearchFlagMaxPrice:   "最高價格",
	        SearchFlagMinRating:  "最低評價",
	        SearchFlagInStock:    "只顯示有庫存的商品",
	        SearchFlagArrival24h: "只顯示支援 24h 到貨的商品",
	        SearchFlagColumns:    "以逗號分隔的顯示欄位 (#,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url,desc)",
	        SearchFlagShowURL:    "在文字輸出中包含商品網址",
	        SearchFlagCompact:    "使用較窄的商品名稱欄位",
	        SearchFlagWide:       "使用較寬的商品名稱欄位",

	        ViewShort:                "檢視商品詳情",
	        ViewExample:              "  pchome view https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O",
	        ViewErrNDJSONUnsupported: "view 指令不支援 ndjson 格式，請改用 --format json 或 text",

	        RecommendShort:       "取得推薦商品",
	        RecommendExample:     "  pchome recommend DMBL53-A900JDNJS --top 8 --why",
	        RecommendFlagTop:     "最多顯示幾筆推薦商品",
	        RecommendFlagColumns: "以逗號分隔的顯示欄位 (#,score,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url,why)",
	        RecommendFlagShowURL: "在文字輸出中包含商品網址",
	        RecommendFlagShowWhy: "在文字輸出中顯示推薦原因",
	        RecommendFlagCompact: "使用較窄的推薦商品名稱欄位",
	        RecommendFlagWide:    "使用較寬的推薦商品名稱欄位",

	        CompareShort:       "比較多項商品",
	        CompareExample:     "  pchome compare DMBL53-A900JDNJS DMBL1C-A900JA04J",
	        CompareFlagColumns: "以逗號分隔的顯示欄位 (#,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url)",
	        CompareFlagShowURL: "在文字輸出中包含商品網址",
	        CompareFlagCompact: "使用較窄的比較商品名稱欄位",
	        CompareFlagWide:    "使用較寬的比較商品名稱欄位",

	        SuggestShort:     "提供搜尋建議",
	        SuggestExample:   "  pchome suggest \"掃地機\"",
	        SuggestFlagLimit: "最多顯示幾筆搜尋建議",

	        RenderSearchSummary:    "搜尋 %q | 傳回 %d 筆（共 %d 筆）| 第 %d 頁 | 已掃描 %d 頁 | 排序=%s\n",
	        RenderFilters:          "篩選條件：%s\n",
	        RenderRecommendSummary: "商品 %s 的推薦結果 | 傳回 %d 筆 | 耗時 %.0f 毫秒\n",
	        RenderCompareSummary:   "比較 %d 項商品\n",
	        RenderSuggestSummary:   "%q 的搜尋建議 | 傳回 %d 筆\n",
	        RenderWarning:          "警告：%s\n",

	        FieldID:             "商品編號",
	        FieldURL:            "網址",
	        FieldBrand:          "品牌",
	        FieldSalesName:      "商品名稱",
	        FieldNick:           "簡稱",
	        FieldTagline:        "促銷標語",
	        FieldDescription:    "商品描述",
	        FieldPrice:          "售價",
	        FieldListPrice:      "原價",
	        FieldLowestObserved: "歷史低價",
	        FieldDiscount:       "折扣",
	        FieldRating:         "評價",
	        FieldAvailability:   "庫存狀態",
	        FieldFlags:          "活動標籤",
	        FieldCategories:     "分類",
	        FieldBrandAliases:   "品牌別名",
	        FieldPrimaryImage:   "主圖",
	        FieldImages:         "商品圖片",
	        FieldMoreImages:     "更多圖片",
	        FieldNotAvailable:   "無",

	        AvailabilitySummary: "%s=%s | 24h=%s | %s=%s\n",
	        AvailabilityStock:   "庫存",
	        AvailabilityShip:    "出貨",
	        FlagsSummary:        "%s=%s | %s=%s\n",
	        FlagPrimeOnly:       "Prime 限定",
	        FlagOrderDiscount:   "結帳折扣",
	        RatingWithReviews:   "（%s 則評價）",

	        FilterCategory:   "分類",
	        FilterBrand:      "品牌",
	        FilterMinPrice:   "最低價",
	        FilterMaxPrice:   "最高價",
	        FilterMinRating:  "最低評價",
	        FilterInStock:    "有庫存",
	        FilterArrival24h: "24h到貨",

	        HeaderIndex:       "#",
	        HeaderPrice:       "網路價",
	        HeaderList:        "市價",
	        HeaderDiscount:    "折扣%",
	        HeaderRating:      "評價",
	        HeaderReviews:     "評價數量",
	        Header24h:         "24h",
	        HeaderQty:         "庫存",
	        HeaderBrand:       "品牌",
	        HeaderName:        "商品名稱",
	        HeaderProductID:   "商品編號",
	        HeaderURL:         "網址",
	        HeaderDescription: "商品描述",
	        HeaderScore:       "相關度",
	        HeaderWhy:         "推薦原因",
	        HeaderSuggestion:  "搜尋建議",

	        ErrInvalidFormat:        "不支援的輸出格式 %q（請使用 text、json 或 ndjson）",
	        ErrInvalidSchemaVersion: "不支援的 schema 版本 %q（請使用 %s）",
	        ErrUnknownColumn:        "未知欄位 %q",
	},	LangEN: {
		RootShort: "Search, inspect, compare, and recommend PChome 24h products",
		RootLong: `Search, inspect, compare, and recommend PChome 24h products.

Examples:
  pchome search "掃地機器人" --min-price 5000 --max-price 15000
  pchome view DMBL53-A900JDNJS
  pchome recommend https://24h.pchome.com.tw/prod/DMBL53-A900JDNJS --why
  pchome compare DMBL53-A900JDNJS DMBL1C-A900JA04J`,

		GroupShopping:  "Shopping Commands",
		GroupDiscovery: "Discovery Commands",

		HelpHeadingUsage:              "Usage:",
		HelpHeadingAliases:            "Aliases:",
		HelpHeadingExamples:           "Examples:",
		HelpHeadingAvailableCommands:  "Available Commands:",
		HelpHeadingAdditionalCommands: "Additional Commands:",
		HelpHeadingFlags:              "Flags:",
		HelpHeadingGlobalFlags:        "Global Flags:",
		HelpHeadingAdditionalTopics:   "Additional help topics:",
		HelpMoreInfo:                  "Use \"%s [command] --help\" for more information about a command.",
		HelpFlagForCommand:            "help for %s",
		HelpFlagForThisCommand:        "help for this command",
		VersionFlagForCommand:         "version for %s",
		VersionFlagForThisCommand:     "version for this command",
		HelpCommandShort:              "Help about any command",
		HelpCommandLong:               "Help provides help for any command in the application.\nSimply type `%s help [path to command]` for full details.",
		CompletionCommandShort:        "Generate the autocompletion script for the specified shell",

		RootFlagFormat:        "Output format: text|json|ndjson",
		RootFlagSchemaVersion: "Machine-readable schema version",
		RootFlagTimeout:       "Request timeout",
		RootFlagNameWidth:     "Default text width for product names",

		SearchShort:          "Search products",
		SearchExample:        "  pchome search \"掃地機器人\" --min-price 5000 --max-price 15000 --in-stock",
		SearchFlagCategory:   "Category id",
		SearchFlagBrand:      "Brand filter (substring match)",
		SearchFlagSort:       "Sort: relevance|price-asc|price-desc|newest|best-selling",
		SearchFlagPage:       "Start page to search from",
		SearchFlagPageSize:   "Page size per upstream request",
		SearchFlagLimit:      "Maximum number of products to return",
		SearchFlagMinPrice:   "Minimum price",
		SearchFlagMaxPrice:   "Maximum price",
		SearchFlagMinRating:  "Minimum rating",
		SearchFlagInStock:    "Only include products with stock available",
		SearchFlagArrival24h: "Only include products eligible for 24h arrival",
		SearchFlagColumns:    "Comma-separated text columns (#,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url,desc)",
		SearchFlagShowURL:    "Include the product URL in text output",
		SearchFlagCompact:    "Use a narrower product name column in text output",
		SearchFlagWide:       "Use a wider product name column in text output",

		ViewShort:                "View a product in detail",
		ViewExample:              "  pchome view https://24h.pchome.com.tw/prod/DRAA5K-A900JOK9O",
		ViewErrNDJSONUnsupported: "view does not support ndjson; use --format json or text",

		RecommendShort:       "Get product recommendations",
		RecommendExample:     "  pchome recommend DMBL53-A900JDNJS --top 8 --why",
		RecommendFlagTop:     "Maximum number of recommendations",
		RecommendFlagColumns: "Comma-separated text columns (#,score,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url,why)",
		RecommendFlagShowURL: "Include the product URL in text output",
		RecommendFlagShowWhy: "Show recommendation reasoning in text output",
		RecommendFlagCompact: "Use a narrower product name column in text output",
		RecommendFlagWide:    "Use a wider product name column in text output",

		CompareShort:       "Compare multiple products",
		CompareExample:     "  pchome compare DMBL53-A900JDNJS DMBL1C-A900JA04J",
		CompareFlagColumns: "Comma-separated text columns (#,price,list,discount,rating,reviews,24h,qty,stock,brand,name,id,url)",
		CompareFlagShowURL: "Include the product URL in text output",
		CompareFlagCompact: "Use a narrower product name column in text output",
		CompareFlagWide:    "Use a wider product name column in text output",

		SuggestShort:     "Suggest likely search queries",
		SuggestExample:   "  pchome suggest \"掃地機\"",
		SuggestFlagLimit: "Maximum number of suggestions",

		RenderSearchSummary:    "Search %q | returned %d of %d | page %d | scanned %d page(s) | sort=%s\n",
		RenderFilters:          "Filters: %s\n",
		RenderRecommendSummary: "Recommendations for %s | returned %d | took %.0fms\n",
		RenderCompareSummary:   "Compare %d product(s)\n",
		RenderSuggestSummary:   "Suggestions for %q | returned %d\n",
		RenderWarning:          "Warning: %s\n",

		FieldID:             "Id",
		FieldURL:            "URL",
		FieldBrand:          "Brand",
		FieldSalesName:      "Sales Name",
		FieldNick:           "Nick",
		FieldTagline:        "Tagline",
		FieldDescription:    "Description",
		FieldPrice:          "Price",
		FieldListPrice:      "List Price",
		FieldLowestObserved: "Lowest Observed",
		FieldDiscount:       "Discount",
		FieldRating:         "Rating",
		FieldAvailability:   "Availability",
		FieldFlags:          "Flags",
		FieldCategories:     "Categories",
		FieldBrandAliases:   "Brand Aliases",
		FieldPrimaryImage:   "Primary Image",
		FieldImages:         "Images",
		FieldMoreImages:     "More Images",
		FieldNotAvailable:   "n/a",

		AvailabilitySummary: "%s=%s | 24h=%s | %s=%s\n",
		AvailabilityStock:   "stock",
		AvailabilityShip:    "ship",
		FlagsSummary:        "%s=%s | %s=%s\n",
		FlagPrimeOnly:       "prime_only",
		FlagOrderDiscount:   "order_discount",
		RatingWithReviews:   " (%s reviews)",

		FilterCategory:   "category",
		FilterBrand:      "brand",
		FilterMinPrice:   "min-price",
		FilterMaxPrice:   "max-price",
		FilterMinRating:  "min-rating",
		FilterInStock:    "in-stock",
		FilterArrival24h: "arrival-24h",

		HeaderIndex:       "#",
		HeaderPrice:       "Price",
		HeaderList:        "List",
		HeaderDiscount:    "Disc%",
		HeaderRating:      "Rating",
		HeaderReviews:     "Reviews",
		Header24h:         "24h",
		HeaderQty:         "Qty",
		HeaderBrand:       "Brand",
		HeaderName:        "Name",
		HeaderProductID:   "Product ID",
		HeaderURL:         "URL",
		HeaderDescription: "Description",
		HeaderScore:       "Score",
		HeaderWhy:         "Why",
		HeaderSuggestion:  "Suggestion",

		ErrInvalidFormat:        "unsupported format %q (use text, json, or ndjson)",
		ErrInvalidSchemaVersion: "unsupported schema version %q (use %s)",
		ErrUnknownColumn:        "unknown column %q",
	},
}
