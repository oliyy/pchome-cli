package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oliy/pchome-cli/pkg/catalog"
	"github.com/oliy/pchome-cli/pkg/config"
	"github.com/oliy/pchome-cli/pkg/httpx"
	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/oliy/pchome-cli/pkg/pchome/hermes"
	"github.com/oliy/pchome-cli/pkg/pchome/prodapi"
	"github.com/oliy/pchome-cli/pkg/pchome/search"
	"github.com/spf13/cobra"
)

var version = "dev"

type App struct {
	Format        string
	SchemaVersion string
	Timeout       time.Duration
	NameWidth     int
	Language      i18n.Language
	HermesToken   string
	Config        config.Config
	ConfigPath    string
	I18N          *i18n.Catalog

	Catalog *catalog.Service
}

func Execute() {
	cmd, err := newRootCmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newRootCmd() (*cobra.Command, error) {
	cfg, path, err := config.Load()
	if err != nil {
		return nil, err
	}

	timeout, err := cfg.TimeoutDuration()
	if err != nil {
		return nil, err
	}

	language, err := i18n.ParseLanguage(cfg.I18N.Language)
	if err != nil {
		return nil, err
	}
	loc := i18n.New(language)

	app := &App{
		Format:        cfg.Output.Format,
		SchemaVersion: cfg.Output.SchemaVersion,
		Timeout:       timeout,
		NameWidth:     cfg.Output.NameWidth,
		Language:      language,
		HermesToken:   cfg.Hermes.Token,
		Config:        cfg,
		ConfigPath:    path,
		I18N:          loc,
	}

	cmd := &cobra.Command{
		Use:           "pchome",
		Short:         loc.Text(i18n.RootShort),
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		Long:          loc.Text(i18n.RootLong),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(app.Format, app.I18N); err != nil {
				return err
			}
			if err := validateSchemaVersion(app.SchemaVersion, app.I18N); err != nil {
				return err
			}
			if app.NameWidth <= 0 {
				app.NameWidth = cfg.Output.NameWidth
			}

			httpClient := httpx.New(httpx.Options{
				Timeout: app.Timeout,
			})
			app.Catalog = catalog.New(
				search.New(httpClient),
				prodapi.New(httpClient),
				hermes.New(httpClient),
			)
			return nil
		},
	}

	cmd.SetContext(context.Background())
	cmd.SetVersionTemplate(fmt.Sprintf("pchome %s\n", version))

	cmd.PersistentFlags().StringVar(&app.Format, "format", app.Format, app.I18N.Text(i18n.RootFlagFormat))
	cmd.PersistentFlags().StringVar(&app.SchemaVersion, "schema-version", app.SchemaVersion, app.I18N.Text(i18n.RootFlagSchemaVersion))
	cmd.PersistentFlags().DurationVar(&app.Timeout, "timeout", app.Timeout, app.I18N.Text(i18n.RootFlagTimeout))
	cmd.PersistentFlags().IntVar(&app.NameWidth, "name-width", app.NameWidth, app.I18N.Text(i18n.RootFlagNameWidth))

	cmd.AddGroup(
		&cobra.Group{ID: "shopping", Title: app.I18N.Text(i18n.GroupShopping)},
		&cobra.Group{ID: "discovery", Title: app.I18N.Text(i18n.GroupDiscovery)},
	)

	cmd.AddCommand(newSearchCmd(app))
	cmd.AddCommand(newViewCmd(app))
	cmd.AddCommand(newRecommendCmd(app))
	cmd.AddCommand(newCompareCmd(app))
	cmd.AddCommand(newSuggestCmd(app))
	localizeRootCommand(cmd, app.I18N)

	return cmd, nil
}

func validateFormat(format string, loc *i18n.Catalog) error {
	switch format {
	case "text", "json", "ndjson":
		return nil
	default:
		return fmt.Errorf(loc.Text(i18n.ErrInvalidFormat), format)
	}
}

func validateSchemaVersion(version string, loc *i18n.Catalog) error {
	if version != catalog.SchemaVersion {
		return fmt.Errorf(loc.Text(i18n.ErrInvalidSchemaVersion), version, catalog.SchemaVersion)
	}
	return nil
}
