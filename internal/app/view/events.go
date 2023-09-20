package view

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewEvents(cmd *cobra.Command, args []string, isPast bool) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	cookiesFile := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "cookies")
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	canvasClient := canvas.NewClient(http.DefaultClient, canvasUrl, accessToken, cookiesFile)
	if accessToken == "" {
		fmt.Printf("\nno cookies found, using stored browser cookies...\n\n")
		if err := canvasClient.ExtractStoredBrowserCookies(); err != nil {
			fmt.Printf("no stored cookies found, extracting browser cookies...\n\n")
			canvasClient.ExtractBrowserCookies()
			canvasClient.StoreDomainBrowserCookies()
		}
	} else {
		fmt.Printf("\nUsing access token starting with: %s", accessToken[:5])
	}

	var events []nodes.EventNode
	var err error
	if isPast {
		events, err = canvasClient.GetRecentCalendarEvents()
		if err != nil {
			fmt.Printf("\nError: failed to fetch all recent calendar events: %s", err.Error())
			if err := canvasClient.ClearStoredBrowserCookies(); err != nil {
				log.Fatalf("\nError: failed to clear stale cookies in %s: %s", cookiesFile, err.Error())
			}
		}
	} else {
		events, err = canvasClient.GetIncomingCalendarEvents()
		if err != nil {
			log.Fatalf("\nError: failed to fetch all upcoming calendar events: %s", err.Error())
		}
	}

	for i, event := range events {
		if event.Plannable.AnnouncementPlannableNode != nil {
			fmt.Printf("\n\nincoming announcement %d: %v, %v", i+1, event.Plannable, event.Plannable.AnnouncementPlannableNode)
		} else if event.Plannable.AssignmentPlannableNode != nil {
			fmt.Printf("\n\nincoming assignment %d: %v, due at %s", i+1, event.Plannable, event.Plannable.AssignmentPlannableNode.DueAt)
		} else if event.Plannable.EventPlannableNode != nil {
			if event.Plannable.ZoomPlannableNode != nil {
				fmt.Printf("\n\nincoming zoom event %d: %v, %v", i+1, event.Plannable, event.Plannable.EventPlannableNode)
			} else {
				fmt.Printf("\n\nincoming live event %d: %v, %v", i+1, event.Plannable, event.Plannable.EventPlannableNode)
			}
		}
	}
}
