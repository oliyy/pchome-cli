package cmd

import (
	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func newSearchCmd(app *App) *cobra.Command {
	var (
		category   string
		brand      string
		sort       string
		page       int
		pageSize   int
		limit      int
		minPrice   int
		maxPrice   int
		minRating  float64
		inStock    bool
		arrival24h bool
		columns    string
		showURL    bool
		compact    bool
		wide       bool
	)

	cmd := &cobra.Command{
		Use:     "search QUERY",
		Short:   app.I18N.Text(i18n.SearchShort),
		GroupID: "shopping",
		Args:    cobra.MinimumNArgs(1),
		Example: app.I18N.Text(i18n.SearchExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(cmd, app.Timeout)
			defer cancel()

			result, err := app.Catalog.Search(ctx, catalog.SearchOptions{
				Query:      joinArgs(args),
				Category:   category,
				Brand:      brand,
				Sort:       sort,
				Page:       page,
				PageSize:   pageSize,
				Limit:      limit,
				MinPrice:   minPrice,
				MaxPrice:   maxPrice,
				MinRating:  minRating,
				InStock:    inStock,
				Arrival24h: arrival24h,
			})
			if err != nil {
				return err
			}

			switch app.Format {
			case "json":
				return printJSON(cmd, result)
			case "ndjson":
				items := make([]any, 0, len(result.Items))
				for _, item := range result.Items {
					items = append(items, item)
				}
				return printNDJSON(cmd, items)
			default:
				return renderSearchResult(
					cmd.OutOrStdout(),
					app.I18N,
					result,
					columns,
					showURL,
					effectiveNameWidth(app.NameWidth, compact, wide),
				)
			}
		},
	}

	cmd.Flags().StringVar(&category, "category", app.Config.Search.Category, app.I18N.Text(i18n.SearchFlagCategory))
	cmd.Flags().StringVar(&brand, "brand", app.Config.Search.Brand, app.I18N.Text(i18n.SearchFlagBrand))
	cmd.Flags().StringVar(&sort, "sort", app.Config.Search.Sort, app.I18N.Text(i18n.SearchFlagSort))
	cmd.Flags().IntVar(&page, "page", 1, app.I18N.Text(i18n.SearchFlagPage))
	cmd.Flags().IntVar(&pageSize, "page-size", app.Config.Search.PageSize, app.I18N.Text(i18n.SearchFlagPageSize))
	cmd.Flags().IntVar(&limit, "limit", app.Config.Search.Limit, app.I18N.Text(i18n.SearchFlagLimit))
	cmd.Flags().IntVar(&minPrice, "min-price", app.Config.Search.MinPrice, app.I18N.Text(i18n.SearchFlagMinPrice))
	cmd.Flags().IntVar(&maxPrice, "max-price", app.Config.Search.MaxPrice, app.I18N.Text(i18n.SearchFlagMaxPrice))
	cmd.Flags().Float64Var(&minRating, "min-rating", app.Config.Search.MinRating, app.I18N.Text(i18n.SearchFlagMinRating))
	cmd.Flags().BoolVar(&inStock, "in-stock", app.Config.Search.InStock, app.I18N.Text(i18n.SearchFlagInStock))
	cmd.Flags().BoolVar(&arrival24h, "arrival-24h", app.Config.Search.Arrival24h, app.I18N.Text(i18n.SearchFlagArrival24h))
	cmd.Flags().StringVar(&columns, "columns", columnsFlagDefault(app.Config.Search.Columns), app.I18N.Text(i18n.SearchFlagColumns))
	cmd.Flags().BoolVar(&showURL, "show-url", app.Config.Search.ShowURL, app.I18N.Text(i18n.SearchFlagShowURL))
	cmd.Flags().BoolVar(&compact, "compact", app.Config.Search.Compact, app.I18N.Text(i18n.SearchFlagCompact))
	cmd.Flags().BoolVar(&wide, "wide", app.Config.Search.Wide, app.I18N.Text(i18n.SearchFlagWide))

	return cmd
}
