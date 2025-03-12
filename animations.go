package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func findManifest(hud string) string {
	// Search for scripts folder
	for _, scripts := range []string{"scripts", "Scripts"} {
		_, err := os.Stat(filepath.Join(hud, scripts))
		if err == nil { // Found scripts, now locate hudanimations_manifest.
			scriptsPath := filepath.Join(hud, scripts)
			for _, manifest := range []string{"hudanimations_manifest.txt", "HudAnimations_Manifest.txt"} {
				_, er := os.Stat(filepath.Join(scriptsPath, manifest))
				if er == nil { // Found hudanimations_manifest
					return filepath.Join(scriptsPath, manifest)
				}
			}
		}
	}

	os.Exit(0)
	return "Did not find hudanimations_manifest.txt"
}

func scanManifest(manifest string) ([]string, []string) {
	// Get file names to search
	input, err := os.Open(manifest)
	if err != nil {
		fmt.Printf("Error opening %v for reading: %v\n", manifest, err)
		os.Exit(1)
	}
	defer input.Close()

	// Create slice containing animation files
	var files []string
	var hudAnimationsManifest []string
	scnr := bufio.NewScanner(input)
	for scnr.Scan() {
		line := scnr.Text()
		// If line is a known PeachREC line
		if strings.Contains(line, "peachrec") {
			continue
		}
		hudAnimationsManifest = append(hudAnimationsManifest, line)
		// Make a list of files to search for code
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
	return hudAnimationsManifest, files
}

func insertPeachRecManifest(hudAnimationsManifest []string) []string {
	var peachRecManifest []string
	// Insert PeachREC file path to the top of the manifest
	for f, line := range hudAnimationsManifest {
		if strings.Contains(line, "{") {
			peachRecManifest = slices.Insert(hudAnimationsManifest, f+1, "\tfile\tscripts/hudanimations_peachrec.txt")
			return peachRecManifest
		}
	}
	return hudAnimationsManifest
}

func scanAnimations(hud string, files []string) ([]string, []string, []string) {
	// Go through each animation file and find the first instance of HintMessageHide and HudTournamentSetupPanelOpen/Close
	var HintMessageHide []string
	var HudTournamentSetupPanelOpen []string
	var HudTournamentSetupPanelClose []string

	for _, file := range files {
		input, err := os.Open(filepath.Join(hud, strings.ReplaceAll(file, "\"", "")))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Error opening %v to scan animations: %v", file, err)
			os.Exit(1)
		}

		// Go through current file
		foundHintMessageHide := 0
		foundHudTournamentSetupPanelOpen := 0
		foundHudTournamentSetupPanelClose := 0
		scnr := bufio.NewScanner(input)
		for scnr.Scan() {
			line := scnr.Text()
			for _, token := range strings.Fields(line) {
				if strings.HasPrefix(token, "//") { // Ignore commented lines
					break
				} else if strings.Contains(line, "PeachRec") { // Skip known PeachREC lines
					break
				}

				if (foundHintMessageHide < 2 && strings.Contains(line, "event HintMessageHide")) || foundHintMessageHide == 1 { // Copy HintMessageHide
					// Found animation header, now copy subsequent lines
					foundHintMessageHide = 1
					HintMessageHide = append(HintMessageHide, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHintMessageHide = 2
					}
					break
				} else if (foundHudTournamentSetupPanelOpen < 2 && strings.Contains(line, "event HudTournamentSetupPanelOpen")) || foundHudTournamentSetupPanelOpen == 1 { // Copy HudTournamentSetupPanelOpen
					// Found animation header, now copy subsequent lines
					foundHudTournamentSetupPanelOpen = 1
					HudTournamentSetupPanelOpen = append(HudTournamentSetupPanelOpen, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHudTournamentSetupPanelOpen = 2
					}
					break
				} else if (foundHudTournamentSetupPanelClose < 2 && strings.Contains(line, "event HudTournamentSetupPanelClose")) || foundHudTournamentSetupPanelClose == 1 { // Copy HudTournamentSetupPanelClose
					// Found animation header, now copy subsequent lines
					foundHudTournamentSetupPanelClose = 1
					HudTournamentSetupPanelClose = append(HudTournamentSetupPanelClose, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHudTournamentSetupPanelClose = 2
					}
					break
				}
			}
			// If all required animations are found
			if foundHintMessageHide == 2 && foundHudTournamentSetupPanelOpen == 2 && foundHudTournamentSetupPanelClose == 2 {
				return HintMessageHide, HudTournamentSetupPanelOpen, HudTournamentSetupPanelClose
			}
		}
		input.Close()
	}
	return HintMessageHide, HudTournamentSetupPanelOpen, HudTournamentSetupPanelClose
}

func insertPeachRecAnimations(HintMessageHide []string, HudTournamentSetupPanelOpen []string, HudTournamentSetupPanelClose []string) []string {
	var peachRecAnimations []string
	// Insert animations for HintMessageHide
	for a, line := range HintMessageHide {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(HintMessageHide)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecSpawn 0.0")
		}
	}
	// Insert animations for HudTournamentSetupPanelOpen
	for a, line := range HudTournamentSetupPanelOpen {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(HudTournamentSetupPanelOpen)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecOpen 0.0")
		}
	}
	// Insert animations for HudTournamentSetupPanelClose
	for a, line := range HudTournamentSetupPanelClose {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(HudTournamentSetupPanelClose)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecClose 0.0")
		}
	}

	// Add custom PeachREC animation events
	peachRecAnimations = append(peachRecAnimations, "event PeachRecSpawn")
	peachRecAnimations = append(peachRecAnimations, "{")
	peachRecAnimations = append(peachRecAnimations, "\tFireCommand 0.0 \"engine peachrec\"")
	peachRecAnimations = append(peachRecAnimations, "}")

	peachRecAnimations = append(peachRecAnimations, "event PeachRecOpen")
	peachRecAnimations = append(peachRecAnimations, "{")
	peachRecAnimations = append(peachRecAnimations, "\tFireCommand 0.001 \"engine pr_open\"")
	peachRecAnimations = append(peachRecAnimations, "}")

	peachRecAnimations = append(peachRecAnimations, "event PeachRecClose")
	peachRecAnimations = append(peachRecAnimations, "{")
	peachRecAnimations = append(peachRecAnimations, "\tFireCommand 0.0 \"engine pr_close\"")
	peachRecAnimations = append(peachRecAnimations, "}")

	return peachRecAnimations
}
