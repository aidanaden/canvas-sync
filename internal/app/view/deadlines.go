package view

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewDeadlines(cmd *cobra.Command, args []string, isPast bool) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(canvasUrl, accessToken)
	if accessToken == "" {
		pterm.Error.Printfln("Invalid config, please run 'canvas-sync init'")
		os.Exit(1)
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
			pterm.Error.Printfln("Error: failed to fetch all recent assignments: %s", err.Error())
			os.Exit(1)
		}
	} else {
		events, err = canvasClient.GetIncomingCalendarEvents()
		if err != nil {
			pterm.Error.Printfln("Error: failed to fetch all upcoming assignments: %s", err.Error())
			os.Exit(1)
		}
	}

	tableData := pterm.TableData{
		{"Course", "Assignment name", "Due date", "Points possible"},
	}

	for _, event := range events {
		if event.Plannable.AssignmentPlannableNode != nil {
			due, err := utils.SGTfromUTC(event.Plannable.AssignmentPlannableNode.DueAt)
			if err != nil {
				pterm.Error.Printfln("Error loading assignment due date: %s", err.Error())
				continue
			}
			dueStr := utils.FormatEventDate(*due)
			pointsStr := strconv.Itoa(int(event.Plannable.AssignmentPlannableNode.PointsPossible))
			var eventCourseCode string
			for _, c := range courses {
				if c.ID == event.CourseId {
					eventCourseCode = c.CourseCode
					break
				}
			}
			tableData = append(tableData, []string{
				eventCourseCode, event.Plannable.Title, dueStr, pointsStr,
			})
		}
	}

	pterm.Println()
	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Render(); err != nil {
		pterm.Error.Printfln("Error rendering assignments: %s", err.Error())
		os.Exit(1)
	}
}
