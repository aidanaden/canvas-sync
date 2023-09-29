package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/pull"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateVideosCmd = &cobra.Command{
	Use:   "videos",
	Short: "Updates locally downloaded course videos from canvas",
	Example: `  canvas-sync update videos - updates all downloaded videos for all courses
  canvas-sync update videos CS3219 - updates all videos for course with course code "CS3219"
  canvas-sync update videos CS3219 CS3230 - updates all videos for courses with course codes "CS3219" or "CS3230"`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		pull.RunPullVideos(cmd, args, true)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		utils.MotivationalAndStarMessage()
	},
}

func init() {
	updateCmd.AddCommand(updateVideosCmd)

	updateVideosCmd.PersistentFlags().StringP("username", "u", "", "canvas username")
	viper.BindPFlag("username", updateVideosCmd.PersistentFlags().Lookup("username"))
	updateVideosCmd.PersistentFlags().StringP("password", "p", "", "canvas password")
	viper.BindPFlag("password", updateVideosCmd.PersistentFlags().Lookup("password"))
}
