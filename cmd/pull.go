package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/pull"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the pull command
var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"download", "get"},
	Short:   "Downloads course data from canvas",
	Long: `
Download course data from canvas to a target directory (defaults to $HOME/canvas-sync/data/files)
Specify target directory in the $HOME/canvas-sync/config.yml file
`,
}

// represents the pull files command
var pullFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Downloads files for a given course (all if none specified)",
	Example: `  canvas-sync pull files - downloads files for all courses
  canvas-sync pull files --data_dir /Users/test - downloads files for all courses in the /Users/test/files directory
  canvas-sync pull files CS3219 CS3230 - downloads files for courses with course codes "CS3219" or "CS3230"`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		pull.RunPullFiles(cmd, args)
	},
}

func init() {
	pullCmd.AddCommand(pullFilesCmd)
	rootCmd.AddCommand(pullCmd)

	rootCmd.PersistentFlags().StringP("data_dir", "d", "~/canvas-data", "downloaded data directory")
	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data_dir"))
}
