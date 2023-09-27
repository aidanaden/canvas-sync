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

func getPage() (playwright.Page, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	CHANNEL := "chrome"
	bw, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Channel: &CHANNEL})
	if err != nil {
		return nil, err
	}
	page, err := bw.NewPage(playwright.BrowserNewPageOptions{Viewport: &playwright.Size{Height: 1080, Width: 1920}})
	if err != nil {
		return nil, err
	}
	return page, nil
}

func RunPullVideos(cmd *cobra.Command, args []string, isUpdate bool) {
	if err := playwright.Install(&playwright.RunOptions{Verbose: false}); err != nil {
		pterm.Error.Println("Error setting up playwright. Please create an issue on https://github.com/aidanaden/canvas-sync/issues")
		os.Exit(1)
	}

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
	}
	pterm.Info.Printfln("canvas url: %s", parsedCanvasUrl.String())

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

	page, err := getPage()
	if err != nil {
		pterm.Error.Printfln("Error getting page: %s", err.Error())
		os.Exit(1)
	}

	page, loginInfo, err := canvas.LoginToCanvas(page, username, password, parsedCanvasUrl)
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
	courseVideoUrls := make(map[string]map[string]*canvas.CourseVideoUrls, len(courses))
	for _, c := range courses {
		videoUrls, err := canvasClient.GetCourseVideos(page, c)
		if err != nil {
			pterm.Error.Printfln("Error getting videos for %s: %s", c.CourseCode, err.Error())
			continue
		}

		if len(videoUrls) == 0 {
			pterm.Warning.Printfln("No videos found for %s", c.CourseCode)
			continue
		}

		filteredVideoUrls := map[string]*canvas.CourseVideoUrls{}
		for video, videoUrlInfos := range videoUrls {
			video = strings.ReplaceAll(video, "/", "-")
			courseVideosDir := filepath.Join(targetDir, c.CourseCode, "videos")
			videoFilename := filepath.Join(courseVideosDir, fmt.Sprintf(`%s.mp4`, video))
			if !isUpdate {
				// add video url if 'pull'
				filteredVideoUrls[video] = videoUrlInfos
			} else {
				// add video url if not exists if 'update'
				if _, err := os.Stat(videoFilename); err != nil {
					filteredVideoUrls[video] = videoUrlInfos
				}
			}
		}

		courseVideosDir := filepath.Join(targetDir, c.CourseCode, "videos")
		if err := os.MkdirAll(courseVideosDir, 0755); err != nil {
			pterm.Error.Printfln("error creating videos directory for %s: %s", c.CourseCode, err.Error())
			continue
		}
		courseVideoUrls[c.CourseCode] = filteredVideoUrls
	}

	pterm.Println()
	var wg sync.WaitGroup
	sm := ysmrr.NewSpinnerManager(
		ysmrr.WithCompleteColor(colors.FgHiGreen),
		ysmrr.WithSpinnerColor(colors.FgHiBlue),
	)

	for code, videoUrls := range courseVideoUrls {
		if len(videoUrls) == 0 {
			sm.AddSpinner(pterm.FgGreen.Sprintf("All videos already downloaded for %s", code)).Complete()
			continue
		}
		wg.Add(1)
		sp := sm.AddSpinner(pterm.FgCyan.Sprintf("Downloading %d videos for %s...", len(videoUrls), code))
		go func(code string, sp *ysmrr.Spinner, videoUrls map[string]*canvas.CourseVideoUrls) {
			defer wg.Done()
			var videoWg sync.WaitGroup
			numDownloaded := 0
			for video, urls := range videoUrls {
				video = strings.ReplaceAll(video, "/", "-")
				courseVideosDir := filepath.Join(targetDir, code, "videos")
				videoFilename := filepath.Join(courseVideosDir, fmt.Sprintf(`%s.mp4`, video))

				if urls.VideoUrl == "" {
					// if only 1 source file, simply download
					err = ffmpeg_go.
						Input(urls.AudioUrl).
						Output(videoFilename, ffmpeg_go.KwArgs{"c": "copy"}).
						Silent(true).
						OverWriteOutput().
						WithOutput(io.Discard).
						Run()
				} else {
					// if 2 source files, merge audio of audio file into video files
					main := ffmpeg_go.Input(urls.VideoUrl)
					overlay := ffmpeg_go.Input(urls.AudioUrl)
					err = ffmpeg_go.Output(
						[]*ffmpeg_go.Stream{main, overlay},
						videoFilename,
						ffmpeg_go.KwArgs{"map": "1:a,0:v", "c:v": "copy"},
					).
						Silent(true).
						OverWriteOutput().
						WithOutput(io.Discard).
						Run()
				}
				numDownloaded += 1
				if err != nil {
					pterm.Error.Printfln("Error downloading video %s: %s", video, err.Error())
				} else {
					sp.UpdateMessagef(pterm.FgCyan.Sprintf("Downloaded %d/%d videos for %s", numDownloaded, len(videoUrls), code))
				}
			}

			videoWg.Wait()
			sp.UpdateMessagef(pterm.FgGreen.Sprintf("Downloaded %d videos for %s", len(videoUrls), code))
			sp.Complete()
		}(code, sp, videoUrls)
	}

	if len(sm.GetSpinners()) > 0 {
		sm.Start()
		wg.Wait()
		sm.Stop()
	}
	pterm.Println()
	pterm.Success.Printfln("Downloaded videos: %s\n", targetDir)
}
