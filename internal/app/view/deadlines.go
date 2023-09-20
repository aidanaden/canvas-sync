package view

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunViewDeadlines(cmd *cobra.Command, args []string) {
	pterm.Error.Printfln("Viewing deadlines soon bro")
}
