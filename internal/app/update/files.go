package update

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

func RunUpdateFiles(cmd *cobra.Command, args []string) {
	targetDir := filepath.Join(fmt.Sprintf("%s", viper.Get("data_dir")), "files")
	targetDir = utils.GetExpandedHomeDirectoryPath(targetDir)

	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	updateStaleFiles := viper.Get("force").(bool)

	providedCodes := utils.GetCourseCodesFromArgs(args)

	pterm.Info.Printfln("Downloading files to: %s", targetDir)
	canvasClient := canvas.NewClient(canvasUrl, accessToken)
	if accessToken == "" {
		pterm.Error.Printfln("Invalid config, please run 'canvas-sync init'")
		os.Exit(1)
	}

	rawCourses, err := canvasClient.GetActiveEnrolledCourses()
	if err != nil {
		pterm.Error.Printfln("Error: failed to fetch all actively enrolled courses: %s", err.Error())
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
		sp := sm.AddSpinner(fmt.Sprintf("Starting files update for %s...", courses[ci].CourseCode))
		go func(i int, sp *ysmrr.Spinner, updateStaleFiles bool) {
			defer wg.Done()
			id := courses[i].ID
			code := courses[i].CourseCode

			rootNode, err := canvasClient.GetCourseRootFolder(id)
			if err != nil {
				pterm.Error.Printfln("Error: failed to fetch course root folder: %s", err.Error())
				os.Exit(1)
			}
			rootNode.Name = fmt.Sprintf("%s/%s", targetDir, code)

			sp.UpdateMessagef(pterm.FgCyan.Sprintf("Pulling files info for %s", code))
			if err := canvasClient.RecurseDirectoryNode(rootNode, nil); err != nil {
				sp.UpdateMessagef(pterm.Error.Sprintf("Error: failed to recurse directories: %s", err.Error()))
				sp.Error()
			}

			sp.UpdateMessagef(pterm.FgCyan.Sprintf("Updating files for %s", code))
			totalFileDownloads := 0

			if err := canvasClient.RecursiveUpdateNode(rootNode, updateStaleFiles, func(numDownloads int) {
				totalFileDownloads += numDownloads
				sp.UpdateMessagef(pterm.FgCyan.Sprintf("Downloading %d files for %s", totalFileDownloads, code))
			}); err != nil {
				sp.UpdateMessagef(pterm.Error.Sprintf("Error: failed to recurse update files: %s", err.Error()))
				sp.Error()
			}

			if totalFileDownloads > 0 {
				sp.UpdateMessagef(pterm.FgGreen.Sprintf("Downloaded %d files for %s", totalFileDownloads, code))
			} else {
				sp.UpdateMessagef(pterm.FgGreen.Sprintf("All files are up-to-date"))
			}

			sp.Complete()
		}(ci, sp, updateStaleFiles)
	}

	sm.Start()
	wg.Wait()
	sm.Stop()
	pterm.Println()
	pterm.Success.Printfln("Updated files: %s\n", targetDir)
}
