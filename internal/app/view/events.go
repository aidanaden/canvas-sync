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
		log.Printf("\nno cookies found, getting auth cookies from browser...\n")
		if err := canvasClient.ExtractStoredBrowserCookies(); err != nil {
			canvasClient.ExtractBrowserCookies()
			canvasClient.StoreDomainBrowserCookies()
		}
	} else {
		log.Printf("\nUsing access token starting with: %s", accessToken[:5])
	}

	pastEvents, err := canvasClient.GetRecentCalendarEvents()
	if err != nil {
		log.Printf("\nError: failed to fetch all recent calendar events: %s", err.Error())
		if err := canvasClient.ClearStoredBrowserCookies(); err != nil {
			log.Fatalf("\nError: failed to clear stale cookies in %s: %s", cookiesFile, err.Error())
		}
	}
	for i, raw := range pastEvents {
		log.Printf("\n\nrecent %d: %v", i+1, raw)
	}

	incomingEvents, err := canvasClient.GetIncomingCalendarEvents()
	if err != nil {
		log.Fatalf("\nError: failed to fetch all upcoming calendar events: %s", err.Error())
	}
	for i, raw := range incomingEvents {
		if raw.Plannable.AnnouncementPlannableNode != nil {
			log.Printf("\n\nincoming announcement %d: %v, %v", i+1, raw.Plannable, raw.Plannable.AnnouncementPlannableNode)
		} else if raw.Plannable.AssignmentPlannableNode != nil {
			log.Printf("\n\nincoming assignment %d: %v, due at %s", i+1, raw.Plannable, raw.Plannable.AssignmentPlannableNode.DueAt)
		} else if raw.Plannable.EventPlannableNode != nil {
			log.Printf("\n\nincoming event %d: %v, %v", i+1, raw.Plannable, raw.Plannable.EventPlannableNode)
		}
	}

}
