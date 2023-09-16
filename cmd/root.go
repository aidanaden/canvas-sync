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
  - downloading files from canvas (files, videos, etc)
  - displaying course calendars
  - more to come...

Report bugs to https://github.com/aidanaden/canvas-sync/issues`,
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

	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfigFile(path string) {
	configDir := filepath.Dir(path)
	dataDir := fmt.Sprintf("%s/data", configDir)
	d1 := []byte(fmt.Sprintf("access_token: \ndata_dir: %s\n", dataDir))
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
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
