package pull

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/chelnak/ysmrr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getCourseCodesFromArgs(rawArgs []string) []string {
	args := make([]string, 0)
	for _, raw := range rawArgs {
		if strings.Contains(raw, ",") {
			splits := strings.Split(raw, ",")
			for _, split := range splits {
				args = append(args, strings.ToLower(split))
			}
		} else {
			args = append(args, strings.ToLower(raw))
		}
	}
	return args
}

func RunPullFiles(cmd *cobra.Command, args []string) {
	targetDir := fmt.Sprintf("%v/files", viper.Get("data_dir"))
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	providedCodes := getCourseCodesFromArgs(args)

	fmt.Printf("files will be downloaded to data_dir: %s", targetDir)

	canvasClient := canvas.NewClient(http.DefaultClient, "canvas.nus.edu.sg", accessToken)
	if accessToken == "" {
		fmt.Printf("\nno cookies found, getting auth cookies from browser...")
		canvasClient.ExtractDomainBrowserCookies()
	}

	fmt.Printf("\naccess token found: %s", accessToken)
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
		sp := sm.AddSpinner(fmt.Sprintf("Starting files download for %s...", courses[ci].CourseCode))
		go func(i int, sp *ysmrr.Spinner) {
			defer wg.Done()
			id := courses[i].ID
			code := courses[i].CourseCode

			rootNode := canvasClient.GetCourseRootFolder(id)
			rootNode.Name = fmt.Sprintf("%s/%s", targetDir, code)

			sp.UpdateMessagef("Pulling files info for %s", code)
			canvasClient.RecurseDirectoryNode(&rootNode, nil)

			sp.UpdateMessagef("Downloading files for %s", code)
			canvasClient.RecursiveCreateNode(&rootNode)

			sp.UpdateMessagef("Downloaded files for %s", code)
			sp.Complete()
		}(ci, sp)
	}
	sm.Start()
	wg.Wait()
	sm.Stop()
	fmt.Printf("\ndownloaded files - view here: %s", targetDir)
}
