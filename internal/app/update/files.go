package update

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/chelnak/ysmrr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunUpdateFiles(cmd *cobra.Command, args []string) {
	targetDir := fmt.Sprintf("%v/files", viper.Get("data_dir"))
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	providedCodes := utils.GetCourseCodesFromArgs(args)

	fmt.Printf("Files will be downloaded to:\n%s\n", targetDir)

	canvasClient := canvas.NewClient(http.DefaultClient, canvasUrl, accessToken)
	if accessToken == "" {
		fmt.Printf("\nno cookies found, getting auth cookies from browser...")
		canvasClient.ExtractDomainBrowserCookies()
	} else {
		fmt.Printf("\nUsing access token starting with: %s", accessToken[:5])
	}

	rawCourses := canvasClient.GetActiveEnrolledCourses()
	courses := make([]nodes.CourseNode, 0)
	for _, raw := range rawCourses {
		if raw.CourseCode == "" {
			continue
		}
		// add any course if no code provided
		if len(providedCodes) == 0 {
			courses = append(courses, raw)
			continue
		}
		// add course if code matches provided codes
		for _, provided := range providedCodes {
			if strings.ToLower(raw.CourseCode) == provided {
				courses = append(courses, raw)
			}
		}
	}

	var wg sync.WaitGroup
	sm := ysmrr.NewSpinnerManager()

	for ci := range courses {
		wg.Add(1)
		sp := sm.AddSpinner(fmt.Sprintf("Starting files update for %s...", courses[ci].CourseCode))
		go func(i int, sp *ysmrr.Spinner) {
			defer wg.Done()
			id := courses[i].ID
			code := courses[i].CourseCode

			rootNode := canvasClient.GetCourseRootFolder(id)
			rootNode.Name = fmt.Sprintf("%s/%s", targetDir, code)

			sp.UpdateMessagef("Pulling files info for %s", code)
			canvasClient.RecurseDirectoryNode(&rootNode, nil)

			sp.UpdateMessagef("Updating files for %s", code)
			totalFileDownloads := 0

			canvasClient.RecursiveUpdateNode(&rootNode, func(numDownloads int) {
				totalFileDownloads += numDownloads
				sp.UpdateMessagef("Downloading %d files for %s", totalFileDownloads, code)
			})

			sp.UpdateMessagef("Downloaded %d files for %s", totalFileDownloads, code)
			sp.Complete()
		}(ci, sp)
	}
	sm.Start()
	wg.Wait()
	sm.Stop()
	fmt.Printf("\nUpdated files:\n%s\n", targetDir)
}
