package pull

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/chelnak/ysmrr"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunPullFiles(cmd *cobra.Command, args []string) {
	targetDir := fmt.Sprintf("%s", viper.Get("data_dir"))
	targetDir = utils.GetExpandedHomeDirectoryPath(targetDir)
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))

	providedCodes := utils.GetCourseCodesFromArgs(args)
	pterm.Info.Printfln("Downloading files to: %s", targetDir)

	if accessToken == "" {
		pterm.Error.Printfln("Invalid config, please run 'canvas-sync init'")
		os.Exit(1)
	}
	canvasClient := canvas.NewClient(canvasUrl, accessToken)

	rawCourses, err := canvasClient.GetActiveEnrolledCourses()
	if err != nil {
		pterm.Error.Printfln("Failed to fetch actively enrolled courses: %s", err.Error())
		os.Exit(1)
	}
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

	pterm.Println()
	var wg sync.WaitGroup
	sm := ysmrr.NewSpinnerManager()

	for ci := range courses {
		wg.Add(1)
		sp := sm.AddSpinner(fmt.Sprintf("Starting files download for %s...", courses[ci].CourseCode))
		go func(i int, sp *ysmrr.Spinner) {
			defer wg.Done()
			id := courses[i].ID
			code := courses[i].CourseCode

			rootNode, err := canvasClient.GetCourseRootFolder(id)
			if err != nil {
				sp.UpdateMessagef(pterm.Error.Sprintf("Failed to fetch course root folder: %s", err.Error()))
				sp.Error()
			}
			rootNode.Name = filepath.Join(targetDir, code, "files")

			sp.UpdateMessagef(pterm.FgCyan.Sprintf("Pulling files info for %s", code))
			if err := canvasClient.RecurseDirectoryNode(rootNode, nil); err != nil {
				sp.UpdateMessagef(pterm.Error.Sprintf("Failed to recurse directories: %s", err.Error()))
				sp.Error()
			}

			sp.UpdateMessagef(pterm.FgCyan.Sprintf("Downloading files for %s", code))
			totalFileDownloads := 0
			if err := canvasClient.RecursiveCreateNode(rootNode, func(numDownloads int) {
				totalFileDownloads += numDownloads
				sp.UpdateMessagef(pterm.FgCyan.Sprintf("Downloading %d files for %s", totalFileDownloads, code))
			}); err != nil {
				sp.UpdateMessagef(pterm.Error.Sprintf("Failed to recurse download files: %s", err.Error()))
				sp.Error()
			}

			sp.UpdateMessagef(pterm.FgGreen.Sprintf("Downloaded %d files for %s", totalFileDownloads, code))
			sp.Complete()
		}(ci, sp)
	}

	sm.Start()
	wg.Wait()
	sm.Stop()
	pterm.Println()
	pterm.Success.Printfln("Downloaded files: %s\n", targetDir)
}
