package view

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewEvents(cmd *cobra.Command, args []string) {
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

	pastEvents, err := canvasClient.GetRecentCalendarEvents()
	if err != nil {
		fmt.Printf("\nError: failed to fetch all recent calendar events: %s", err.Error())
		if err := canvasClient.ClearStoredBrowserCookies(); err != nil {
			log.Fatalf("\nError: failed to clear stale cookies in %s: %s", cookiesFile, err.Error())
		}
	}
	for i, raw := range pastEvents {
		fmt.Printf("recent %d: %v\n\n", i+1, raw)
	}

	incomingEvents, err := canvasClient.GetIncomingCalendarEvents()
	if err != nil {
		log.Fatalf("\nError: failed to fetch all upcoming calendar events: %s", err.Error())
	}
	for i, raw := range incomingEvents {
		if raw.Plannable.AnnouncementPlannableNode != nil {
			fmt.Printf("\n\nincoming announcement %d: %v, %v", i+1, raw.Plannable, raw.Plannable.AnnouncementPlannableNode)
		} else if raw.Plannable.AssignmentPlannableNode != nil {
			fmt.Printf("\n\nincoming assignment %d: %v, due at %s", i+1, raw.Plannable, raw.Plannable.AssignmentPlannableNode.DueAt)
		} else if raw.Plannable.EventPlannableNode != nil {
			if raw.Plannable.ZoomPlannableNode != nil {
				fmt.Printf("\n\nincoming zoom event %d: %v, %v", i+1, raw.Plannable, raw.Plannable.EventPlannableNode)
			} else {
				fmt.Printf("\n\nincoming live event %d: %v, %v", i+1, raw.Plannable, raw.Plannable.EventPlannableNode)
			}
		}
	}

}
