package view

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunViewCourseAnnouncements(cmd *cobra.Command, args []string) {
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

	courseAnnouncements, err := canvasClient.GetCourseAnnouncements(courseCode)
	if err != nil {
		pterm.Error.Printfln("Error: failed to fetch all course announcements: %s", err.Error())
	}

	tableData := pterm.TableData{
		{"Title", "Posted", "Author", "Message"},
	}

	for i := len(courseAnnouncements) - 1; i > 0; i-- {
		announcement := courseAnnouncements[i]
		postedAtStr := utils.FormatEventDate(announcement.PostedAt)
		tableData = append(tableData, []string{
			announcement.Title, postedAtStr, announcement.PosterName, pterm.DefaultParagraph.WithMaxWidth(80).Sprint(strip.StripTags(announcement.Message)),
		})
	}
	pterm.Println()
	if err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Render(); err != nil {
		pterm.Error.Printfln("Error rendering events: %s", err.Error())
		os.Exit(1)
	}
}
