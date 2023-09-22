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

func RunViewEvents(cmd *cobra.Command, args []string, isPast bool) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	cookiesFile := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "cookies")
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(http.DefaultClient, canvasUrl, accessToken, cookiesFile)
	if accessToken == "" {
		pterm.Info.Printfln("No access token found, using cookies...")
		canvasClient.ExtractCookies()
	} else {
		pterm.Info.Printfln("Using access token starting with: %s", accessToken[:5])
	}

	courses, err := canvasClient.GetActiveEnrolledCourses()
	if err != nil {
		pterm.Error.Printfln("Error: failed to get actively enrolled courses: %s", err.Error())
		os.Exit(1)
	}

	var events []nodes.EventNode
	if isPast {
		events, err = canvasClient.GetRecentCalendarEvents()
		if err != nil {
			pterm.Error.Printfln("Error: failed to fetch all recent calendar events: %s", err.Error())
			if err := canvasClient.ClearStoredBrowserCookies(); err != nil {
				pterm.Error.Printfln("Error: failed to clear stale cookies in %s: %s", cookiesFile, err.Error())
				os.Exit(1)
			}
		}
	} else {
		events, err = canvasClient.GetIncomingCalendarEvents()
		if err != nil {
			pterm.Error.Printfln("Error: failed to fetch all upcoming calendar events: %s", err.Error())
			os.Exit(1)
		}
	}

	tableData := pterm.TableData{
		{"Course", "Event Type", "Name", "Start date", "Location"},
	}

	for _, event := range events {
		var eventCourseCode string
		for _, c := range courses {
			if c.ID == event.CourseId {
				eventCourseCode = c.CourseCode
				break
			}
		}
		if event.Plannable.AnnouncementPlannableNode != nil {
			tableData = append(tableData, []string{
				eventCourseCode, "Announcement", event.Plannable.Title, "",
			})
		} else if event.Plannable.EventPlannableNode != nil {
			startAt, err := utils.SGTfromUTC(event.Plannable.EventPlannableNode.StartAt)
			if err != nil {
				pterm.Error.Printfln("Error loading event start date: %s", err.Error())
				continue
			}
			startAtStr := utils.FormatEventDate(*startAt)
			if event.Plannable.ZoomPlannableNode != nil {
				tableData = append(tableData, []string{
					eventCourseCode, "Zoom", event.Plannable.Title, startAtStr, event.Plannable.ZoomPlannableNode.OnlineMeetingUrl,
				})
			} else {
				tableData = append(tableData, []string{
					eventCourseCode, "Live", event.Plannable.Title, startAtStr, event.Plannable.LocationName,
				})
			}
		}
	}

	pterm.Println()
	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Render(); err != nil {
		pterm.Error.Printfln("Error rendering events: %s", err.Error())
		os.Exit(1)
	}
}
