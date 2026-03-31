package cmd

import (
	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func newSuggestCmd(app *App) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:     "suggest QUERY",
		Short:   app.I18N.Text(i18n.SuggestShort),
		GroupID: "discovery",
		Args:    cobra.MinimumNArgs(1),
		Example: app.I18N.Text(i18n.SuggestExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(cmd, app.Timeout)
			defer cancel()

			result, err := app.Catalog.Suggest(ctx, catalog.SuggestOptions{
				Query: joinArgs(args),
				Limit: limit,
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
				renderSuggestResult(cmd.OutOrStdout(), app.I18N, result)
				return nil
			}
		},
	}

	cmd.Flags().IntVar(&limit, "limit", app.Config.Suggest.Limit, app.I18N.Text(i18n.SuggestFlagLimit))
	return cmd
}
