package cmd

import (
	"errors"

	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func newViewCmd(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view PRODUCT",
		Short:   app.I18N.Text(i18n.ViewShort),
		GroupID: "shopping",
		Args:    cobra.ExactArgs(1),
		Example: app.I18N.Text(i18n.ViewExample),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(cmd, app.Timeout)
			defer cancel()

			detail, err := app.Catalog.View(ctx, args[0])
			if err != nil {
				return err
			}

			switch app.Format {
			case "json":
				return printJSON(cmd, detail)
			case "ndjson":
				return errors.New(app.I18N.Text(i18n.ViewErrNDJSONUnsupported))
			default:
				renderProductDetail(cmd.OutOrStdout(), app.I18N, detail)
				return nil
			}
		},
	}

	return cmd
}
