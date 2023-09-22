package cmd

import (
	"errors"

	"github.com/aidanaden/canvas-sync/internal/app/view"
	"github.com/spf13/cobra"
)

var FUTURE_ALIASES = []string{"future", "coming", "next"}
var PAST_ALIASES = []string{"completed", "done"}
var VIEW_ALIASES = []string{"display", "print"}
var VIEW_PEOPLE_ALIASES = []string{"classmates", "class", "students"}

// represents the view command
var viewCmd = &cobra.Command{
	Use:     "view",
	Aliases: VIEW_ALIASES,
	Short:   "View data from canvas (events, deadlines, people)",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// represents the view people command
var viewPeopleCmd = &cobra.Command{
	Use:     "people",
	Aliases: VIEW_PEOPLE_ALIASES,
	Short:   "View people from a given course",
	Long:    `View all students from a given course code (case-insensitive).`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		view.RunViewCoursePeople(cmd, args)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no valid course code provided")
		}
		return nil
	},
	Example: "  canvas-sync view people cs3230",
}

// represents the view events command
var viewEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "View upcoming/past events",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// view upcoming events command
var viewUpcomingEventsCmd = &cobra.Command{
	Use:     "upcoming",
	Aliases: FUTURE_ALIASES,
	Short:   "View upcoming events",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		view.RunViewEvents(cmd, args, false)
	},
}

// view past events command
var viewPastEventsCmd = &cobra.Command{
	Use:     "past",
	Aliases: PAST_ALIASES,
	Short:   "View past/completed events",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		view.RunViewEvents(cmd, args, true)
	},
}

// deadlinesCmd represents the deadlines command
var viewDeadlinesCmd = &cobra.Command{
	Use:     "deadlines",
	Aliases: []string{"assignments"},
	Short:   "View past/future assignment deadlines",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// view upcoming deadlines command
var viewUpcomingDeadlinesCmd = &cobra.Command{
	Use:     "upcoming",
	Aliases: FUTURE_ALIASES,
	Short:   "View upcoming/future assignment deadlines",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		view.RunViewDeadlines(cmd, args, false)
	},
}

// view past deadlines command
var viewPastDeadlinesCmd = &cobra.Command{
	Use:     "past",
	Aliases: PAST_ALIASES,
	Short:   "View completed/past assignment deadlines",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		preRun(cmd)
		view.RunViewDeadlines(cmd, args, true)
	},
}

func init() {
	viewEventsCmd.AddCommand(viewUpcomingEventsCmd)
	viewEventsCmd.AddCommand(viewPastEventsCmd)
	viewCmd.AddCommand(viewEventsCmd)

	viewDeadlinesCmd.AddCommand(viewUpcomingDeadlinesCmd)
	viewDeadlinesCmd.AddCommand(viewPastDeadlinesCmd)
	viewCmd.AddCommand(viewDeadlinesCmd)

	viewCmd.AddCommand(viewPeopleCmd)
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
