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
		pull.RunPullVideos(cmd, args, false)
	},
}

func init() {
	pullCmd.AddCommand(pullVideosCmd)

	pullVideosCmd.PersistentFlags().StringP("canvas_username", "u", "", "canvas username")
	viper.BindPFlag("canvas_username", pullVideosCmd.PersistentFlags().Lookup("canvas_username"))
	pullVideosCmd.PersistentFlags().StringP("canvas_password", "p", "", "canvas password")
	viper.BindPFlag("canvas_password", pullVideosCmd.PersistentFlags().Lookup("canvas_password"))
}
