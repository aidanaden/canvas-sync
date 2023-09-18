/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/initialise"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise new config file with defaults",
	Long: `Initialise new directory '~/.canvas-sync' with default config '~/.canvas-sync/config.yaml'.
Default values:
  - access_token: ""
  - data_dir: ~/.canvas-sync/data
  - canvas_url: https://canvas.nus.edu.sg`,
	Run: func(cmd *cobra.Command, args []string) {
		initialise.RunInit(true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
