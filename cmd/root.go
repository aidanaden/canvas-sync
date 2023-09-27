package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aidanaden/canvas-sync/internal/app/initialise"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
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

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
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

	rootCmd.PersistentFlags().StringP("access_token", "a", "", "canvas access token")
	viper.BindPFlag("access_token", rootCmd.PersistentFlags().Lookup("access_token"))
	rootCmd.PersistentFlags().StringP("canvas_url", "c", "https://canvas.nus.edu.sg", "canvas url e.g. canvas.nus.edu.sg")
	viper.BindPFlag("canvas_url", rootCmd.PersistentFlags().Lookup("canvas_url"))

	viper.SetDefault("author", "ryan aidan aidan@u.nus.edu")
	viper.SetDefault("license", "MIT")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func isLatestVersion(rawVersionStr string) bool {
	splits := strings.Split(strings.ReplaceAll(rawVersionStr, ")", ""), " ")
	currHash := splits[len(splits)-1]
	latestReleaseInfo, err := utils.GetCavasSyncLatestVersionHash()
	if err != nil {
		return false
	}
	if currHash == "" {
		return false
	}
	if currHash != latestReleaseInfo.CommitHash {
		return false
	}
	return true
}

func latestVersionCheck(versionStr string) {
	if !isLatestVersion(rootCmd.Version) {
		pterm.Warning.Printfln("New version of canvas-sync available, run 'canvas-sync upgrade' to update!")
		pterm.Println()
	}
}

// preRun reads in config file and ENV variables if set, verifies current app version
func preRun(cmd *cobra.Command) {
	latestVersionCheck(rootCmd.Version)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		_, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := initialise.RunInit(false)

		// Search config in home directory with name "canvas-sync/config" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		pterm.Info.Printfln("Using config file: %s", viper.ConfigFileUsed())
	} else {
		pterm.Error.Printfln("Error reading config: %s", err.Error())
		os.Exit(1)
	}
}
