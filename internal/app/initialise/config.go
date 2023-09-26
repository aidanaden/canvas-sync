package initialise

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/playwright-community/playwright-go"
	"github.com/pterm/pterm"
)

type Config struct {
	DataDir     string
	CanvasUrl   string
	AccessToken string
}

const cfg_template = `
# file directory to store canvas data e.g. $HOME/canvas-sync/data
data_dir: {{ .DataDir }}
# canvas website url e.g. https://canvas.nus.edu.sg
canvas_url: {{ .CanvasUrl }}
# your canvas access token
access_token: {{ .AccessToken }}
`

func generateConfigYaml(config *Config) string {
	t, err := template.New("yaml generator").Parse(cfg_template)
	if err != nil {
		pterm.Error.Printfln("Error generating yaml generator: %s", err.Error())
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, config)
	if err != nil {
		pterm.Error.Printfln("Error generating yaml from config: %s", err.Error())
	}
	return buf.String()
}

func generateCanvasAccessToken(page playwright.Page, canvasUrl *url.URL) (string, error) {
	page, err := canvas.LoginToCanvas(page, canvasUrl)
	if err != nil {
		return "", err
	}

	profileUrl := url.URL{
		Host:   canvasUrl.Host,
		Scheme: canvasUrl.Scheme,
		Path:   "/profile/settings",
	}

	if _, err := page.Goto(profileUrl.String()); err != nil {
		return "", fmt.Errorf("failed to navigate to profile url %s: %s", profileUrl.String(), err.Error())
	}

	// check if existing canvas-sync token exists
	tokenLocs, err := page.Locator(".access_token").All()
	if err != nil {
		return "", fmt.Errorf("failed to get access tokens table: %s", err.Error())
	}

	page.On("dialog", func(al playwright.Dialog) {
		if err := al.Accept(); err != nil {
			pterm.Error.Printfln("Failed to accept dialog with message %s: %s", al.Message(), err.Error())
		}
	})

	// delete existing canvas-sync token
	for _, tokenLoc := range tokenLocs {
		tokenPurpose, err := tokenLoc.Locator(".purpose").TextContent()
		if err != nil {
			continue
		}
		if tokenPurpose != "canvas-sync" {
			continue
		}
		if err := tokenLoc.Locator(".delete_key_link").Click(); err != nil {
			return "", fmt.Errorf("failed to click existing canvas-sync token details: %s", err.Error())
		}
	}

	if err := page.Locator(".add_access_token_link").Click(); err != nil {
		return "", fmt.Errorf("failed to open 'new access token' button: %s", err.Error())
	}

	if err := page.Locator("#access_token_purpose").Fill("canvas-sync"); err != nil {
		return "", fmt.Errorf("failed to fill access token purpose: %s", err.Error())
	}

	if err := page.Locator("xpath=//html/body/div[6]/div[4]/div/button[2]").Click(); err != nil {
		return "", fmt.Errorf("failed to click 'generate token': %s", err.Error())
	}

	tokenLoc := page.Locator(".visible_token")
	if err := tokenLoc.WaitFor(); err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}
	accessToken, err := tokenLoc.TextContent()
	if err != nil {
		return "", fmt.Errorf("failed to get generated token value: %s", err.Error())
	}
	return accessToken, nil
}

func initConfigFile(path string) error {
	configDir := filepath.Dir(path)
	accessToken := ""
	dataDir := ""
	canvasUrl := ""
	var err error

	for err != nil || dataDir == "" {
		dataDir, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter location to store downloaded canvas data (default: $HOME/canvas-sync/data)")
		r := strings.NewReplacer(`"`, "", `'`, "")
		dataDir = r.Replace(dataDir)
		if err != nil {
			return err
		}
		// set default data dir
		if dataDir == "" {
			dataDir = filepath.Join(configDir, "data")
		}
	}

	var parsedCanvasUrl *url.URL
	for err != nil || parsedCanvasUrl == nil {
		canvasUrl, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter canvas url (default: https://canvas.nus.edu.sg)")
		if err != nil {
			return err
		}
		// set default canvas url
		if canvasUrl == "" {
			canvasUrl = "canvas.nus.edu.sg"
		}
		parsedCanvasUrl, err = url.Parse(canvasUrl)
		if err != nil {
			return err
		}
		if parsedCanvasUrl.Scheme == "" {
			parsedCanvasUrl.Scheme = "https"
		}
		// re-parse url to ensure host and schema values are populated
		parsedCanvasUrl, err = url.Parse(parsedCanvasUrl.String())
		if err != nil {
			return err
		}
	}

	if err := playwright.Install(&playwright.RunOptions{Verbose: false}); err != nil {
		pterm.Warning.Println("Failed to install headless chrome, login via username/password disabled.")
	}
	pterm.Info.Println("Logging in to canvas...")
	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	bw, err := pw.Chromium.Launch()
	if err != nil {
		return err
	}
	page, err := bw.NewPage()
	if err != nil {
		return err
	}

	accessToken, err = generateCanvasAccessToken(page, parsedCanvasUrl)
	if err != nil {
		return err
	}

	d1 := []byte(generateConfigYaml(&Config{
		DataDir:     dataDir,
		CanvasUrl:   parsedCanvasUrl.String(),
		AccessToken: accessToken,
	}))

	if err := os.WriteFile(path, d1, 0755); err != nil {
		return err
	}

	pterm.Println()
	pterm.Success.Printfln("Successfully created config file: %s", path)
	pterm.Info.Printfln("You should get an email about an access token being created - if you don't, please contact me on tele @ryaidan")
	pterm.Println()
	td := [][]string{
		{pterm.FgCyan.Sprint("data_dir"), pterm.FgGreen.Sprint(dataDir)},
		{pterm.FgCyan.Sprint("canvas_url"), pterm.FgGreen.Sprint(parsedCanvasUrl.String())},
		{pterm.FgCyan.Sprint("access_token"), pterm.FgGreen.Sprint(accessToken)},
	}
	tablePrint, err := pterm.DefaultTable.WithHasHeader().WithData(td).Srender()
	if err != nil {
		return err
	}
	box := pterm.DefaultBox.WithTitle(pterm.Sprintf("Config: %s", path)).Sprint(tablePrint)
	pterm.Println(box)

	if err = bw.Close(); err != nil {
		return err
	}
	if err = pw.Stop(); err != nil {
		return err
	}

	return nil
}

var DEFAULT_CONFIG_DIR = "canvas-sync"
var DEFAULT_CONFIG_FILE = "config.yaml"

func RunInit(isInitCommand bool) string {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		pterm.Error.Printfln("Error getting user home directory: %s", err.Error())
		os.Exit(1)
	}
	cfgDir := filepath.Join(home, DEFAULT_CONFIG_DIR)
	cfgFilePath := filepath.Join(cfgDir, DEFAULT_CONFIG_FILE)

	// create config directory + file if not exist
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0755); err != nil {
			pterm.Error.Printfln("Error creating config directory: %s", err.Error())
			os.Exit(1)
		}
		if err := initConfigFile(cfgFilePath); err != nil {
			pterm.Error.Printfln("Error creating config file: %s", err.Error())
			os.Exit(1)
		}
		if isInitCommand {
			pterm.Success.Printfln("Successfully created config file: %s\n", cfgFilePath)
		}
		return cfgDir
	}

	_, err = os.Stat(cfgFilePath)
	// overwrite existing
	if err == nil && isInitCommand {
		pterm.Info.Println("Existing config file found - create new config file?")
		res, err := pterm.DefaultInteractiveConfirm.Show()
		if err != nil {
			pterm.Error.Printfln("Error requesting confirmation: %s", err.Error())
			os.Exit(1)
		}
		if res {
			if err := initConfigFile(cfgFilePath); err != nil {
				pterm.Error.Printfln("Error creating config file: %s", err.Error())
				os.Exit(1)
			}
		} else {
			pterm.Error.Println("Init command cancelled.")
		}
	} else if os.IsNotExist(err) {
		if !isInitCommand {
			pterm.Warning.Println("No config file found, creating default config...")
		}
		if err := initConfigFile(cfgFilePath); err != nil {
			pterm.Error.Printfln("Error creating config file: %s", err.Error())
			os.Exit(1)
		}
	} else if err != nil {
		pterm.Error.Printfln("Error getting config file info: %s", err.Error())
		os.Exit(1)
	}

	return cfgDir
}
