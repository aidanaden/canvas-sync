package view

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewDeadlines(cmd *cobra.Command, args []string, isPast bool) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	cookiesFile := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "cookies")
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(http.DefaultClient, canvasUrl, accessToken, cookiesFile)
	if accessToken == "" {
		pterm.Info.Printfln("No cookies found, using stored browser cookies...")
		if err := canvasClient.ExtractStoredBrowserCookies(); err != nil {
			pterm.Info.Printfln("No stored cookies found, extracting browser cookies...")
			canvasClient.ExtractBrowserCookies()
			canvasClient.StoreDomainBrowserCookies()
		}
	} else {
		pterm.Info.Printfln("Using access token starting with: %s", accessToken[:5])
	}

	var events []nodes.EventNode
	var err error
	if isPast {
		events, err = canvasClient.GetRecentCalendarEvents()
		if err != nil {
			pterm.Error.Printfln("Error: failed to fetch all recent assignments: %s", err.Error())
			if err := canvasClient.ClearStoredBrowserCookies(); err != nil {
				pterm.Error.Printfln("Error: failed to clear stale cookies in %s: %s", cookiesFile, err.Error())
				os.Exit(1)
			}
		}
	} else {
		events, err = canvasClient.GetIncomingCalendarEvents()
		if err != nil {
			pterm.Error.Printfln("Error: failed to fetch all upcoming assignments: %s", err.Error())
			os.Exit(1)
		}
	}

	tableData := pterm.TableData{
		{"Assignment name", "Due date", "Points possible"},
	}

	for _, event := range events {
		if event.Plannable.AssignmentPlannableNode != nil {
			due, err := utils.SGTfromUTC(event.Plannable.AssignmentPlannableNode.DueAt)
			if err != nil {
				pterm.Error.Printfln("Error loading assignment due date: %s", err.Error())
				continue
			}
			dueStr := utils.FormatEventDate(*due)
			tableData = append(tableData, []string{
				event.Plannable.Title, dueStr, fmt.Sprintf("%d", event.Plannable.AssignmentPlannableNode.PointsPossible),
			})
		}
	}

	pterm.Println()
	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Render(); err != nil {
		pterm.Error.Printfln("Error rendering assignments: %s", err.Error())
		os.Exit(1)
	}
}
