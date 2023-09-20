package initialise

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
)

func initConfigFile(path string) {
	configDir := filepath.Dir(path)
	dataDir := filepath.Join(configDir, "data")
	d1 := []byte(
		fmt.Sprintf("# paste your access token below\naccess_token: \n# replace with your preferred location to store canvas data\ndata_dir: %s\n# replace with your own canvas url if not from nus\ncanvas_url: %s\n", dataDir, "https://canvas.nus.edu.sg"),
	)
	if err := os.WriteFile(path, d1, 0755); err != nil {
		pterm.Error.Printfln("Error creating config file: %s", err.Error())
		os.Exit(1)
	}
}

var DEFAULT_CONFIG_DIR = ".canvas-sync"
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
			pterm.Error.Printfln("Error creating default config directory: %s", err.Error())
			os.Exit(1)
		}
		initConfigFile(cfgFilePath)
		if isInitCommand {
			pterm.Success.Printfln("Successfully created config file: %s", cfgFilePath)
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
			initConfigFile(cfgFilePath)
			if isInitCommand {
				pterm.Success.Printfln("Successfully created config file: %s", cfgFilePath)
			}
		} else {
			pterm.Warning.Println("Init command cancelled.")
		}
	} else if os.IsNotExist(err) {
		if !isInitCommand {
			pterm.Warning.Println("No config file found, creating default config...")
		}
		initConfigFile(cfgFilePath)
		if isInitCommand {
			pterm.Success.Println("Successfully created config file.")
		}
	} else if err != nil {
		pterm.Error.Printfln("Error getting config file info: %s", err.Error())
		os.Exit(1)
	}

	return cfgDir
}
