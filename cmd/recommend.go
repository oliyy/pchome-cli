package cmd

import (
	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func newRecommendCmd(app *App) *cobra.Command {
	var (
		top     int
		columns string
		showURL bool
		showWhy bool
		compact bool
		wide    bool
	)

	cmd := &cobra.Command{
		Use:     "recommend PRODUCT",
		Short:   app.I18N.Text(i18n.RecommendShort),
		GroupID: "shopping",
		Args:    cobra.ExactArgs(1),
		Example: app.I18N.Text(i18n.RecommendExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(cmd, app.Timeout)
			defer cancel()

			result, err := app.Catalog.Recommend(ctx, catalog.RecommendOptions{
				Ref:   args[0],
				Top:   top,
				Token: app.HermesToken,
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
				return renderRecommendResult(
					cmd.OutOrStdout(),
					app.I18N,
					result,
					columns,
					showURL,
					showWhy,
					effectiveNameWidth(app.NameWidth, compact, wide),
				)
			}
		},
	}

	cmd.Flags().IntVar(&top, "top", app.Config.Recommend.Top, app.I18N.Text(i18n.RecommendFlagTop))
	cmd.Flags().StringVar(&columns, "columns", columnsFlagDefault(app.Config.Recommend.Columns), app.I18N.Text(i18n.RecommendFlagColumns))
	cmd.Flags().BoolVar(&showURL, "show-url", app.Config.Recommend.ShowURL, app.I18N.Text(i18n.RecommendFlagShowURL))
	cmd.Flags().BoolVar(&showWhy, "why", app.Config.Recommend.ShowWhy, app.I18N.Text(i18n.RecommendFlagShowWhy))
	cmd.Flags().BoolVar(&compact, "compact", app.Config.Recommend.Compact, app.I18N.Text(i18n.RecommendFlagCompact))
	cmd.Flags().BoolVar(&wide, "wide", app.Config.Recommend.Wide, app.I18N.Text(i18n.RecommendFlagWide))

	return cmd
}
