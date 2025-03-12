package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ProgramDir string = "tf" + string(filepath.Separator) + "custom"

func main() {
	// Verify location
	workingDir := locationCheck()

	// Locate custom Hud
	hud := findHud(workingDir)

	// Locate hudanimations_manifest.txt
	manifest := findManifest(hud)

	// Scan manifest for animation files
	hudAnimationsManifest, animationFiles := scanManifest(manifest)

	// Insert PeachREC animations file in manifest
	peachRecManifest := insertPeachRecManifest(hudAnimationsManifest)

	// Scan animation files for HudTournamentSetupPanelOpen/Close
	HintMessageHide, HudTournamentSetupPanelOpen, HudTournamentSetupPanelClose := scanAnimations(hud, animationFiles)

	// Insert PeachREC animations
	peachRecAnimations := insertPeachRecAnimations(HintMessageHide, HudTournamentSetupPanelOpen, HudTournamentSetupPanelClose)

	// Generate updated hudanimations_manifest.txt
	generateManifest(manifest, peachRecManifest)

	// Generate hudanimations_peachrect.txt
	generateAnimations(hud, peachRecAnimations)

	// Generate cfg/peachrec.cfg
	generateConfig(hud)
}

func locationCheck() string {
	// Get filepath of working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error checking location of program.")
	}

	// Check working directory is the desired location
	if strings.HasSuffix(workingDir, ProgramDir) {
		fmt.Println("Location check passed.")
	} else {
		fmt.Printf("Location check failed.\nProgram must be placted in your tf\\custom folder. Program is currently in:\n%v\n\n", workingDir)
		//os.Exit(1)
	}

	return workingDir
}

func findHud(workingDir string) string {
	// Get list of all mods in custom
	mods, err := os.ReadDir(workingDir)
	if err != nil {
		fmt.Println("Error getting list of mods in custom")
	}

	// Search each mod for info.vdf and return first result
	for i := range mods {
		if mods[i].IsDir() {
			target := filepath.Join(mods[i].Name(), "info.vdf")
			if _, err := os.Stat(target); err == nil {
				return filepath.Join(workingDir, mods[i].Name())
			}
		}
	}

	fmt.Println("No HUD found.")
	// TODO generate with default code if no custom hud is found.
	os.Exit(0)
	return ""
}

func findAutoExec() {

}
