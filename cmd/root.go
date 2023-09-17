package cmd

import (
	"fmt"
	"os"

	"path/filepath"

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
	cobra.OnInitialize(initConfig)

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

func initConfigFile(path string) {
	configDir := filepath.Dir(path)
	dataDir := fmt.Sprintf("%s/data", configDir)
	d1 := []byte(
		fmt.Sprintf("access_token: \ndata_dir: %s\ncanvas_url: %s\n", dataDir, "https://canvas.nus.edu.sg"),
	)
	if err := os.WriteFile(path, d1, 0644); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := fmt.Sprintf("%s/.canvas-sync", home)
		cfgFilePath := fmt.Sprintf("%s/config.yaml", cfgDir)
		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cfgDir, os.ModePerm); err != nil {
				panic(err)
			}
			initConfigFile(cfgFilePath)
		}

		if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
			initConfigFile(cfgFilePath)
		}

		// Search config in home directory with name ".canvas-sync/config" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Using config file:\n%s\n", viper.ConfigFileUsed()))
	}
}
