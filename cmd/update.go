package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/update"
	"github.com/spf13/cobra"
)

// represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates locally downloaded course data from canvas",
	Long: `Updates downloaded course data from canvas in a target directory (defaults to $HOME/.canvas-sync/data/files)
Specify target directory in the $HOME/.canvas-sync/config.yml file
`,
}

var updateFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Updates locally downloaded course files from canvas",
	Long: `Updates downloaded files from canvas to a target directory (defaults to $HOME/.canvas-sync/data/files)
Specify target directory in the $HOME/.canvas-sync/config.yml file
`,
	Run: update.RunUpdateFiles,
}

func init() {
	updateCmd.AddCommand(updateFilesCmd)
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
