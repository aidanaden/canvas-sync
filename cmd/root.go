package cmd

import (
	"os"

	"github.com/aidanaden/canvas-sync/internal/app/initialise"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "canvas-sync",
	Short: "Canvas-sync is a CLI tool to interact with canvas",
	Long: `Canvas-sync is a CLI alternative to the canvas website.

Features:
  - download from canvas (files, videos, etc)
  - display canvas info (deadlines, announcements, etc)
  - upload/submit assignments (only if i get > 10 stars on github)
  - more to come (tm)...`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String("access_token", "", "canvas access token; configurable in $HOME/.canvas-sync/config.yaml")
	viper.BindPFlag("access_token", rootCmd.PersistentFlags().Lookup("access_token"))
	rootCmd.PersistentFlags().String("canvas_url", "", "canvas url e.g. canvas.nus.edu.sg; configurable in $HOME/.canvas-sync/config.yaml")
	viper.BindPFlag("canvas_url", rootCmd.PersistentFlags().Lookup("canvas_url"))

	viper.SetDefault("author", "ryan aidan aidan@u.nus.edu")
	viper.SetDefault("license", "MIT")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		_, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := initialise.RunInit(false)

		// Search config in home directory with name ".canvas-sync/config" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		pterm.Println()
		pterm.Info.Printfln("Using config file: %s", viper.ConfigFileUsed())
	}
}
