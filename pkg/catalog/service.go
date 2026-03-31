package catalog

import (
	"context"
	"crypto/rand"
	"fmt"
	"html"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/oliy/pchome-cli/pkg/pchome/hermes"
	"github.com/oliy/pchome-cli/pkg/pchome/prodapi"
	"github.com/oliy/pchome-cli/pkg/pchome/search"
)

type SearchOptions struct {
	Query      string
	Category   string
	Sort       string
	Page       int
	PageSize   int
	Limit      int
	Brand      string
	MinPrice   int
	MaxPrice   int
	MinRating  float64
	InStock    bool
	Arrival24h bool
}

type RecommendOptions struct {
	Ref   string
	Top   int
	Token string
}

type SuggestOptions struct {
	Query string
	Limit int
}

type Service struct {
	search  *search.Client
	prodapi *prodapi.Client
	hermes  *hermes.Client
}

func New(searchClient *search.Client, prodClient *prodapi.Client, hermesClient *hermes.Client) *Service {
	return &Service{
		search:  searchClient,
		prodapi: prodClient,
		hermes:  hermesClient,
	}
}

type productRecord struct {
	id             string
	searchProduct  *search.ResultProduct
	recommendation *hermes.GoodsRankItem
	product        *prodapi.Product
	stockQty       *int
	inStock        *bool
	primeOnly      *bool
	orderDiscount  *bool
}

func (s *Service) Search(ctx context.Context, opt SearchOptions) (*SearchResult, error) {
	query := strings.TrimSpace(opt.Query)
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	page := opt.Page
	if page <= 0 {
		page = 1
	}

	pageSize := opt.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	limit := opt.Limit
	if limit <= 0 {
		limit = pageSize
	}

	sortLabel, apiSort, err := normalizeSort(opt.Sort)
	if err != nil {
		return nil, err
	}

	fetchSize := pageSize
	if limit > fetchSize {
		fetchSize = limit
	}
	if fetchSize > 50 {
		fetchSize = 50
	}

	result := &SearchResult{
		SchemaVersion: SchemaVersion,
		Query:         query,
		Sort:          sortLabel,
		Page:          page,
		PageSize:      pageSize,
		Limit:         limit,
		Filters: SearchFilters{
			Category:   strings.TrimSpace(opt.Category),
			Brand:      strings.TrimSpace(opt.Brand),
			MinPrice:   opt.MinPrice,
			MaxPrice:   opt.MaxPrice,
			MinRating:  opt.MinRating,
			InStock:    opt.InStock,
			Arrival24h: opt.Arrival24h,
		},
	}

	for len(result.Items) < limit {
		raw, err := s.search.Results(ctx, search.ResultsParams{
			Q:         query,
			CateID:    strings.TrimSpace(opt.Category),
			Page:      page,
			PageCount: fetchSize,
			Sort:      apiSort,
			Price:     formatPriceRange(opt.MinPrice, opt.MaxPrice),
		})
		if err != nil {
			return nil, err
		}

		if result.TotalRows == 0 {
			result.TotalRows = raw.TotalRows
			result.TotalPages = raw.TotalPage
		}
		result.PagesScanned++

		records, warnings := s.enrichRecords(ctx, idsFromSearchProducts(raw.Prods))
		result.Warnings = append(result.Warnings, warnings...)

		for _, rawProd := range raw.Prods {
			rec := ensureProductRecord(records, rawProd.Id)
			copyRaw := rawProd
			rec.searchProduct = &copyRaw

			item := summaryFromRecord(rec)
			if !matchesSearchFilters(item, opt) {
				continue
			}
			result.Items = append(result.Items, item)
			if len(result.Items) >= limit {
				break
			}
		}

		if len(result.Items) >= limit || page >= raw.TotalPage || len(raw.Prods) == 0 {
			break
		}
		page++
	}

	result.Returned = len(result.Items)
	return result, nil
}

func (s *Service) View(ctx context.Context, ref string) (*ProductDetail, error) {
	id, err := NormalizeProductRef(ref)
	if err != nil {
		return nil, err
	}

	records, warnings := s.enrichRecords(ctx, []string{id})
	rec := records[id]
	if rec == nil || rec.product == nil {
		return nil, fmt.Errorf("product %s not found", id)
	}

	detail := detailFromRecord(rec)
	detail.SchemaVersion = SchemaVersion
	detail.Warnings = warnings
	return &detail, nil
}

func (s *Service) Compare(ctx context.Context, refs []string) (*CompareResult, error) {
	ids, err := normalizeRefs(refs)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one product id or URL is required")
	}

	records, warnings := s.enrichRecords(ctx, ids)
	out := &CompareResult{
		SchemaVersion: SchemaVersion,
		Warnings:      warnings,
		Items:         make([]ProductSummary, 0, len(ids)),
	}

	for _, id := range ids {
		rec := records[id]
		if rec == nil || rec.product == nil {
			out.Warnings = append(out.Warnings, fmt.Sprintf("product %s was not found", id))
			continue
		}
		out.Items = append(out.Items, summaryFromRecord(rec))
	}

	if len(out.Items) == 0 {
		return nil, fmt.Errorf("no products were found")
	}

	out.Returned = len(out.Items)
	return out, nil
}

func (s *Service) Recommend(ctx context.Context, opt RecommendOptions) (*RecommendationResult, error) {
	id, err := NormalizeProductRef(opt.Ref)
	if err != nil {
		return nil, err
	}

	top := opt.Top
	if top <= 0 {
		top = 10
	}

	venGuid, err := randomUUID()
	if err != nil {
		return nil, err
	}
	venSession, err := randomUUID()
	if err != nil {
		return nil, err
	}

	raw, err := s.hermes.GoodsRank(ctx, hermes.GoodsRankRequest{
		Token:      hermesToken(opt.Token),
		RecType:    hermes.DefaultRecType,
		RecPos:     hermes.DefaultRecPos,
		TopK:       top,
		Device:     "pc",
		UID:        "0",
		VenGuid:    venGuid,
		VenSession: venSession,
		GID:        id,
	})
	if err != nil {
		return nil, err
	}
	if raw.Error != "" {
		return nil, fmt.Errorf("recommendation API error: %s", raw.Error)
	}

	ids := make([]string, 0, len(raw.RecomdList))
	for _, item := range raw.RecomdList {
		if item.ID != "" {
			ids = append(ids, item.ID)
		}
	}

	records, warnings := s.enrichRecords(ctx, ids)
	out := &RecommendationResult{
		SchemaVersion: SchemaVersion,
		SeedID:        id,
		SeedURL:       prodapi.ProdURL(id),
		TookMs:        raw.Took,
		Warnings:      warnings,
		Items:         make([]Recommendation, 0, len(raw.RecomdList)),
	}

	for _, item := range raw.RecomdList {
		rec := ensureProductRecord(records, item.ID)
		copyItem := item
		rec.recommendation = &copyItem

		reason := copyItem.Msg
		var sourceItem, sourceModel, sourceAlg string
		if len(copyItem.RefItemList) > 0 {
			sourceItem = copyItem.RefItemList[0].ItemName
			sourceModel = copyItem.RefItemList[0].RefModel
			sourceAlg = copyItem.RefItemList[0].AlgName
		}

		out.Items = append(out.Items, Recommendation{
			Product:     summaryFromRecord(rec),
			Score:       copyItem.Score,
			Reason:      reason,
			MessageType: copyItem.MsgType,
			SourceItem:  sourceItem,
			SourceModel: sourceModel,
			SourceAlg:   sourceAlg,
		})
	}

	out.Returned = len(out.Items)
	return out, nil
}

func (s *Service) Suggest(ctx context.Context, opt SuggestOptions) (*SuggestResult, error) {
	query := strings.TrimSpace(opt.Query)
	if query == "" {
		return nil, fmt.Errorf("suggest query is required")
	}

	limit := opt.Limit
	if limit <= 0 {
		limit = 10
	}

	raw, err := s.search.SuggestWords(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]SuggestItem, 0, min(limit, len(raw)))
	for _, item := range raw {
		if len(items) >= limit {
			break
		}
		items = append(items, SuggestItem{
			Name:     item.Name,
			URL:      item.Rurl,
			ImageURL: item.Src,
			Source:   item.Src,
			BU:       item.BU,
		})
	}

	return &SuggestResult{
		SchemaVersion: SchemaVersion,
		Query:         query,
		Returned:      len(items),
		Items:         items,
	}, nil
}

func (s *Service) enrichRecords(ctx context.Context, ids []string) (map[string]*productRecord, []string) {
	ids = dedupeNormalizedIDs(ids)
	records := make(map[string]*productRecord, len(ids))
	for _, id := range ids {
		records[id] = &productRecord{id: id}
	}
	if len(ids) == 0 {
		return records, nil
	}

	fields := []string{
		"Id", "Name", "Nick", "SalesName", "Tagline", "Qty", "ShipDay", "ShipType",
		"Price", "Pic", "PicExtra", "BrandList", "CategoryIds",
		"isArrival24h", "isOnSale", "isOrderDiscount", "isPrimeOnly",
		"RatingValue", "ReviewCount",
	}

	warnings := []string{}

	products, err := s.prodapi.Products(ctx, prodapi.ProductsOptions{
		IDs:    ids,
		Fields: fields,
	})
	if err != nil {
		return records, append(warnings, fmt.Sprintf("product enrichment unavailable: %v", err))
	}

	for key, product := range products {
		id := prodapi.NormalizeID(product.Id)
		if id == "" {
			id = prodapi.NormalizeID(key)
		}
		rec := ensureProductRecord(records, id)
		copyProduct := product
		rec.product = &copyProduct
	}

	buttons, err := s.prodapi.Button(ctx, prodapi.ButtonOptions{
		IDs:    ids,
		Fields: []string{"Id", "Qty", "SaleStatus", "ButtonType", "isPrimeOnly", "isOrderDiscount"},
	})
	if err != nil {
		return records, append(warnings, fmt.Sprintf("stock enrichment unavailable: %v", err))
	}

	seen := map[string]bool{}
	for _, button := range buttons {
		id := prodapi.NormalizeID(button.Id)
		if id == "" {
			continue
		}
		seen[id] = true
		rec := ensureProductRecord(records, id)

		if button.IsPrimeOnly != nil {
			value := *button.IsPrimeOnly == prodapi.Bool1
			rec.primeOnly = &value
		}
		if button.IsOrderDiscount != nil {
			value := *button.IsOrderDiscount == prodapi.Bool1
			rec.orderDiscount = &value
		}
		if button.SaleStatus == 1 && button.Qty != nil {
			if rec.stockQty == nil || *button.Qty > *rec.stockQty {
				value := *button.Qty
				rec.stockQty = &value
			}
		}
	}

	for id := range seen {
		rec := ensureProductRecord(records, id)
		if rec.stockQty == nil {
			value := 0
			rec.stockQty = &value
		}
		if rec.stockQty != nil {
			inStock := *rec.stockQty > 0
			rec.inStock = &inStock
		}
	}

	return records, warnings
}

func summaryFromRecord(rec *productRecord) ProductSummary {
	product := rec.product
	searchProduct := rec.searchProduct
	reco := rec.recommendation

	id := rec.id
	name := ""
	brand := ""
	description := ""
	imageURL := ""
	categoryIDs := []string{}

	if searchProduct != nil {
		if id == "" {
			id = prodapi.NormalizeID(searchProduct.Id)
		}
		name = searchProduct.Name
		brand = strings.TrimSpace(searchProduct.Brand)
		description = strings.TrimSpace(searchProduct.Describe)
		imageURL = fullImageURL(searchProduct.PicB)
		categoryIDs = append(categoryIDs, searchProduct.PCateID...)
	}

	if product != nil {
		if id == "" {
			id = prodapi.NormalizeID(product.Id)
		}
		if product.Name != "" {
			name = product.Name
		}
		if len(product.BrandList) > 0 {
			brand = strings.TrimSpace(product.BrandList[0])
		}
		if description == "" {
			description = strings.TrimSpace(product.Tagline)
		}
		if imageURL == "" && product.Pic != nil {
			imageURL = firstNonEmpty(
				fullImageURL(product.Pic.W),
				fullImageURL(product.Pic.B),
				fullImageURL(product.Pic.S),
			)
		}
		categoryIDs = append(categoryIDs, product.CategoryIds...)
	}

	if reco != nil {
		if id == "" {
			id = prodapi.NormalizeID(reco.ID)
		}
		if name == "" {
			name = html.UnescapeString(reco.Name)
		}
		if imageURL == "" {
			imageURL = strings.TrimSpace(reco.GoodsImgURL)
		}
	}

	price := priceFromRecord(rec)

	var arrival24h *bool
	var shipType string
	var rating Rating
	if product != nil {
		if product.IsArrival24h != nil {
			value := *product.IsArrival24h == prodapi.Bool1
			arrival24h = &value
		}
		shipType = product.ShipType
		rating = mergeRating(product.RatingValue, product.ReviewCount, nil, nil)
	}
	if searchProduct != nil {
		rating = mergeRating(rating.Score, rating.Reviews, searchProduct.RatingValue, searchProduct.ReviewCount)
	}

	availability := Availability{
		InStock:       rec.inStock,
		StockQty:      rec.stockQty,
		Arrival24h:    arrival24h,
		ShipType:      shipType,
		PrimeOnly:     rec.primeOnly,
		OrderDiscount: rec.orderDiscount,
	}
	if availability.InStock == nil && product != nil && product.Qty != nil {
		value := *product.Qty > 0
		availability.InStock = &value
		availability.StockQty = product.Qty
	}

	return ProductSummary{
		ID:           id,
		Name:         name,
		Brand:        brand,
		Description:  description,
		Price:        price,
		Availability: availability,
		Rating:       rating,
		CategoryIDs:  dedupeStrings(categoryIDs),
		URL:          prodapi.ProdURL(id),
		ImageURL:     imageURL,
	}
}

func detailFromRecord(rec *productRecord) ProductDetail {
	summary := summaryFromRecord(rec)
	product := rec.product
	if product == nil {
		return ProductDetail{
			SchemaVersion: SchemaVersion,
			Product:       summary,
		}
	}

	images := []string{}
	if product.Pic != nil {
		images = append(images,
			fullImageURL(product.Pic.W),
			fullImageURL(product.Pic.B),
			fullImageURL(product.Pic.S),
		)
	}
	for _, image := range product.PicExtra {
		images = append(images, fullImageURL(image))
	}

	return ProductDetail{
		SchemaVersion: SchemaVersion,
		Product:       summary,
		Nick:          product.Nick,
		SalesName:     product.SalesName,
		Tagline:       product.Tagline,
		BrandAliases:  dedupeStrings(product.BrandList),
		Images:        dedupeStrings(images),
	}
}

func priceFromRecord(rec *productRecord) Money {
	current := 0
	original := 0
	lowest := 0

	if rec.searchProduct != nil {
		current = rec.searchProduct.Price
		original = rec.searchProduct.OriginPrice
	}

	if rec.product != nil && rec.product.Price != nil {
		if rec.product.Price.P > 0 {
			current = rec.product.Price.P
		}
		if rec.product.Price.Low != nil && *rec.product.Price.Low > 0 {
			lowest = *rec.product.Price.Low
			current = *rec.product.Price.Low
		}
		if rec.product.Price.M > 0 {
			original = rec.product.Price.M
		}
	}

	if rec.recommendation != nil && current == 0 {
		current = int(math.Round(rec.recommendation.SalePrice))
	}

	if original > 0 && original < current {
		original = 0
	}

	money := Money{
		Current:  current,
		Original: original,
		Lowest:   lowest,
	}
	if original > 0 && current > 0 && original > current {
		money.DiscountPercent = int(math.Round((1 - float64(current)/float64(original)) * 100))
	}
	return money
}

func mergeRating(baseScore *float64, baseReviews *int, nextScore *float64, nextReviews *int) Rating {
	out := Rating{
		Score:   baseScore,
		Reviews: baseReviews,
	}
	if out.Score == nil && nextScore != nil {
		value := *nextScore
		out.Score = &value
	}
	if out.Reviews == nil && nextReviews != nil {
		value := *nextReviews
		out.Reviews = &value
	}
	return out
}

func matchesSearchFilters(item ProductSummary, opt SearchOptions) bool {
	if opt.Brand != "" {
		needle := strings.ToLower(strings.TrimSpace(opt.Brand))
		haystacks := []string{strings.ToLower(item.Brand), strings.ToLower(item.Name)}
		found := false
		for _, haystack := range haystacks {
			if strings.Contains(haystack, needle) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if opt.MinRating > 0 {
		if item.Rating.Score == nil || *item.Rating.Score < opt.MinRating {
			return false
		}
	}

	if opt.InStock {
		if item.Availability.InStock == nil || !*item.Availability.InStock {
			return false
		}
	}

	if opt.Arrival24h {
		if item.Availability.Arrival24h == nil || !*item.Availability.Arrival24h {
			return false
		}
	}

	return true
}

func idsFromSearchProducts(products []search.ResultProduct) []string {
	ids := make([]string, 0, len(products))
	for _, product := range products {
		if product.Id != "" {
			ids = append(ids, product.Id)
		}
	}
	return ids
}

func normalizeRefs(refs []string) ([]string, error) {
	out := make([]string, 0, len(refs))
	seen := map[string]bool{}
	for _, ref := range refs {
		id, err := NormalizeProductRef(ref)
		if err != nil {
			return nil, err
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out, nil
}

func dedupeNormalizedIDs(ids []string) []string {
	out := make([]string, 0, len(ids))
	seen := map[string]bool{}
	for _, id := range ids {
		id = prodapi.NormalizeID(strings.TrimSpace(id))
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func dedupeStrings(values []string) []string {
	out := make([]string, 0, len(values))
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func ensureProductRecord(records map[string]*productRecord, id string) *productRecord {
	id = prodapi.NormalizeID(id)
	if rec, ok := records[id]; ok {
		return rec
	}
	rec := &productRecord{id: id}
	records[id] = rec
	return rec
}

func fullImageURL(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if strings.Contains(path, "://") {
		return path
	}
	if strings.HasPrefix(path, "/") {
		return "https://cs-c.ecimg.tw" + path
	}
	return path
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeSort(value string) (string, string, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", "relevance":
		return "relevance", "", nil
	case "price-asc":
		return "price-asc", "prc/ac", nil
	case "price-desc":
		return "price-desc", "prc/dc", nil
	case "newest":
		return "newest", "new/dc", nil
	case "best-selling":
		return "best-selling", "sale/dc", nil
	default:
		return "", "", fmt.Errorf("unsupported sort %q (use relevance, price-asc, price-desc, newest, best-selling)", value)
	}
}

func formatPriceRange(minPrice, maxPrice int) string {
	if minPrice <= 0 && maxPrice <= 0 {
		return ""
	}
	if minPrice < 0 {
		minPrice = 0
	}
	if maxPrice <= 0 {
		maxPrice = 99999999
	}
	return fmt.Sprintf("%d-%d", minPrice, maxPrice)
}

func randomUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5],
		b[6], b[7],
		b[8], b[9],
		b[10], b[11], b[12], b[13], b[14], b[15],
	), nil
}

func hermesToken(configToken string) string {
	if token := strings.TrimSpace(configToken); token != "" {
		return token
	}
	if token := strings.TrimSpace(os.Getenv("PCHOME_HERMES_TOKEN")); token != "" {
		return token
	}
	return hermes.DefaultToken
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sortSummariesByPriceThenName(items []ProductSummary) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Price.Current == items[j].Price.Current {
			return items[i].Name < items[j].Name
		}
		return items[i].Price.Current < items[j].Price.Current
	})
}
