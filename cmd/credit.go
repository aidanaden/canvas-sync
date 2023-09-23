package cmd

import (
	"time"

	"github.com/aidanaden/canvas-sync/internal/app/credit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// represents the credit command
var creditCmd = &cobra.Command{
	Use:   "credits",
	Short: "Show a list of contributors to canvas-sync project.",
	Long: `Show a list of contributors to canvas-sync project. First one is list of user's with at least one commit in 
the past week. Then a list of all time contributors sorted by number of commits descending.

Examples:
  canvas-sync credits`,
	Run: func(cmd *cobra.Command, args []string) {
        println()

        // Past week contributors
        pterm.FgGreen.Println("Past week contributors:")
        sevenDaysAgo := time.Now().AddDate(0, 0, -7)
        latestContributors := credit.ContributorsSince(sevenDaysAgo)
        for _, c := range latestContributors {
            pterm.FgGreen.Printfln("- %s", c)
        }
        println()

        // Top 500 contributors
        pterm.FgCyan.Println("All contributors:")
        allContributors := credit.Contributors()
        for _, c := range allContributors {
            pterm.FgCyan.Printfln("- %s", c)
        }
        println()
	},
}

func init() {
	rootCmd.AddCommand(creditCmd)
}

