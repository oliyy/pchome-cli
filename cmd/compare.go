package cmd

import (
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func newCompareCmd(app *App) *cobra.Command {
	var (
		columns string
		showURL bool
		compact bool
		wide    bool
	)

	cmd := &cobra.Command{
		Use:     "compare PRODUCT [PRODUCT...]",
		Short:   app.I18N.Text(i18n.CompareShort),
		GroupID: "shopping",
		Args:    cobra.MinimumNArgs(2),
		Example: app.I18N.Text(i18n.CompareExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(cmd, app.Timeout)
			defer cancel()

			result, err := app.Catalog.Compare(ctx, args)
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
				return renderCompareResult(
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

	cmd.Flags().StringVar(&columns, "columns", columnsFlagDefault(app.Config.Compare.Columns), app.I18N.Text(i18n.CompareFlagColumns))
	cmd.Flags().BoolVar(&showURL, "show-url", app.Config.Compare.ShowURL, app.I18N.Text(i18n.CompareFlagShowURL))
	cmd.Flags().BoolVar(&compact, "compact", app.Config.Compare.Compact, app.I18N.Text(i18n.CompareFlagCompact))
	cmd.Flags().BoolVar(&wide, "wide", app.Config.Compare.Wide, app.I18N.Text(i18n.CompareFlagWide))

	return cmd
}
