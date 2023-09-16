/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// videosCmd represents the videos command
var pullVideosCmd = &cobra.Command{
	Use:   "videos",
	Short: "Downloads videos for a given course (all if none specified)",
	Long: `Downloads videos for the given course code(s) case insensitive. If none is specified, all will be downloaded.

Examples:
  canvas-sync pull videos - downloads videos for all courses
  canvas-sync pull videos CS3219 - downloads videos for course with course code "CS3219"
  canvas-sync pull videos CS3219,CS3230,CS1101S - downloads videos for courses with course codes "CS3219", "CS3230" and "CS1101S"`,
	Run: runPullVideos,
}

func init() {
	pullCmd.AddCommand(pullVideosCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// videosCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// videosCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPullVideos(cmd *cobra.Command, args []string) {
	fmt.Println("NOT IMPLEMENTED YET:")
	fmt.Println("Downloading videos is p fking annoying cuz ill need to simulate a browser to get the video urls (thank u canvas) - will add before 1.0 release!")
}
