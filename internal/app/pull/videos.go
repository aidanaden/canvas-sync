package pull

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/config"
	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
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

func RunPullVideos(cmd *cobra.Command, args []string) {
	if err := playwright.Install(&playwright.RunOptions{Verbose: false}); err != nil {
		pterm.Error.Println("Error setting up playwright. Please create an issue on https://github.com/aidanaden/canvas-sync/issues")
		os.Exit(1)
	}

	targetDir := fmt.Sprintf("%s", viper.Get("data_dir"))
	targetDir = utils.GetExpandedHomeDirectoryPath(targetDir)
	username := fmt.Sprintf("%v", viper.Get("username"))
	password := fmt.Sprintf("%v", viper.Get("password"))

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

	pterm.Info.Printfln("Getting videos for %d courses", len(courses))
	for _, c := range courses {
		courseVideosDir := filepath.Join(targetDir, c.CourseCode, "videos")
		if err := os.MkdirAll(courseVideosDir, 0755); err != nil {
			pterm.Error.Printfln("error creating videos directory for %s: %s", c.CourseCode, err.Error())
			continue
		}
		videoUrls, err := canvasClient.GetCourseVideos(page, c)
		if err != nil {
			pterm.Error.Printfln("error getting videos for %s: %s", c.CourseCode, err.Error())
		}
		var wg sync.WaitGroup
		for video, urls := range videoUrls {
			wg.Add(1)
			go func(video string, urls *canvas.CourseVideoUrls) {
				defer wg.Done()
				video = strings.ReplaceAll(video, "/", "-")
				videoFilename := filepath.Join(courseVideosDir, fmt.Sprintf(`%s.mp4`, video))
				if _, err := os.Stat(videoFilename); err == nil {
					pterm.Info.Printfln("%s already downloaded in %s", video, videoFilename)
					return
				}
				if urls.VideoUrl == "" {
					if err := ffmpeg_go.
						Input(urls.AudioUrl).
						Output(videoFilename, ffmpeg_go.KwArgs{"c": "copy"}).
						ErrorToStdOut().
						Run(); err != nil {
						pterm.Error.Printfln("Error downloading video %s: %s", video, err.Error())
					} else {
						pterm.Success.Println("successfully downloaded n merged video: ", time.Now())
					}
				} else {
					// merge audio n video
					main := ffmpeg_go.Input(urls.VideoUrl)
					overlay := ffmpeg_go.Input(urls.AudioUrl)
					err := ffmpeg_go.Output(
						[]*ffmpeg_go.Stream{main, overlay},
						videoFilename,
						ffmpeg_go.KwArgs{"map": "1:a,0:v", "c:v": "copy"},
					).
						ErrorToStdOut().
						Run()
					if err != nil {
						pterm.Error.Printfln("Error downloading n merging video %s: %s", video, err.Error())
					} else {
						pterm.Success.Printfln("successfully downloaded n merged video %s: %s", video, time.Now())
					}
				}
			}(video, urls)
		}
		wg.Wait()
	}

	// 1. convert main screen to 30fps
	// ffmpeg -i index1.mp4 -filter:v fps=29.72 index1-30fps.mp4
	// 2. compress lecturer screen to 204 x 116
	// ffmpeg -i index2.mp4 -vf "scale=204:116" index2-mini.mp4
	// 3. overlay compressed lecturer screen on main screen
	// ffmpeg -i index1-30fps.mp4 -i index2-mini.mp4 -filter_complex 'overlay=main_w-overlay_w-10:10' overlayed.mp4

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	pterm.Info.Println("starting index1 download: ", time.Now())
	// 	// main := ffmpeg_go.Input(`https://s-cloudfront.cdn.ap.panopto.com/sessions/ad21a9bd-ac56-4fac-844c-b05e006a054b/dfa7e8f4-8164-49cd-af7c-b0820062e24f.object.hls/228532/index.m3u8`).Filter("fps", ffmpeg_go.Args{"2972/100"})
	// 	// overlay := ffmpeg_go.Input(`https://s-cloudfront.cdn.ap.panopto.com/sessions/ad21a9bd-ac56-4fac-844c-b05e006a054b/1f0c2f51-c6b1-4b17-9595-b05e006a0553-10eaf9d2-8529-49e2-966c-b08200835cf8.hls/1088712/index.m3u8`).Filter("scale", ffmpeg_go.Args{"256:144"})
	// 	// err := ffmpeg_go.Filter([]*ffmpeg_go.Stream{
	// 	// 	main,
	// 	// 	overlay,
	// 	// }, "overlay", ffmpeg_go.Args{"main_w-overlay_w-10:10"}).
	// 	// 	Output(fmt.Sprintf("%s.mp4", "merged"), ffmpeg_go.KwArgs{"map": "1:a"}).
	// 	// 	OverWriteOutput().
	// 	// 	ErrorToStdOut().
	// 	// 	Run()
	// 	main := ffmpeg_go.Input(`https://s-cloudfront.cdn.ap.panopto.com/sessions/ad21a9bd-ac56-4fac-844c-b05e006a054b/dfa7e8f4-8164-49cd-af7c-b0820062e24f.object.hls/228532/index.m3u8`)
	// 	overlay := ffmpeg_go.Input(`https://s-cloudfront.cdn.ap.panopto.com/sessions/ad21a9bd-ac56-4fac-844c-b05e006a054b/1f0c2f51-c6b1-4b17-9595-b05e006a0553-10eaf9d2-8529-49e2-966c-b08200835cf8.hls/1088712/index.m3u8`)
	// 	err := ffmpeg_go.Output([]*ffmpeg_go.Stream{main, overlay}, fmt.Sprintf("%s.mp4", "merged2"), ffmpeg_go.KwArgs{"map": "1:a,0:v", "c:v": "copy"}).
	// 		OverWriteOutput().
	// 		ErrorToStdOut().
	// 		Run()
	// 	if err != nil {
	// 		pterm.Error.Printfln("Error downloading n merging videos: %s", err.Error())
	// 	} else {
	// 		pterm.Success.Println("successfully downloaded n merged video: ", time.Now())
	// 	}
	// }()
	// wg.Wait()

	// pterm.Error.Println("NOT IMPLEMENTED YET:")
	// pterm.Info.Println("Downloading videos is p fking annoying cuz ill need to simulate a browser to get the video urls (thank u canvas) - will add before 1.0 release!")
}
