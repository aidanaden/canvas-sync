package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/update"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates locally downloaded course data from canvas",
	Long: `Updates downloaded course data from canvas in a target directory (defaults to $HOME/canvas-sync/data/files)
`,
}

var updateFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Updates locally downloaded course files from canvas",
	Example: `  canvas-sync update files - updates all downloaded files for all courses
  canvas-sync update files CS3219 - updates all files for course with course code "CS3219"
  canvas-sync update files CS3219 CS3230 - updates all files for courses with course codes "CS3219" or "CS3230"`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		update.RunUpdateFiles(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		utils.MotivationalAndStarMessage()
	},
}

func init() {
	updateCmd.AddCommand(updateFilesCmd)
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")
	updateCmd.PersistentFlags().BoolP("force", "f", false, "overwrite downloaded files if there's a newer version on canvas")
	viper.BindPFlag("force", updateCmd.PersistentFlags().Lookup("force"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("force", "f", false, "Overwrite downloaded files if there's a newer version on canvas")
}
