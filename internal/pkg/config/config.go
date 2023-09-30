package config

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

type Config struct {
	DataDir     string
	CanvasUrl   string
	Username    string
	Password    string
	AccessToken string
}

const cfg_template = `
data_dir: "{{ .DataDir }}"
canvas_url: "{{ .CanvasUrl }}"
canvas_username: "{{ .Username }}"
canvas_password: "{{ .Password }}"
access_token: "{{ .AccessToken }}"`

func GenerateConfigYaml(config *Config) string {
	cleanedTemplate := strings.ReplaceAll(cfg_template, "\t", " ")
	t, err := template.New("yaml generator").Parse(cleanedTemplate)
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

type GenerateAccessTokenInfo struct {
	AccessToken string
	Username    string
	Password    string
}

func SaveConfig(filepath string, config *Config, verbose bool) error {
	d1 := []byte(GenerateConfigYaml(config))
	if err := os.WriteFile(filepath, d1, 0755); err != nil {
		return err
	}
	pterm.Println()
	pterm.Success.Printfln("Successfully saved config file: %s", filepath)
	if verbose {
		if err := PrintConfig(filepath, config); err != nil {
			return err
		}
	}
	return nil
}

func PrintConfig(path string, config *Config) error {
	td := [][]string{
		{pterm.FgCyan.Sprint("data_dir"), pterm.FgGreen.Sprint(config.DataDir)},
		{pterm.FgCyan.Sprint("canvas_url"), pterm.FgGreen.Sprint(config.CanvasUrl)},
	}
	if config.Username != "" {
		td = append(td, []string{pterm.FgCyan.Sprint("canvas_username"), pterm.FgGreen.Sprint(config.Username)})
	}
	if config.Password != "" {
		var rawTruncated = []byte{}
		for i := 0; i < len(config.Password); i++ {
			if i == 0 || i == len(config.Password)-1 {
				rawTruncated = append(rawTruncated, []byte(config.Password)[i])
			} else {
				rawTruncated = append(rawTruncated, '*')
			}
		}
		truncated := string(rawTruncated)
		td = append(td, []string{pterm.FgCyan.Sprint("canvas_password"), pterm.FgGreen.Sprint(truncated)})
	}
	td = append(td, []string{pterm.FgCyan.Sprint("access_token"), pterm.FgGreen.Sprint(config.AccessToken)})

	tablePrint, err := pterm.DefaultTable.WithHasHeader().WithData(td).Srender()
	if err != nil {
		return err
	}
	box := pterm.DefaultBox.WithTitle(pterm.Sprintf("Config: %s", path)).Sprint(tablePrint)
	pterm.Println(box)
	return nil
}

var DEFAULT_CONFIG_DIR = "canvas-sync"
var DEFAULT_CONFIG_FILE = "config.yaml"

type ConfigPaths struct {
	CfgDirPath  string
	CfgFilePath string
}

func GetConfigPaths() ConfigPaths {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		pterm.Error.Printfln("Error getting user home directory: %s", err.Error())
		os.Exit(1)
	}
	cfgDir := filepath.Join(home, DEFAULT_CONFIG_DIR)
	cfgFilePath := filepath.Join(cfgDir, DEFAULT_CONFIG_FILE)
	return ConfigPaths{
		CfgDirPath:  cfgDir,
		CfgFilePath: cfgFilePath,
	}
}
