package view

import (
	"fmt"
	"os"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewCoursePeople(cmd *cobra.Command, args []string) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	courseCode := args[0]
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(canvasUrl, accessToken)
	if accessToken == "" {
		pterm.Error.Printfln("Invalid config, please run 'canvas-sync init'")
		os.Exit(1)
	}

	coursePeople, err := canvasClient.GetCoursePeople(courseCode)
	if err != nil {
		pterm.Error.Printfln("Failed to fetch all upcoming calendar events: %s", err.Error())
	}

	tableData := pterm.TableData{
		{"Name"},
	}

	for _, person := range coursePeople {
		tableData = append(tableData, []string{
			person.Name,
		})
	}
	pterm.Println()
	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Render(); err != nil {
		pterm.Error.Printfln("Error rendering people: %s", err.Error())
		os.Exit(1)
	}
	pterm.Info.Printfln("Showing %d people from %s", len(coursePeople), courseCode)
}
