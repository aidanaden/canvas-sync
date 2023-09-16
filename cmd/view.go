/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View data from canvas (upcoming lectures, deadlines, etc)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// view events: https://canvas.nus.edu.sg/api/v1/calendar_events?per_page=100&type=assignment&context_codes%5B%5D=course_45742&all_events=1&excludes%5B%5D=assignment&excludes%5B%5D=description&excludes%5B%5D=child_events
// view people: https://canvas.nus.edu.sg/api/v1/courses/45742/users?include%5B%5D=avatar_url&include%5B%5D=enrollments&include%5B%5D=email&include%5B%5D=observed_users&include%5B%5D=can_be_removed&include%5B%5D=custom_links&include_inactive=true&page=2&per_page=50
// view syllabus: server-rendered table, will need to scrape
// view grades: client-rendered table, will need to use playwright
