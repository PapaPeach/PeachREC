package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ProgramDir string = "tf" + string(filepath.Separator) + "custom"

func main() {
	workingDir := locationCheck()
	hud := findHud(workingDir)

	findManifest(hud)
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
		fmt.Printf("Location check failed.\nProgram must be placted in your tf\\custom folder. Program is currently in:\n%v\n", workingDir)
		//os.Exit(0)
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
			fmt.Println("Searching: " + target)
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
