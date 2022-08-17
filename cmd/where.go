package cmd

import (
	"github.com/metafates/mangal/style"
	"github.com/metafates/mangal/where"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(whereCmd)

	whereCmd.Flags().BoolP("config", "c", false, "configuration path")
	whereCmd.Flags().BoolP("sources", "s", false, "sources path")
	whereCmd.Flags().BoolP("logs", "l", false, "logs path")
	whereCmd.MarkFlagsMutuallyExclusive("config", "sources", "logs")
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths for a files related to the mangal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)

		headerStyle := style.Combined(style.Bold, style.HiBlue)

		whereConfig := lo.Must(cmd.Flags().GetBool("config"))
		whereSources := lo.Must(cmd.Flags().GetBool("sources"))
		wherLogs := lo.Must(cmd.Flags().GetBool("logs"))

		title := func(do bool, what, arg string) {
			if do {
				cmd.Printf("%s %s\n", headerStyle(what+"?"), style.Yellow(arg))
			}
		}

		printConfigPath := func(header bool) {
			title(header, "Configuration", "--config")
			cmd.Println(where.Config())
		}

		printSourcesPath := func(header bool) {
			title(header, "Sources", "--sources")
			cmd.Println(where.Sources())
		}

		printLogsPath := func(header bool) {
			title(header, "Logs", "--logs")
			cmd.Println(where.Logs())
		}

		switch {
		case whereConfig:
			printConfigPath(false)
		case whereSources:
			printSourcesPath(false)
		case wherLogs:
			printLogsPath(false)
		default:
			printConfigPath(true)
			cmd.Println()
			printSourcesPath(true)
			cmd.Println()
			printLogsPath(true)
		}
	},
}
