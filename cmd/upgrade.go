/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	return false
}

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade canvas-sync to the latest available version",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if isLatestVersion(rootCmd.Version) {
			pterm.Info.Printfln("Canvas-sync is up-to-date!")
			return nil
		}

		cmdToRun := ""
		if IsUnderHomebrew() {
			cmdToRun = "brew upgrade canvas-sync"
		} else {
			return fmt.Errorf("only installs via brew can be upgraded via 'canvas-sync upgrade' :(")
		}

		command := exec.Command("sh", "-c", cmdToRun)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to upgrade canvas-sync with command: %w", err)
		}
		return nil
	},
}
