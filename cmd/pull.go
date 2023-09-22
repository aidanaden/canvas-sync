package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/pull"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the pull command
var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"download", "get"},
	Short:   "Downloads course data from canvas",
	Long: `
Download course data from canvas to a target directory (defaults to $HOME/.canvas-sync/data/files)
Specify target directory in the $HOME/.canvas-sync/config.yml file
`,
}

// represents the pull files command
var pullFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Downloads files for a given course (all if none specified)",
	Long:  `Downloads files for the given course code(s) case insensitive. If none is specified, all will be downloaded.`,
	Example: `  canvas-sync pull files - downloads files for all courses
  canvas-sync pull files --data_dir /Users/test - downloads files for all courses in the /Users/test/files directory
  canvas-sync pull files CS3219 CS3230 - downloads files for courses with course codes "CS3219" or "CS3230"`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		pull.RunPullFiles(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		utils.MotivationalAndStarMessage()
	},
}

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
	pullCmd.AddCommand(pullFilesCmd)
	pullCmd.AddCommand(pullVideosCmd)
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawTargetDir := flag.String("target-dir", ".canvas-sync", "target directory to store downloaded canvas files & directories")
	// rawAccessToken := flag.String("access-token", "", "access token to bypass canvas oauth process")
	rootCmd.PersistentFlags().StringP("data_dir", "d", "", "downloaded data directory")
	viper.BindPFlag("data_dir", rootCmd.PersistentFlags().Lookup("data_dir"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pullCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// pull video: client-rendered panopto thingy - will need to use playwright
