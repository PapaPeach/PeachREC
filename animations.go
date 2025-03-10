package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func findManifest(hud string) string {
	// Search for scripts folder
	for _, scripts := range []string{"scripts", "Scripts"} {
		_, err := os.Stat(filepath.Join(hud, scripts))
		if err == nil { // Found scripts, now locate hudanimations_manifest.
			scriptsPath := filepath.Join(hud, scripts)
			for _, manifest := range []string{"hudanimations_manifest.txt", "HudAnimations_Manifest.txt"} {
				fmt.Println("Searching: " + filepath.Join(scriptsPath, manifest))
				_, er := os.Stat(filepath.Join(scriptsPath, manifest))
				if er == nil { // Found hudanimations_manifest
					fmt.Println("Found: " + filepath.Join(scriptsPath, manifest))
					return filepath.Join(scriptsPath, manifest)
				}
			}
		}
	}

	os.Exit(0)
	return "Did not find hudanimations_manifest.txt"
}

func scanManifest(manifest string) []string {
	// Get file names to search
	input, err := os.Open(manifest)
	if err != nil {
		fmt.Printf("Error opening %v for reading: %v\n", manifest, err)
		os.Exit(1)
	}
	defer input.Close()

	// Create slice containing animation files
	var files []string
	scnr := bufio.NewScanner(input)
	for scnr.Scan() {
		line := scnr.Text()
		for _, token := range strings.Fields(line) {
			// If line is commented or incompatible, skip
			if strings.HasPrefix(token, "//") || strings.HasPrefix(token, "../") {
				break
			}
			// If token is an animations file add it to slice
			if strings.Contains(token, ".txt") {
				files = append(files, token)
			}
		}
	}
	return files
}

func scanAnimations(hud string, files []string) []string {
	// Go through each animation file and find the first instance of HintMessageHide and HudTournamentSetupPanelOpen/Close
	for _, file := range files {
		input, err := os.Open(filepath.Join(hud, file))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Error opening %v to scan animations: %v", file, err)
		}

		// Go through current file
		foundHintMessageHide := false
		foundHudTournamentSetupPanelOpen := false
		foundHudTournamentSetupPanelClose := false
		scnr := bufio.NewScanner(input)
		for scnr.Scan() {
			line := scnr.Text()
			for _, token := range strings.Fields(line) {
				if strings.HasPrefix(token, "//") {
					break
				}
				if !foundHintMessageHide && strings.Contains(line, "event HintMessageHide") {
					fmt.Println(file + ": " + line)
					break
				} else if !foundHudTournamentSetupPanelOpen && strings.Contains(line, "event HudTournamentSetupPanelOpen") {
					fmt.Println(file + ": " + line)
					break
				} else if !foundHudTournamentSetupPanelClose && strings.Contains(line, "event HudTournamentSetupPanelClose") {
					fmt.Println(file + ": " + line)
					break
				}
			}
			// If all required animations are found
			if foundHintMessageHide && foundHudTournamentSetupPanelOpen && foundHudTournamentSetupPanelClose {
				break
			}
		}
		input.Close()
	}
	return nil
}
