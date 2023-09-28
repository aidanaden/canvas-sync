/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func IsUnderHomebrew() bool {
	binary, err := os.Executable()
	if err != nil {
		return false
	}

	brewExe, err := exec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(binary, brewBinPrefix)
}

func IsUnderScoop() bool {
	binary, err := os.Executable()
	if err != nil {
		return false
	}

	scoopExe, err := exec.LookPath("scoop")
	if err != nil {
		return false
	}

	rawScoopCanvasPath, err := exec.Command(scoopExe, "which", "canvas-sync").Output()
	if err != nil {
		return false
	}

	scoopCanvasPath := utils.GetExpandedHomeDirectoryPath(string(rawScoopCanvasPath))
	return strings.EqualFold(binary, scoopCanvasPath)
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade canvas-sync to the latest available version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if isLatestVersion(rootCmd.Version) {
			pterm.Success.Printfln("Canvas-sync is up-to-date!")
			return
		}

		cmdToRun := ""
		var command *exec.Cmd

		if IsUnderHomebrew() {
			cmdToRun = "brew update && brew upgrade canvas-sync"
			command = exec.Command("sh", "-c", cmdToRun)
		} else if IsUnderScoop() {
			cmdToRun = "scoop update; scoop update canvas-sync"
			command = exec.Command("cmd", "/C", cmdToRun)
		} else {
			pterm.Error.Printfln("Only installs via brew/scoop can be upgraded via 'canvas-sync upgrade' :(")
			return
		}

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			pterm.Error.Printfln("Failed to upgrade canvas-sync with command: %s", err.Error())
		} else {
			pterm.Success.Printfln("Canvas-sync successfully updated!")
		}
	},
}
