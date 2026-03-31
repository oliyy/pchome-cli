package cmd

import (
	"fmt"

	"github.com/oliy/pchome-cli/pkg/i18n"
	"github.com/spf13/cobra"
)

func localizeRootCommand(root *cobra.Command, loc *i18n.Catalog) {
	root.SetUsageTemplate(localizedUsageTemplate(loc))

	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	root.InitDefaultVersionFlag()

	localizeCommandTree(root, root, loc)
	localizeBuiltinCommands(root, loc)
}

func localizeCommandTree(root, cmd *cobra.Command, loc *i18n.Catalog) {
	cmd.InitDefaultHelpFlag()

	if flag := cmd.Flags().Lookup("help"); flag != nil {
		if cmd == root {
			flag.Usage = loc.Sprintf(i18n.HelpFlagForCommand, cmd.CommandPath())
		} else {
			flag.Usage = loc.Text(i18n.HelpFlagForThisCommand)
		}
	}

	if cmd == root {
		if flag := cmd.Flags().Lookup("version"); flag != nil {
			flag.Usage = loc.Sprintf(i18n.VersionFlagForCommand, cmd.CommandPath())
		}
	}

	for _, child := range cmd.Commands() {
		localizeCommandTree(root, child, loc)
	}
}

func localizeBuiltinCommands(root *cobra.Command, loc *i18n.Catalog) {
	for _, cmd := range root.Commands() {
		switch cmd.Name() {
		case "help":
			cmd.Short = loc.Text(i18n.HelpCommandShort)
			cmd.Long = loc.Sprintf(i18n.HelpCommandLong, root.CommandPath())
		case "completion":
			cmd.Short = loc.Text(i18n.CompletionCommandShort)
		}
	}
}

func localizedUsageTemplate(loc *i18n.Catalog) string {
	return fmt.Sprintf(`%s{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{if .IsAvailableCommand}}%s{{else}}%s{{end}}
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s
{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

{{printf %q .CommandPath}}
{{end}}`,
		loc.Text(i18n.HelpHeadingUsage),
		loc.Text(i18n.HelpHeadingAliases),
		loc.Text(i18n.HelpHeadingExamples),
		loc.Text(i18n.HelpHeadingAvailableCommands),
		loc.Text(i18n.HelpHeadingAdditionalCommands),
		loc.Text(i18n.HelpHeadingFlags),
		loc.Text(i18n.HelpHeadingGlobalFlags),
		loc.Text(i18n.HelpHeadingAdditionalTopics),
		loc.Text(i18n.HelpMoreInfo),
	)
}
