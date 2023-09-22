package view

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewCoursePeople(cmd *cobra.Command, args []string) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	cookiesFile := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "cookies")
	courseCode := args[0]
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(http.DefaultClient, canvasUrl, accessToken, cookiesFile)
	if accessToken == "" {
		pterm.Info.Printfln("No access token found, using cookies...")
		canvasClient.ExtractCookies()
	} else {
		pterm.Info.Printfln("Using access token starting with: %s", accessToken[:5])
	}

	coursePeople, err := canvasClient.GetCoursePeople(courseCode)
	if err != nil {
		pterm.Error.Printfln("Error: failed to fetch all upcoming calendar events: %s", err.Error())
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
		pterm.Error.Printfln("Error rendering events: %s", err.Error())
		os.Exit(1)
	}
	pterm.Info.Printfln("Showing %d people from %s", len(coursePeople), courseCode)
}
