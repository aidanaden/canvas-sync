/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"flag"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Downloads course files from canvas",
	Long: `
Download files from canvas to a target directory (defaults to $HOME/.canvas-sync/data/files)
Specify target directory in the $HOME/.canvas-sync/config.yml file
`,
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawTargetDir := flag.String("target-dir", ".canvas-sync", "target directory to store downloaded canvas files & directories")
	// rawAccessToken := flag.String("access-token", "", "access token to bypass canvas oauth process")
	rootCmd.PersistentFlags().StringP("data_dir", "d", "", "downloaded data directory; default is $HOME/.canvas-sync/data, configurable in $HOME/.canvas-sync/config.yaml")
	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data_dir"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var DEFAULT_DIR = ".canvas-sync"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [file]\n", os.Args[0])
	flag.PrintDefaults()
}

// pull video: client-rendered panopto thingy - will need to use playwright
