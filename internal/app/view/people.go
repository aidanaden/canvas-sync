package view

import (
	"fmt"
	"net/http"
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
		pterm.Info.Printfln("No cookies found, using stored browser cookies...")
		if err := canvasClient.ExtractStoredBrowserCookies(); err != nil {
			pterm.Info.Printfln("No stored cookies found, extracting browser cookies...")
			canvasClient.ExtractBrowserCookies()
			canvasClient.StoreDomainBrowserCookies()
		}
	} else {
		pterm.Info.Printfln("Using access token starting with: %s", accessToken[:5])
	}

	coursePeople, err := canvasClient.GetCoursePeople(courseCode)
	if err != nil {
		pterm.Error.Printfln("Error: failed to fetch all upcoming calendar events: %s", err.Error())
	}
	for _, people := range coursePeople {
		pterm.Info.Printfln("Name: %s\nImage url: %s", people.ShortName, people.AvatarUrl)
	}
}
