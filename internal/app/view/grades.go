package view

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func RunViewGrades(cmd *cobra.Command, args []string) {
	// https://canvas.nus.edu.sg/courses/45767/grades
	pterm.Error.Printfln("Viewing grades soon bro")
}
