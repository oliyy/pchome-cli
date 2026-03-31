package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/oliy/pchome-cli/pkg/output"
)

func renderSearchResult(w io.Writer, loc *i18n.Catalog, result *catalog.SearchResult, columnsCSV string, showURL bool, nameWidth int) error {
	columns, err := selectColumns(
		parseCSV(columnsCSV),
		searchColumnSpecs(loc, nameWidth),
		defaultSearchColumns(showURL),
		loc,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		w,
		loc.Text(i18n.RenderSearchSummary),
		result.Query,
		result.Returned,
		result.TotalRows,
		result.Page,
		result.PagesScanned,
		result.Sort,
	)
	if filters := formatSearchFilters(loc, result.Filters); filters != "" {
		fmt.Fprintf(w, loc.Text(i18n.RenderFilters), filters)
	}

	rows := [][]string{columnHeaders(columns)}
	for index, item := range result.Items {
		rows = append(rows, columnValues(columns, summaryColumnContext{
			index: index + 1,
			item:  item,
		}))
	}

	output.PrintTable(w, rows)
	renderWarnings(w, loc, result.Warnings)
	return nil
}

func renderRecommendResult(w io.Writer, loc *i18n.Catalog, result *catalog.RecommendationResult, columnsCSV string, showURL, showWhy bool, nameWidth int) error {
	defaultColumns := defaultRecommendColumns(showURL, showWhy)
	columns, err := selectColumns(
		parseCSV(columnsCSV),
		recommendColumnSpecs(loc, nameWidth),
		defaultColumns,
		loc,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		w,
		loc.Text(i18n.RenderRecommendSummary),
		result.SeedID,
		result.Returned,
		result.TookMs,
	)

	rows := [][]string{columnHeaders(columns)}
	for index, item := range result.Items {
		rows = append(rows, columnValues(columns, summaryColumnContext{
			index:         index + 1,
			item:          item.Product,
			score:         item.Score,
			reason:        joinNonEmpty(" | ", item.Reason, item.SourceModel),
			includeScore:  true,
			includeReason: true,
		}))
	}

	output.PrintTable(w, rows)
	renderWarnings(w, loc, result.Warnings)
	return nil
}

func renderCompareResult(w io.Writer, loc *i18n.Catalog, result *catalog.CompareResult, columnsCSV string, showURL bool, nameWidth int) error {
	columns, err := selectColumns(
		parseCSV(columnsCSV),
		searchColumnSpecs(loc, nameWidth),
		defaultCompareColumns(showURL),
		loc,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, loc.Text(i18n.RenderCompareSummary), result.Returned)

	rows := [][]string{columnHeaders(columns)}
	for index, item := range result.Items {
		rows = append(rows, columnValues(columns, summaryColumnContext{
			index: index + 1,
			item:  item,
		}))
	}

	output.PrintTable(w, rows)
	renderWarnings(w, loc, result.Warnings)
	return nil
}

func renderProductDetail(w io.Writer, loc *i18n.Catalog, detail *catalog.ProductDetail) {
	product := detail.Product
	notAvailable := loc.Text(i18n.FieldNotAvailable)

	fmt.Fprintln(w, product.Name)
	fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldID), product.ID)
	fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldURL), product.URL)
	if product.Brand != "" {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldBrand), product.Brand)
	}
	if detail.SalesName != "" && detail.SalesName != product.Name {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldSalesName), detail.SalesName)
	}
	if detail.Nick != "" && detail.Nick != product.Name {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldNick), detail.Nick)
	}
	if detail.Tagline != "" {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldTagline), detail.Tagline)
	}
	if product.Description != "" {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldDescription), product.Description)
	}

	fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldPrice), formatMoney(product.Price.Current))
	if product.Price.Original > 0 {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldListPrice), formatMoney(product.Price.Original))
	}
	if product.Price.Lowest > 0 && product.Price.Lowest != product.Price.Current {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldLowestObserved), formatMoney(product.Price.Lowest))
	}
	if product.Price.DiscountPercent > 0 {
		fmt.Fprintf(w, "%s: %d%%\n", loc.Text(i18n.FieldDiscount), product.Price.DiscountPercent)
	}

	if product.Rating.Score != nil || product.Rating.Reviews != nil {
		fmt.Fprintf(
			w,
			"%s: %s",
			loc.Text(i18n.FieldRating),
			emptyFallback(formatRating(product.Rating.Score), notAvailable),
		)
		if product.Rating.Reviews != nil {
			fmt.Fprintf(w, loc.Text(i18n.RatingWithReviews), formatInt(product.Rating.Reviews))
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "%s: ", loc.Text(i18n.FieldAvailability))
	fmt.Fprintf(
		w,
		loc.Text(i18n.AvailabilitySummary),
		loc.Text(i18n.AvailabilityStock),
		emptyFallback(formatInt(product.Availability.StockQty), "?"),
		formatBoolMarker(product.Availability.Arrival24h),
		loc.Text(i18n.AvailabilityShip),
		emptyFallback(product.Availability.ShipType, "?"),
	)
	if product.Availability.PrimeOnly != nil || product.Availability.OrderDiscount != nil {
		fmt.Fprintf(w, "%s: ", loc.Text(i18n.FieldFlags))
		fmt.Fprintf(
			w,
			loc.Text(i18n.FlagsSummary),
			loc.Text(i18n.FlagPrimeOnly),
			formatBoolMarker(product.Availability.PrimeOnly),
			loc.Text(i18n.FlagOrderDiscount),
			formatBoolMarker(product.Availability.OrderDiscount),
		)
	}
	if len(product.CategoryIDs) > 0 {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldCategories), strings.Join(product.CategoryIDs, ", "))
	}
	if len(detail.BrandAliases) > 1 {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldBrandAliases), strings.Join(detail.BrandAliases, ", "))
	}
	if product.ImageURL != "" {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldPrimaryImage), product.ImageURL)
	}
	if len(detail.Images) > 1 {
		fmt.Fprintf(w, "%s: %s\n", loc.Text(i18n.FieldImages), strings.Join(detail.Images[:min(3, len(detail.Images))], ", "))
		if len(detail.Images) > 3 {
			fmt.Fprintf(w, "%s: %d\n", loc.Text(i18n.FieldMoreImages), len(detail.Images)-3)
		}
	}

	renderWarnings(w, loc, detail.Warnings)
}

func renderSuggestResult(w io.Writer, loc *i18n.Catalog, result *catalog.SuggestResult) {
	fmt.Fprintf(w, loc.Text(i18n.RenderSuggestSummary), result.Query, result.Returned)
	rows := [][]string{{loc.Text(i18n.HeaderIndex), loc.Text(i18n.HeaderSuggestion), loc.Text(i18n.HeaderURL)}}
	for index, item := range result.Items {
		rows = append(rows, []string{
			fmt.Sprintf("%d", index+1),
			item.Name,
			item.URL,
		})
	}
	output.PrintTable(w, rows)
}

type summaryColumnContext struct {
	index         int
	item          catalog.ProductSummary
	score         float64
	reason        string
	includeScore  bool
	includeReason bool
}

type columnSpec struct {
	header string
	value  func(summaryColumnContext) string
}

func searchColumnSpecs(loc *i18n.Catalog, nameWidth int) map[string]columnSpec {
	specs := map[string]columnSpec{
		"#":        {header: loc.Text(i18n.HeaderIndex), value: func(ctx summaryColumnContext) string { return zeroAwareIndex(ctx.index) }},
		"price":    {header: loc.Text(i18n.HeaderPrice), value: func(ctx summaryColumnContext) string { return formatMoney(ctx.item.Price.Current) }},
		"list":     {header: loc.Text(i18n.HeaderList), value: func(ctx summaryColumnContext) string { return zeroAwareMoney(ctx.item.Price.Original) }},
		"discount": {header: loc.Text(i18n.HeaderDiscount), value: func(ctx summaryColumnContext) string { return zeroAwarePercent(ctx.item.Price.DiscountPercent) }},
		"rating": {header: loc.Text(i18n.HeaderRating), value: func(ctx summaryColumnContext) string {
			return formatRating(ctx.item.Rating.Score)
		}},
		"reviews": {header: loc.Text(i18n.HeaderReviews), value: func(ctx summaryColumnContext) string { return formatInt(ctx.item.Rating.Reviews) }},
		"24h":     {header: loc.Text(i18n.Header24h), value: func(ctx summaryColumnContext) string { return formatBoolMarker(ctx.item.Availability.Arrival24h) }},
		"stock": {header: loc.Text(i18n.HeaderQty), value: func(ctx summaryColumnContext) string {
			return emptyFallback(formatInt(ctx.item.Availability.StockQty), "?")
		}},
		"qty": {header: loc.Text(i18n.HeaderQty), value: func(ctx summaryColumnContext) string {
			return emptyFallback(formatInt(ctx.item.Availability.StockQty), "?")
		}},
		"brand": {header: loc.Text(i18n.HeaderBrand), value: func(ctx summaryColumnContext) string { return output.TruncateEnd(ctx.item.Brand, 18) }},
		"name":  {header: loc.Text(i18n.HeaderName), value: func(ctx summaryColumnContext) string { return output.TruncateEnd(ctx.item.Name, nameWidth) }},
		"id":    {header: loc.Text(i18n.HeaderProductID), value: func(ctx summaryColumnContext) string { return ctx.item.ID }},
		"url":   {header: loc.Text(i18n.HeaderURL), value: func(ctx summaryColumnContext) string { return ctx.item.URL }},
		"desc":  {header: loc.Text(i18n.HeaderDescription), value: func(ctx summaryColumnContext) string { return output.TruncateEnd(ctx.item.Description, 48) }},
	}
	return specs
}

func defaultSearchColumns(showURL bool) []string {
	columns := []string{"#", "name", "price", "rating", "reviews", "24h", "brand", "qty"}
	if showURL {
		columns = append(columns, "url")
	}
	return columns
}

func defaultRecommendColumns(showURL, showWhy bool) []string {
	columns := defaultSearchColumns(showURL)
	if showWhy {
		columns = append(columns, "why")
	}
	return columns
}

func defaultCompareColumns(showURL bool) []string {
	return defaultSearchColumns(showURL)
}

func recommendColumnSpecs(loc *i18n.Catalog, nameWidth int) map[string]columnSpec {
	specs := searchColumnSpecs(loc, nameWidth)
	specs["score"] = columnSpec{
		header: loc.Text(i18n.HeaderScore),
		value: func(ctx summaryColumnContext) string {
			if !ctx.includeScore {
				return ""
			}
			return fmt.Sprintf("%.1f", ctx.score)
		},
	}
	specs["why"] = columnSpec{
		header: loc.Text(i18n.HeaderWhy),
		value: func(ctx summaryColumnContext) string {
			if !ctx.includeReason {
				return ""
			}
			return output.TruncateEnd(ctx.reason, 28)
		},
	}
	return specs
}

func selectColumns(requested []string, specs map[string]columnSpec, defaults []string, loc *i18n.Catalog) ([]columnSpec, error) {
	if len(requested) == 0 {
		requested = defaults
	}

	columns := make([]columnSpec, 0, len(requested))
	for _, key := range requested {
		spec, ok := specs[strings.ToLower(key)]
		if !ok {
			return nil, fmt.Errorf(loc.Text(i18n.ErrUnknownColumn), key)
		}
		columns = append(columns, spec)
	}
	return columns, nil
}

func columnHeaders(columns []columnSpec) []string {
	headers := make([]string, 0, len(columns))
	for _, column := range columns {
		headers = append(headers, column.header)
	}
	return headers
}

func columnValues(columns []columnSpec, ctx summaryColumnContext) []string {
	values := make([]string, 0, len(columns))
	for _, column := range columns {
		values = append(values, column.value(ctx))
	}
	return values
}

func formatSearchFilters(loc *i18n.Catalog, filters catalog.SearchFilters) string {
	parts := []string{}
	if filters.Category != "" {
		parts = append(parts, loc.Text(i18n.FilterCategory)+"="+filters.Category)
	}
	if filters.Brand != "" {
		parts = append(parts, loc.Text(i18n.FilterBrand)+"="+filters.Brand)
	}
	if filters.MinPrice > 0 {
		parts = append(parts, loc.Text(i18n.FilterMinPrice)+"="+formatMoney(filters.MinPrice))
	}
	if filters.MaxPrice > 0 {
		parts = append(parts, loc.Text(i18n.FilterMaxPrice)+"="+formatMoney(filters.MaxPrice))
	}
	if filters.MinRating > 0 {
		parts = append(parts, fmt.Sprintf("%s=%.1f", loc.Text(i18n.FilterMinRating), filters.MinRating))
	}
	if filters.InStock {
		parts = append(parts, loc.Text(i18n.FilterInStock))
	}
	if filters.Arrival24h {
		parts = append(parts, loc.Text(i18n.FilterArrival24h))
	}
	return strings.Join(parts, " | ")
}

func renderWarnings(w io.Writer, loc *i18n.Catalog, warnings []string) {
	for _, warning := range warnings {
		if strings.TrimSpace(warning) == "" {
			continue
		}
		fmt.Fprintf(w, loc.Text(i18n.RenderWarning), warning)
	}
}

func zeroAwareIndex(value int) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func zeroAwareMoney(value int) string {
	if value == 0 {
		return ""
	}
	return formatMoney(value)
}

func zeroAwarePercent(value int) string {
	if value == 0 {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

func emptyFallback(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func joinNonEmpty(sep string, values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		parts = append(parts, value)
	}
	return strings.Join(parts, sep)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
