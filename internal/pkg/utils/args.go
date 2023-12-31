package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

func GetCourseCodesFromArgs(rawArgs []string) []string {
	args := make([]string, 0)
	for _, raw := range rawArgs {
		if strings.Contains(raw, ",") {
			splits := strings.Split(raw, ",")
			for _, split := range splits {
				args = append(args, strings.ToLower(split))
			}
		} else {
			args = append(args, strings.ToLower(raw))
		}
	}
	return args
}

func GetExpandedHomeDirectoryPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		pterm.Error.Printfln("Error expanding home directory: %s", err.Error())
		os.Exit(1)
	}
	if strings.HasPrefix(path, "~") {
		return filepath.Join(home, path[2:])
	}
	return path
}
