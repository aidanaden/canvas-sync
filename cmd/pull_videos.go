package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/pull"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the pull videos command
var pullVideosCmd = &cobra.Command{
	Use:   "videos",
	Short: "Downloads videos for a given course (all if none specified)",
	Long:  `Downloads videos for the given course code(s) case insensitive. If none is specified, all will be downloaded.`,
	Example: `  canvas-sync pull videos - downloads videos for all courses
  canvas-sync pull videos CS3219 CS3230 - downloads videos for courses with course codes "CS3219" or "CS3230"`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		pull.RunPullVideos(cmd, args)
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawTargetDir := flag.String("target-dir", "canvas-sync", "target directory to store downloaded canvas files & directories")
	// rawAccessToken := flag.String("access-token", "", "access token to bypass canvas oauth process")

	pullVideosCmd.PersistentFlags().StringP("username", "u", "", "canvas username")
	viper.BindPFlag("username", pullVideosCmd.PersistentFlags().Lookup("username"))
	pullVideosCmd.PersistentFlags().StringP("password", "p", "", "canvas password")
	viper.BindPFlag("password", pullVideosCmd.PersistentFlags().Lookup("password"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
