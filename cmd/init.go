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
	Long: `Initialise directory '$HOME/canvas-sync' with default config '$HOME/canvas-sync/config.yaml'.
Default values:
  - access_token: ""
  - data_dir: $HOME/canvas-sync/data
  - canvas_url: https://canvas.nus.edu.sg`,
	Run: func(cmd *cobra.Command, args []string) {
		latestVersionCheck(rootCmd.Version)
		initialise.RunInit(true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
