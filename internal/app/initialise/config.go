package initialise

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aidanaden/canvas-sync/internal/pkg/input"
)

func initConfigFile(path string) {
	configDir := filepath.Dir(path)
	dataDir := filepath.Join(configDir, "data")
	d1 := []byte(
		fmt.Sprintf("# paste your access token below\naccess_token: \n# replace with your preferred location to store canvas data\ndata_dir: %s\n# replace with your own canvas url if not from nus\ncanvas_url: %s\n", dataDir, "https://canvas.nus.edu.sg"),
	)
	if err := os.WriteFile(path, d1, 0755); err != nil {
		log.Fatalf("\nError creating config file: %s", err.Error())
	}
}

var DEFAULT_CONFIG_DIR = ".canvas-sync"
var DEFAULT_CONFIG_FILE = "config.yaml"

func RunInit(isInitCommand bool) string {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("\nError getting user home directory: %s", err.Error())
	}
	cfgDir := filepath.Join(home, DEFAULT_CONFIG_DIR)
	cfgFilePath := filepath.Join(cfgDir, DEFAULT_CONFIG_FILE)

	// create config directory + file if not exist
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfgDir, 0755); err != nil {
			log.Fatalf("\nError creating default config directory: %s", err.Error())
		}
		initConfigFile(cfgFilePath)
		if isInitCommand {
			fmt.Printf("\nSuccessfully created config file.")
		}
	} else {
		_, err = os.Stat(cfgFilePath)
		// overwrite existing
		if err == nil && isInitCommand {
			fmt.Printf("Existing config file found - confirm create new config file? (y/N): ")
			res := input.GetYesOrNoFromUser()
			if res {
				initConfigFile(cfgFilePath)
				if isInitCommand {
					fmt.Printf("\nSuccessfully created config file.")
				}
			} else {
				fmt.Println("Init command cancelled.")
			}
		} else if os.IsNotExist(err) {
			if !isInitCommand {
				fmt.Printf("\nNo config file found, creating default config...\n\n")
			}
			initConfigFile(cfgFilePath)
			if isInitCommand {
				fmt.Printf("\nSuccessfully created config file.")
			}
		} else if err != nil {
			log.Fatalf("\nError getting config file info: %s", err.Error())
		}
	}

	return cfgDir
}
