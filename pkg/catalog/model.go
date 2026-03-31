package catalog

const SchemaVersion = "v1"

type Money struct {
	Current         int `json:"current"`
	Original        int `json:"original,omitempty"`
	Lowest          int `json:"lowest,omitempty"`
	DiscountPercent int `json:"discount_percent,omitempty"`
}

type Availability struct {
	InStock       *bool  `json:"in_stock,omitempty"`
	StockQty      *int   `json:"stock_qty,omitempty"`
	Arrival24h    *bool  `json:"arrival_24h,omitempty"`
	ShipType      string `json:"ship_type,omitempty"`
	PrimeOnly     *bool  `json:"prime_only,omitempty"`
	OrderDiscount *bool  `json:"order_discount,omitempty"`
}

type Rating struct {
	Score   *float64 `json:"score,omitempty"`
	Reviews *int     `json:"reviews,omitempty"`
}

type ProductSummary struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Brand        string       `json:"brand,omitempty"`
	Description  string       `json:"description,omitempty"`
	Price        Money        `json:"price"`
	Availability Availability `json:"availability"`
	Rating       Rating       `json:"rating"`
	CategoryIDs  []string     `json:"category_ids,omitempty"`
	URL          string       `json:"url"`
	ImageURL     string       `json:"image_url,omitempty"`
}

type ProductDetail struct {
	SchemaVersion string         `json:"schema_version"`
	Product       ProductSummary `json:"product"`
	Nick          string         `json:"nick,omitempty"`
	SalesName     string         `json:"sales_name,omitempty"`
	Tagline       string         `json:"tagline,omitempty"`
	BrandAliases  []string       `json:"brand_aliases,omitempty"`
	Images        []string       `json:"images,omitempty"`
	Warnings      []string       `json:"warnings,omitempty"`
}

type SearchFilters struct {
	Category   string  `json:"category,omitempty"`
	Brand      string  `json:"brand,omitempty"`
	MinPrice   int     `json:"min_price,omitempty"`
	MaxPrice   int     `json:"max_price,omitempty"`
	MinRating  float64 `json:"min_rating,omitempty"`
	InStock    bool    `json:"in_stock,omitempty"`
	Arrival24h bool    `json:"arrival_24h,omitempty"`
}

type SearchResult struct {
	SchemaVersion string           `json:"schema_version"`
	Query         string           `json:"query"`
	Sort          string           `json:"sort"`
	Page          int              `json:"page"`
	PageSize      int              `json:"page_size"`
	Limit         int              `json:"limit"`
	PagesScanned  int              `json:"pages_scanned"`
	TotalRows     int              `json:"total_rows"`
	TotalPages    int              `json:"total_pages"`
	Returned      int              `json:"returned"`
	Filters       SearchFilters    `json:"filters"`
	Warnings      []string         `json:"warnings,omitempty"`
	Items         []ProductSummary `json:"items"`
}

type Recommendation struct {
	Product     ProductSummary `json:"product"`
	Score       float64        `json:"score"`
	Reason      string         `json:"reason,omitempty"`
	MessageType string         `json:"message_type,omitempty"`
	SourceItem  string         `json:"source_item,omitempty"`
	SourceModel string         `json:"source_model,omitempty"`
	SourceAlg   string         `json:"source_algorithm,omitempty"`
}

type RecommendationResult struct {
	SchemaVersion string           `json:"schema_version"`
	SeedID        string           `json:"seed_id"`
	SeedURL       string           `json:"seed_url"`
	TookMs        float64          `json:"took_ms"`
	Returned      int              `json:"returned"`
	Warnings      []string         `json:"warnings,omitempty"`
	Items         []Recommendation `json:"items"`
}

type CompareResult struct {
	SchemaVersion string           `json:"schema_version"`
	Returned      int              `json:"returned"`
	Warnings      []string         `json:"warnings,omitempty"`
	Items         []ProductSummary `json:"items"`
}

type SuggestItem struct {
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	Source   string `json:"source,omitempty"`
	BU       string `json:"bu,omitempty"`
}

type SuggestResult struct {
	SchemaVersion string        `json:"schema_version"`
	Query         string        `json:"query"`
	Returned      int           `json:"returned"`
	Items         []SuggestItem `json:"items"`
}
