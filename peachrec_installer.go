package main

import (
	"bufio"
	"errors"
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
	generateAnimations(workingDir, peachRecAnimations)

	// Generate cfg/peachrec.cfg
	generateConfig(workingDir)

	// Try to added PeachREC to autoexec
	findAutoExec(workingDir)

	pressToExit()
}

func pressToExit() {
	fmt.Print("Press enter to exit.\n")
	fmt.Scanln()
	os.Exit(0)
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
			if _, err := os.Stat(target); err == nil {
				fmt.Println("Found custom HUD:", mods[i].Name())
				return filepath.Join(workingDir, mods[i].Name())
			}
		}
	}

	fmt.Println("No HUD found.")
	// TODO generate with default code if no custom hud is found.
	pressToExit()
	return ""
}

func findAutoExec(workingDir string) {
	// Ask user for permission to append PeachREC to autoexec
	var allowGenerateAutoexec bool = false
	var validGenerateResponse bool = false
	var reader = bufio.NewReader(os.Stdin)

	fmt.Print("Allow program to add \"exec peachrec\" to your autoexec? [Y]/[N]: ")

	// Loop until we get a response that matches our expectation
	for !validGenerateResponse {
		var response string
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline

		// Compare user input to allowed responses and only proceed if allowed
		if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") {
			allowGenerateAutoexec = true
			validGenerateResponse = true
			break
		} else if strings.EqualFold(response, "n") || strings.EqualFold(response, "no") {
			allowGenerateAutoexec = false
			validGenerateResponse = true
			break
		}

		// No valid response
		if !validGenerateResponse {
			fmt.Printf("%v is not a valid option. [Y]/[N]\n", response)
		}
	}

	// User denied permission
	if !allowGenerateAutoexec {
		fmt.Println("Skipped adding PeachREC to autoexec.")
		fmt.Println("Either add \"exec peachrec\" to your autoexec,\nOR add \"+exec peachrec\" to your launch options.")
		fmt.Println("\nPeachREC installed successfully.")
		pressToExit()
	}

	// Check for mastercomfig
	var mastercomfig bool = false
	tfPath := filepath.Dir(workingDir)
	cfgPath := filepath.Join(tfPath, "cfg")
	cfgAutoexecPath := filepath.Join(cfgPath, "autoexec.cfg")
	overridesPath := filepath.Join(cfgPath, "overrides")
	overrideAutoexecPath := filepath.Join(overridesPath, "autoexec.cfg")

	_, err := os.Stat(overrideAutoexecPath)
	if err == nil { // overrides/autoexec.cfg exists
		var validCfgResponse bool = false

		fmt.Print("Program detected mastercomfig is present. Is this correct? [Y]/[N]: ")

		// Loop until we get a response that matches our expectation
		for !validCfgResponse {
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline

			// Compare user input to allowed responses and only proceed if allowed
			if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") { // mastercomfig in use
				mastercomfig = true
				validCfgResponse = true
				generateAutoexec(overrideAutoexecPath)
				break
			} else if strings.EqualFold(response, "n") || strings.EqualFold(response, "no") { // mastercomfig not in use
				validCfgResponse = true
				generateAutoexec(cfgAutoexecPath)
				break
			}

			// No valid response
			if !validCfgResponse {
				fmt.Printf("%v is not a valid option. [Y]/[N]\n", response)
			}
		}
	} else if !mastercomfig || errors.Is(err, os.ErrNotExist) { // If overrides/autoexec.cfg does not exist
		// Check if cfg/autoexec.cfg exists
		_, err := os.Stat(cfgAutoexecPath)
		if err == nil { // cfg/autoexec.cfg
			generateAutoexec(cfgAutoexecPath)
		} else if errors.Is(err, os.ErrNotExist) { // cfg/autoexec.cfg does not exist
			var validCfgResponse bool = false

			fmt.Print("No cfg/autoexist.cfg detected. Generate new one? [Y]/[N]: ")

			// Loop until we get a response that matches our expectation
			for !validCfgResponse {
				var response string
				response, _ = reader.ReadString('\n')  // Read to newline
				response = strings.TrimSpace(response) // Remove newline

				// Compare user input to allowed responses and only proceed if allowed
				if strings.EqualFold(response, "y") || strings.EqualFold(response, "yes") { // Generate new cfg/autoexec.cfg
					validCfgResponse = true
					generateAutoexec(cfgAutoexecPath)
					break
				} else if strings.EqualFold(response, "n") || strings.EqualFold(response, "no") { // Don't generate new cfg/autoexec.cfg
					fmt.Println("Skipped adding PeachREC to autoexec.")
					fmt.Println("Either add \"exec peachrec\" to your autoexec,\nOR add \"+exec peachrec\" to your launch options.")
					fmt.Println("\nPeachREC installed successfully.")
					pressToExit()
				}

				// No valid response
				if !validCfgResponse {
					fmt.Printf("%v is not a valid option. [Y]/[N]\n", response)
				}
			}
		}
	} else { // Unexpected error
		fmt.Println("Error locating mastercomfig autoexec:", err)
		pressToExit()
	}
}
