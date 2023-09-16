/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/aidanaden/canvas-sync/internal/app/pull_files"
	"github.com/spf13/cobra"
)

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "Downloads files for a given course (all if none specified)",
	Long: `Downloads files for the given course code(s) case insensitive. If none is specified, all will be downloaded.

Examples:
  canvas-sync pull files - downloads files for all courses
  canvas-sync pull files --data_dir /Users/test - downloads files for all courses in the /Users/test/files directory
  canvas-sync pull files CS3219 - downloads files for course with course code "CS3219"
  canvas-sync pull files CS3219,CS3230,CS1101S - downloads files for courses with course codes "CS3219", "CS3230" and "CS1101S"`,
	Run: pull_files.Run,
}

func init() {
	pullCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
