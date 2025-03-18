package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const ProgramDir string = "tf" + string(filepath.Separator) + "custom"
const modNameWindows string = "_PeachREC"
const modNameUnix string = "_peachrec"

var modName string

func main() {
	// Verify location
	workingDir := locationCheck()

	// Detect system OS
	if runtime.GOOS == "windows" {
		modName = modNameWindows
		fmt.Println("Windows OS detected, using mod name \"_PeachREC\"")
	} else {
		modName = modNameUnix
		fmt.Println("Unix OS detected, using mod name \"_peachrec\"")
	}

	// Locate custom HUD
	hud := findHud(workingDir)
	// If no custom HUD is found, use default values
	if hud == "" {
		generateDefaultManifest(workingDir)
		hud = filepath.Join(workingDir, modName)
	}

	// Locate hudanimations_manifest.txt
	manifest := findManifest(hud)
	fmt.Println()

	// Scan manifest for animation files
	hudAnimationsManifest, animationFiles := scanManifest(manifest)

	// Insert PeachREC animations file in manifest
	peachRecManifest := insertPeachRecManifest(hudAnimationsManifest)

	// Scan animation files for HudTournamentSetupPanelOpen/Close
	hintMessageHide, hudTournamentSetupPanelOpen, hudTournamentSetupPanelClose, hudReadyPulseEnd := scanAnimations(hud, animationFiles)
	fmt.Println()

	// Insert PeachREC animations
	peachRecAnimations := insertPeachRecAnimations(hintMessageHide, hudTournamentSetupPanelOpen, hudTournamentSetupPanelClose, hudReadyPulseEnd)

	// Generate updated hudanimations_manifest.txt
	generateManifest(manifest, peachRecManifest)

	// Generate hudanimations_peachrect.txt
	generateAnimations(workingDir, peachRecAnimations)

	// Generate cfg/peachrec.cfg
	generateConfig(workingDir)

	// Try to added PeachREC to autoexec
	findAutoExec(workingDir)

	// Tell user program is done
	fmt.Println("\nPeachREC installed successfully.\nEnjoy <3")
	pressToExit()
}

func pressToExit() {
	fmt.Println("\nPress enter to exit")
	_, err := fmt.Scanln()
	if err != nil {
		// Intentionally ignore error because error or not, closing program is the goal
		os.Exit(0)
	}
	os.Exit(0)
}

func locationCheck() string {
	// Get filepath of working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error checking location of program:", err)
		pressToExit()
	}

	// Check working directory is the desired location
	if strings.HasSuffix(workingDir, ProgramDir) {
		fmt.Println("Verified program is in the correct location")
	} else {
		fmt.Printf("Program must be placted in your tf\\custom folder. Program is currently in:\n%v\n", workingDir)
		pressToExit()
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
			if _, err = os.Stat(target); err == nil {
				fmt.Println("Found custom HUD:", mods[i].Name())
				return filepath.Join(workingDir, mods[i].Name())
			}
		}
	}

	fmt.Println("No custom HUD found, generating with default hud values")
	return ""
}

func findAutoExec(workingDir string) {
	// Ask user for permission to append PeachREC to autoexec
	var allowGenerateAutoexec bool
	var validGenerateResponse bool
	var reader = bufio.NewReader(os.Stdin)

	fmt.Print("\nAllow program to add \"exec peachrec\" to your autoexec? [Y]/[N]: ")

	// Loop until we get a response that matches our expectation
	for !validGenerateResponse {
		var response string
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline

		// Compare user input to allowed responses and only proceed if allowed
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			allowGenerateAutoexec = true
			validGenerateResponse = true
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response, "no") {
			allowGenerateAutoexec = false
			validGenerateResponse = true
		}

		// No valid response
		if !validGenerateResponse {
			fmt.Printf("%v is not a valid option. [Y]/[N]: ", response)
		}
	}

	// User denied permission
	if !allowGenerateAutoexec {
		fmt.Println("Skipped adding PeachREC to autoexec")
		fmt.Println("Either add \"exec peachrec\" to your autoexec,\nOR add \"+exec peachrec\" to your launch options")
		fmt.Println("\nPeachREC installed successfully\nEnjoy <3")
		pressToExit()
	}

	// Check for mastercomfig
	var mastercomfig bool
	tfPath := filepath.Dir(workingDir)
	cfgPath := filepath.Join(tfPath, "cfg")
	cfgAutoexecPath := filepath.Join(cfgPath, "autoexec.cfg")
	overridesPath := filepath.Join(cfgPath, "overrides")
	overrideAutoexecPath := filepath.Join(overridesPath, "autoexec.cfg")

	_, err := os.Stat(overrideAutoexecPath)
	switch {
	case err == nil: // overrides/autoexec.cfg exists
		var validCfgResponse bool

		fmt.Print("Program detected mastercomfig is present. Is this correct? [Y]/[N]: ")

		// Loop until we get a response that matches our expectation
		for !validCfgResponse {
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline

			// Compare user input to allowed responses and only proceed if allowed
			switch {
			// mastercomfig in use
			case strings.EqualFold(response, "y"), strings.EqualFold(response, "yes"):
				validCfgResponse = true
				generateAutoexec(overrideAutoexecPath)
			// mastercomfig not in use
			case strings.EqualFold(response, "n"), strings.EqualFold(response, "no"):
				validCfgResponse = true
				generateAutoexec(cfgAutoexecPath)
			default:
				fmt.Printf("%v is not a valid option. [Y]/[N]: ", response)
			}
		}
	case !mastercomfig || errors.Is(err, os.ErrNotExist): // If overrides/autoexec.cfg does not exist
		// Check if cfg/autoexec.cfg exists
		_, err = os.Stat(cfgAutoexecPath)
		switch {
		case err == nil: // cfg/autoexec.cfg
			generateAutoexec(cfgAutoexecPath)
		case errors.Is(err, os.ErrNotExist): // cfg/autoexec.cfg does not exist
			var validCfgResponse bool

			fmt.Print("No tf/cfg/autoexec.cfg detected. Generate a new one? [Y]/[N]: ")

			// Loop until we get a response that matches our expectation
			for !validCfgResponse {
				var response string
				response, _ = reader.ReadString('\n')  // Read to newline
				response = strings.TrimSpace(response) // Remove newline

				// Compare user input to allowed responses and only proceed if allowed
				switch {
				// Generate new cfg/autoexec.cfg
				case strings.EqualFold(response, "y"), strings.EqualFold(response, "yes"):
					validCfgResponse = true
					generateAutoexec(cfgAutoexecPath)
				// Don't generate new cfg/autoexec.cfg
				case strings.EqualFold(response, "n"), strings.EqualFold(response, "no"):
					fmt.Println("Skipped adding PeachREC to autoexec")
					fmt.Println("Either add \"exec peachrec\" to your autoexec,\nOR add \"+exec peachrec\" to your launch options")
					fmt.Println("\nPeachREC installed successfully\nEnjoy <3")
					pressToExit()
				default:
					fmt.Printf("%v is not a valid option. [Y]/[N]\n", response)
				}
			}
		default: // Unexpected error
			fmt.Println("Error locating mastercomfig autoexec:", err)
			pressToExit()
		}
	default: // Unexpected error
		fmt.Println("Error locating mastercomfig autoexec:", err)
		pressToExit()
	}
}
