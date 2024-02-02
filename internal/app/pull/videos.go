package pull

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/config"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/playwright-community/playwright-go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func getBrowser() (playwright.Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	CHANNEL := "chrome"
	bw, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Channel: &CHANNEL})
	if err != nil {
		return nil, err
	}
	return bw, nil
}

func getPage(bw playwright.Browser) (playwright.Page, error) {
	page, err := bw.NewPage(playwright.BrowserNewPageOptions{Viewport: &playwright.Size{Height: 1600, Width: 1920}})
	if err != nil {
		return nil, err
	}
	return page, nil
}

func extractVideosFromDirectory(folder *canvas.CourseVideoFolder) []*canvas.CourseVideoFile {
	flattened := []*canvas.CourseVideoFile{}
	flattened = append(flattened, folder.Videos...)
	for _, fold := range folder.Folders {
		flattened = append(flattened, extractVideosFromDirectory(fold)...)
	}
	return flattened
}

const MAX_DOWNLOAD_ATTEMPTS = 5

func RunPullVideos(cmd *cobra.Command, args []string, isUpdate bool) {
	targetDir := fmt.Sprintf("%s", viper.Get("data_dir"))
	targetDir = utils.GetExpandedHomeDirectoryPath(targetDir)
	username := fmt.Sprintf("%v", viper.Get("canvas_username"))
	password := fmt.Sprintf("%v", viper.Get("canvas_password"))

	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))
	if accessToken == "" {
		pterm.Error.Printfln("Invalid config, please run 'canvas-sync init'")
		os.Exit(1)
	}
	providedCodes := utils.GetCourseCodesFromArgs(args)

	canvasUrl := fmt.Sprintf("%v", viper.Get("canvas_url"))
	parsedCanvasUrl, err := url.Parse(canvasUrl)
	if err != nil {
		pterm.Error.Printfln("%s is an invalid canvas url: %s", canvasUrl, err.Error())
		os.Exit(1)
	}

	canvasClient := canvas.NewClient(canvasUrl, accessToken)
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

	bw, err := getBrowser()
	if err != nil {
		pterm.Error.Printfln("Error getting browser: %s", err.Error())
		os.Exit(1)
	}
	defer bw.Close()

	page, err := getPage(bw)
	if err != nil {
		pterm.Error.Printfln("Error getting page: %s", err.Error())
		os.Exit(1)
	}

	_, loginInfo, err := canvas.LoginToCanvas(page, username, password, parsedCanvasUrl)
	if err != nil {
		pterm.Error.Printfln("Error logging in to canvas: %s", err.Error())
		os.Exit(1)
	}

	existingConfig := config.Config{
		DataDir:     targetDir,
		CanvasUrl:   parsedCanvasUrl.String(),
		AccessToken: accessToken,
		Username:    loginInfo.Username,
		Password:    loginInfo.Password,
	}

	// only save credentials if none provided
	if username == "" || password == "" {
		saveCredentials, err := pterm.DefaultInteractiveConfirm.Show("Login is required to download videos - save credentials to config?")
		if err != nil {
			pterm.Error.Printfln("Error getting save credentials user input: %s", err.Error())
			os.Exit(1)
		}
		configPaths := config.GetConfigPaths()
		if saveCredentials {
			if err := config.SaveConfig(configPaths.CfgFilePath, &existingConfig, false); err != nil {
				pterm.Error.Printfln("Error saving credentials to config: %s", err.Error())
				os.Exit(1)
			}
		}
	}

	pterm.Info.Printfln("Getting videos for %d courses", len(courses))

	sm := ysmrr.NewSpinnerManager(
		ysmrr.WithCompleteColor(colors.FgHiGreen),
		ysmrr.WithSpinnerColor(colors.FgHiBlue),
	)

	type SpinnerCount struct {
		course      string
		sp          *ysmrr.Spinner
		fileCount   int
		folderCount int
	}

	courseSpinners := make(map[string]*SpinnerCount)
	incrementCount := func(spc *SpinnerCount, isFile bool) {
		if isFile {
			spc.fileCount += 1
			spc.sp.UpdateMessage(pterm.FgCyan.Sprintf("Extracting %d file(s) from %s", spc.fileCount, spc.course))
		} else {
			spc.folderCount += 1
			spc.sp.UpdateMessage(pterm.FgCyan.Sprintf("Extracting %d folder(s) from %s", spc.folderCount, spc.course))
		}
	}

	for _, c := range courses {
		sp := sm.AddSpinner(pterm.FgCyan.Sprintf("Extracting video files for %s", c.CourseCode))
		spc := SpinnerCount{
			course:      c.CourseCode,
			sp:          sp,
			fileCount:   0,
			folderCount: 0,
		}
		courseSpinners[c.CourseCode] = &spc
	}

	sm.Start()
	var wg sync.WaitGroup

	pterm.Println()

	for _, course := range courses {
		wg.Add(1)
		go func(c nodes.CourseNode, spc *SpinnerCount) {
			defer wg.Done()

			code := c.CourseCode
			page, err := getPage(bw)
			if err != nil {
				return
			}
			page, _, err = canvas.LoginToCanvas(page, loginInfo.Username, loginInfo.Password, parsedCanvasUrl)
			if err != nil {
				courseSpinners[code].sp.UpdateMessage(pterm.Error.Sprintf("Error logging in to canvas: %s", err.Error()))
				return
			}

			rootFolder, err := canvasClient.GetCourseVideos(page, targetDir, c, func(isFile bool) {
				incrementCount(courseSpinners[code], isFile)
			})
			if err != nil {
				courseSpinners[code].sp.UpdateMessage(pterm.FgRed.Sprintf("No videos found for %s", code))
				courseSpinners[code].sp.Error()
				return
			}

			filtered := []*canvas.CourseVideoFile{}
			files := extractVideosFromDirectory(rootFolder)
			for _, fil := range files {
				if !isUpdate {
					filtered = append(filtered, fil)
				} else if !fil.Downloaded && isUpdate {
					filtered = append(filtered, fil)
				}
			}

			if len(filtered) == 0 {
				if len(files) > 0 {
					courseSpinners[code].sp.UpdateMessage(pterm.FgGreen.Sprintf("All videos already downloaded for %s", code))
				} else {
					courseSpinners[code].sp.UpdateMessage(pterm.FgGreen.Sprintf("No videos available for %s", code))
				}
				courseSpinners[code].sp.Complete()
				return
			}

			// reset spinner count
			courseSpinners[code].fileCount = 0
			courseSpinners[code].sp.UpdateMessage(pterm.FgCyan.Sprintf("Downloading %d videos for %s...", len(files), c.CourseCode))

			for _, fil := range filtered {
				path := fil.Path
				audioUrl := fil.AudioUrl
				videoUrl := fil.VideoUrl

				parent := filepath.Dir(path)
				if err := os.MkdirAll(parent, 0755); err != nil {
					pterm.Error.Printf("failed to create parent directory for file %s, skipping", path)
				}

				if videoUrl == "" {
					// if only 1 source file, simply download
					for i := 0; i < MAX_DOWNLOAD_ATTEMPTS; i++ {
						err = ffmpeg_go.
							Input(audioUrl).
							Output(path, ffmpeg_go.KwArgs{"c": "copy"}).
							Silent(true).
							OverWriteOutput().
							WithOutput(io.Discard).
							Run()
						if err == nil {
							break
						}
					}
				} else {
					// if 2 source files, merge audio of audio file into video files
					main := ffmpeg_go.Input(videoUrl)
					overlay := ffmpeg_go.Input(audioUrl)
					for i := 0; i < MAX_DOWNLOAD_ATTEMPTS; i++ {
						err = ffmpeg_go.Output(
							[]*ffmpeg_go.Stream{main, overlay},
							path,
							ffmpeg_go.KwArgs{"map": "1:a,0:v", "c:v": "copy"},
						).
							Silent(true).
							OverWriteOutput().
							WithOutput(io.Discard).
							Run()
						if err == nil {
							break
						}
					}
				}

				if err != nil {
					pterm.Error.Printfln("Error downloading video %s: %s", path, err.Error())
				}

				spc.fileCount += 1
				spc.sp.UpdateMessagef(pterm.FgCyan.Sprintf("Downloaded %d/%d videos for %s", spc.fileCount, len(filtered), code))
			}

			spc.sp.UpdateMessagef(pterm.FgGreen.Sprintf("Downloaded %d videos for %s", len(filtered), code))
			spc.sp.Complete()
		}(course, courseSpinners[course.CourseCode])
	}

	wg.Wait()
	if len(sm.GetSpinners()) > 0 {
		sm.Stop()
	}

	pterm.Println()
	pterm.Success.Printfln("Downloaded videos: %s", targetDir)
}
