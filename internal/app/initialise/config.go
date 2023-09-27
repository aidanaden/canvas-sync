package initialise

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aidanaden/canvas-sync/internal/pkg/canvas"
	"github.com/aidanaden/canvas-sync/internal/pkg/config"
	"github.com/playwright-community/playwright-go"
	"github.com/pterm/pterm"
)

func generateCanvasAccessToken(page playwright.Page, canvasUrl *url.URL) (*config.GenerateAccessTokenInfo, error) {
	page, loginInfo, err := canvas.LoginToCanvas(page, "", "", canvasUrl)
	if err != nil {
		return nil, err
	}
	if loginInfo == nil {
		return nil, fmt.Errorf("already logged in to canvas")
	}

	profileUrl := url.URL{
		Host:   canvasUrl.Host,
		Scheme: canvasUrl.Scheme,
		Path:   "/profile/settings",
	}

	if _, err := page.Goto(profileUrl.String()); err != nil {
		return nil, fmt.Errorf("failed to navigate to profile url %s: %s", profileUrl.String(), err.Error())
	}

	// check if existing canvas-sync token exists
	tokenLocs, err := page.Locator(".access_token").All()
	if err != nil {
		return nil, fmt.Errorf("failed to get access tokens table: %s", err.Error())
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
			return nil, fmt.Errorf("failed to click existing canvas-sync token details on %s: %s", page.URL(), err.Error())
		}
	}

	if err := page.Locator(".add_access_token_link").Click(); err != nil {
		return nil, fmt.Errorf("failed to open 'new access token' button: %s", err.Error())
	}

	if err := page.Locator("#access_token_purpose").Fill("canvas-sync"); err != nil {
		return nil, fmt.Errorf("failed to fill access token purpose: %s", err.Error())
	}

	if err := page.Locator(".ui-dialog-buttonset").Locator("xpath=/button[2]").Click(); err != nil {
		return nil, fmt.Errorf("failed to click 'generate token' on %s: %s", page.URL(), err.Error())
	}

	tokenLoc := page.Locator(".visible_token")
	if err := tokenLoc.WaitFor(); err != nil {
		return nil, fmt.Errorf("failed to generate token: %s", err.Error())
	}

	accessToken, err := tokenLoc.TextContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get generated token value: %s", err.Error())
	}

	return &config.GenerateAccessTokenInfo{
		Username:    loginInfo.Username,
		Password:    loginInfo.Password,
		AccessToken: accessToken,
	}, nil
}

func initConfigFile(path string) error {
	configDir := filepath.Dir(path)
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

	accessTokenInfo, err := generateCanvasAccessToken(page, parsedCanvasUrl)
	if err != nil {
		return err
	}
	if accessTokenInfo == nil {
		return fmt.Errorf("error generating canvas access token")
	}

	saveCredentials, err := pterm.DefaultInteractiveConfirm.Show("Login is required to download videos - save credentials to config?")
	if err != nil {
		return err
	}

	savedConfig := config.Config{
		DataDir:     dataDir,
		CanvasUrl:   parsedCanvasUrl.String(),
		AccessToken: accessTokenInfo.AccessToken,
	}
	if saveCredentials {
		savedConfig.Username = accessTokenInfo.Username
		savedConfig.Password = accessTokenInfo.Password
	}

	if err := config.SaveConfig(path, &savedConfig, true); err != nil {
		return err
	}

	if err = bw.Close(); err != nil {
		return err
	}
	if err = pw.Stop(); err != nil {
		return err
	}

	return nil
}

func RunInit(isInitCommand bool) string {
	cfgPaths := config.GetConfigPaths()

	// create config directory + file if not exist
	if _, err := os.Stat(cfgPaths.CfgDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgPaths.CfgDirPath, 0755); err != nil {
			pterm.Error.Printfln("Error creating config directory: %s", err.Error())
			os.Exit(1)
		}
		if err := initConfigFile(cfgPaths.CfgFilePath); err != nil {
			pterm.Error.Printfln("Error creating config file: %s", err.Error())
			os.Exit(1)
		}
		if isInitCommand {
			pterm.Success.Printfln("Successfully created config file: %s\n", cfgPaths.CfgFilePath)
		}
		return cfgPaths.CfgDirPath
	}

	_, err := os.Stat(cfgPaths.CfgFilePath)

	// overwrite existing
	if err == nil && isInitCommand {
		res, err := pterm.DefaultInteractiveConfirm.Show("Existing config file found - create new config file?")
		if err != nil {
			pterm.Error.Printfln("Error requesting confirmation: %s", err.Error())
			os.Exit(1)
		}
		if res {
			if err := initConfigFile(cfgPaths.CfgFilePath); err != nil {
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
		if err := initConfigFile(cfgPaths.CfgFilePath); err != nil {
			pterm.Error.Printfln("Error creating config file: %s", err.Error())
			os.Exit(1)
		}
	} else if err != nil {
		pterm.Error.Printfln("Error getting config file info: %s", err.Error())
		os.Exit(1)
	}

	return cfgPaths.CfgDirPath
}
