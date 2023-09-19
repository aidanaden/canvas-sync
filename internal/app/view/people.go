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

func RunViewCoursePeople(cmd *cobra.Command, args []string) {
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	cookiesFile := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "cookies")
	courseCode := args[0]
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

	coursePeople, err := canvasClient.GetCoursePeople(courseCode)
	if err != nil {
		log.Fatalf("\nError: failed to fetch all upcoming calendar events: %s", err.Error())
	}
	for _, people := range coursePeople {
		fmt.Printf("Name: %s\nImage url: %s\n\n", people.ShortName, people.AvatarUrl)
	}
}
