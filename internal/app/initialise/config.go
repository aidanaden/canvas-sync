package initialise

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pterm/pterm"
)

type Config struct {
	DataDir     string
	CanvasUrl   string
	AccessToken string
}

const cfg_template = `
# file directory to store canvas data e.g. ~/.canvas-sync/data
data_dir: {{ .DataDir }}
# canvas website url e.g. https://canvas.nus.edu.sg
canvas_url: {{ .CanvasUrl }}
# your canvas access token
access_token: {{ .AccessToken }}
`

func generate(config *Config) string {
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

func initConfigFile(path string) {
	configDir := filepath.Dir(path)
	accessToken := ""
	dataDir := ""
	canvasUrl := ""
	var err error

	pterm.Println()
	dataDir, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter location to store downloaded canvas data (default: ~/.canvas-sync/data)")
	if err != nil {
		pterm.Error.Printfln("Error getting input for data location: %s", err.Error())
		os.Exit(1)
	}
	canvasUrl, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter canvas url (default: https://canvas.nus.edu.sg)")
	if err != nil {
		pterm.Error.Printfln("Error getting input for canvas url: %s", err.Error())
		os.Exit(1)
	}
	accessToken, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter access token (optional, skip if u dont have one)")
	if err != nil {
		pterm.Error.Printfln("Error getting input for access token: %s", err.Error())
		os.Exit(1)
	}

	// set defaults
	if dataDir == "" {
		dataDir = filepath.Join(configDir, "data")
	}
	if canvasUrl == "" {
		canvasUrl = "https://canvas.nus.edu.sg"
	}

	d1 := []byte(generate(&Config{
		DataDir:     dataDir,
		CanvasUrl:   canvasUrl,
		AccessToken: accessToken,
	}))

	if err := os.WriteFile(path, d1, 0755); err != nil {
		pterm.Error.Printfln("Error creating config file: %s", err.Error())
		os.Exit(1)
	}

	pterm.Println()
	pterm.Success.Printfln("Successfully created config file: %s", path)
	pterm.Println()
	td := [][]string{
		{pterm.FgCyan.Sprint("data_dir"), pterm.FgGreen.Sprint(dataDir)},
		{pterm.FgCyan.Sprint("canvas_url"), pterm.FgGreen.Sprint(canvasUrl)},
		{pterm.FgCyan.Sprint("access_token"), pterm.FgGreen.Sprint(accessToken)},
	}
	tablePrint, err := pterm.DefaultTable.WithHasHeader().WithData(td).Srender()
	if err != nil {
		pterm.Error.Printfln("Error rendering config table: %s", err.Error())
		os.Exit(1)
	}
	box := pterm.DefaultBox.WithTitle(pterm.Sprintf("Config: %s", path)).Sprint(tablePrint)
	pterm.Println(box)
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
			pterm.Success.Printfln("Successfully created config file: %s\n", cfgFilePath)
		}
		return cfgDir
	}

	_, err = os.Stat(cfgFilePath)
	// overwrite existing
	if err == nil && isInitCommand {
		pterm.Println()
		pterm.Info.Println("Existing config file found - create new config file?")
		res, err := pterm.DefaultInteractiveConfirm.Show()
		if err != nil {
			pterm.Error.Printfln("Error requesting confirmation: %s", err.Error())
			os.Exit(1)
		}
		if res {
			initConfigFile(cfgFilePath)
		} else {
			pterm.Println()
			pterm.Error.Println("Init command cancelled.")
		}
	} else if os.IsNotExist(err) {
		if !isInitCommand {
			pterm.Warning.Println("No config file found, creating default config...")
		}
		initConfigFile(cfgFilePath)
	} else if err != nil {
		pterm.Error.Printfln("Error getting config file info: %s", err.Error())
		os.Exit(1)
	}

	return cfgDir
}
