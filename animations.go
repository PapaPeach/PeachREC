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

	fmt.Println("Did not find hudanimations_manifest.txt")
	pressToExit()
	return ""
}

func scanManifest(manifest string) ([]string, []string) {
	// Get file names to search
	input, err := os.Open(manifest)
	if err != nil {
		fmt.Printf("Error opening %v for reading: %v\n", manifest, err)
		pressToExit()
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
			peachRecManifest = slices.Insert(hudAnimationsManifest, f+1, "\tfile\thudanimations_peachrec.txt")
			return peachRecManifest
		}
	}
	return hudAnimationsManifest
}

func scanAnimations(hud string, files []string) ([]string, []string, []string, []string) {
	// Go through each animation file and find the first instance of HintMessageHide and HudTournamentSetupPanelOpen/Close
	var hintMessageHide []string
	var hudTournamentSetupPanelOpen []string
	var hudTournamentSetupPanelClose []string
	var hudReadyPulseEnd []string

	foundHintMessageHide := 0
	foundHudTournamentSetupPanelOpen := 0
	foundHudTournamentSetupPanelClose := 0
	foundHudReadyPulseEnd := 0

	for _, file := range files {
		input, err := os.Open(filepath.Join(hud, strings.ReplaceAll(file, "\"", "")))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Error opening %v to scan animations: %v", file, err)
			pressToExit()
		}

		// Go through current file

		scnr := bufio.NewScanner(input)
		for scnr.Scan() {
			line := scnr.Text()
			for _, token := range strings.Fields(line) {
				if strings.HasPrefix(token, "//") { // Ignore commented lines
					break
				} else if strings.Contains(line, "PeachRec") { // Skip known PeachREC lines
					break
				}
				switch {
				// Copy HintMessageHide
				case (foundHintMessageHide < 2 && strings.Contains(line, "event HintMessageHide")) || foundHintMessageHide == 1:
					// Found animation header, now copy subsequent lines
					foundHintMessageHide = 1
					hintMessageHide = append(hintMessageHide, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHintMessageHide = 2
						fmt.Println("Found custom HintMessageHide animation, using that")
					}

				// Copy HudTournamentSetupPanelOpen
				case (foundHudTournamentSetupPanelOpen < 2 && strings.Contains(line, "event HudTournamentSetupPanelOpen")) || foundHudTournamentSetupPanelOpen == 1:
					// Found animation header, now copy subsequent lines
					foundHudTournamentSetupPanelOpen = 1
					hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHudTournamentSetupPanelOpen = 2
						fmt.Println("Found custom HudTournamentSetupPanelOpen animation, using that")
					}

				// Copy HudTournamentSetupPanelClose
				case (foundHudTournamentSetupPanelClose < 2 && strings.Contains(line, "event HudTournamentSetupPanelClose")) || foundHudTournamentSetupPanelClose == 1:
					// Found animation header, now copy subsequent lines
					foundHudTournamentSetupPanelClose = 1
					hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHudTournamentSetupPanelClose = 2
						fmt.Println("Found custom HudTournamentSetupPanelClose animation, using that")
					}

				// Copy HudReadyPulseEnd
				case (foundHudReadyPulseEnd < 2 && strings.Contains(line, "event HudReadyPulseEnd")) || foundHudReadyPulseEnd == 1:
					// Found animation header, now copy subsequent lines
					foundHudReadyPulseEnd = 1
					hudReadyPulseEnd = append(hudReadyPulseEnd, line)
					if strings.Contains(line, "}") { // Stop copying once a close brace is copies
						foundHudReadyPulseEnd = 2
						fmt.Println("Found custom HudReadyPulseEnd animation, using that")
					}
				}
			}
			// If all required animations are found
			if foundHintMessageHide == 2 && foundHudTournamentSetupPanelOpen == 2 && foundHudTournamentSetupPanelClose == 2 && foundHudReadyPulseEnd == 2 {
				return hintMessageHide, hudTournamentSetupPanelOpen, hudTournamentSetupPanelClose, hudReadyPulseEnd
			}
		}
		input.Close()
	}

	// If no custom HintMessageHide animation is found, use default code
	if foundHintMessageHide == 0 {
		hintMessageHide = append(hintMessageHide, "event HintMessageHide")
		hintMessageHide = append(hintMessageHide, "{")
		hintMessageHide = append(hintMessageHide, "\tAnimate HudHintDisplay\tFgColor\t\"255 220 0 0\"\tLinear\t0.0\t0.2")
		hintMessageHide = append(hintMessageHide, "\tAnimate HudHintDisplay\tHintSize\t\"0\"\tDeaccel 0.2\t0.3")
		hintMessageHide = append(hintMessageHide, "}")

		fmt.Println("Did not find custom HintMessageHide animation, using default")
	} else if foundHintMessageHide == 1 {
		fmt.Println("Error: Found HintMessageHide animation header, but did not close")
		pressToExit()
	}

	// If no custom HudTournamentSetupPanelOpen animation is found, use default code
	if foundHudTournamentSetupPanelOpen == 0 {
		hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, "event HudTournamentSetupPanelOpen")
		hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, "{")
		hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, "\tAnimate HudTournamentSetupt\tPosition\t\"c-90 -70\"\tLinear 0.0 0.001")
		hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, "\tAnimate HudTournamentSetup\tPosition\t\"c-90 70\"\tSpline 0.001 0.2")
		hudTournamentSetupPanelOpen = append(hudTournamentSetupPanelOpen, "}")

		fmt.Println("Did not find custom HudTournamentSetupPanelOpen animation, using default")
	} else if foundHudTournamentSetupPanelOpen == 1 {
		fmt.Println("Error: Found HudTournamentSetupPanelOpen animation header, but did not close")
		pressToExit()
	}

	// If no custom HudTournamentSetupPanelClose animation is found, use default code
	if foundHudTournamentSetupPanelClose == 0 {
		hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, "event HudTournamentSetupPanelClose")
		hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, "{")
		hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, "\tAnimate HudTournamentSetup\tPosition\t\"c-90 70\"\tLinear 0.0 0.001")
		hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, "\tAnimate HudTournamentSetup\tPosition\t\"c-90 -70\"\tSpline 0.001 0.2")
		hudTournamentSetupPanelClose = append(hudTournamentSetupPanelClose, "}")

		fmt.Println("Did not find custom HudTournamentSetupPanelClose animation, using default")
	} else if foundHudTournamentSetupPanelClose == 1 {
		fmt.Println("Error: Found HudTournamentSetupPanelClose animation header, but did not close")
		pressToExit()
	}

	// If no custom HudReadyPulseEnd animation is found, use default code
	if foundHudReadyPulseEnd == 0 {
		hudReadyPulseEnd = append(hudReadyPulseEnd, "event HudReadyPulseEnd")
		hudReadyPulseEnd = append(hudReadyPulseEnd, "{")
		hudReadyPulseEnd = append(hudReadyPulseEnd, "\tAnimate TournamentInstructionsLabel\tFgColor\t\"TanLight\"\tLinear 0.0 0.1")
		hudReadyPulseEnd = append(hudReadyPulseEnd, "\tStopEvent HudReadyPulse\t\t0.0")
		hudReadyPulseEnd = append(hudReadyPulseEnd, "\tStopEvent HudReadyPulseLoop\t0.0")
		hudReadyPulseEnd = append(hudReadyPulseEnd, "}")

		fmt.Println("Did not find custom HudReadyPulseEnd animation, using default")
	} else if foundHudReadyPulseEnd == 1 {
		fmt.Println("Error: Found HudReadyPulseEnd animation header, but did not close")
		pressToExit()
	}

	return hintMessageHide, hudTournamentSetupPanelOpen, hudTournamentSetupPanelClose, hudReadyPulseEnd
}

func insertPeachRecAnimations(hintMessageHide []string, hudTournamentSetupPanelOpen []string, hudTournamentSetupPanelClose []string, hudReadyPulseEnd []string) []string {
	var peachRecAnimations []string
	// Insert animations for HintMessageHide
	for a, line := range hintMessageHide {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(hintMessageHide)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecSpawn 0.0")
		}
	}
	// Insert animations for HudTournamentSetupPanelOpen
	for a, line := range hudTournamentSetupPanelOpen {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(hudTournamentSetupPanelOpen)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecOpen 0.0")
		}
	}
	// Insert animations for HudTournamentSetupPanelClose
	for a, line := range hudTournamentSetupPanelClose {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(hudTournamentSetupPanelClose)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecClose 0.0")
		}
	}
	// Insert animations for HudReadyPulseEnd
	for a, line := range hudReadyPulseEnd {
		peachRecAnimations = append(peachRecAnimations, line)
		// Add PeachREC line just before closing
		if a == len(hudReadyPulseEnd)-2 {
			peachRecAnimations = append(peachRecAnimations, "\n\tRunEventChild MainMenuOverride PeachRecMvM 0.0 //Must run 3 times")
			peachRecAnimations = append(peachRecAnimations, "\tRunEventChild MainMenuOverride PeachRecMvM 0.0")
			peachRecAnimations = append(peachRecAnimations, "\tRunEventChild MainMenuOverride PeachRecMvM 0.0")
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

	peachRecAnimations = append(peachRecAnimations, "event PeachRecMvM")
	peachRecAnimations = append(peachRecAnimations, "{")
	peachRecAnimations = append(peachRecAnimations, "\tFireCommand 0.0 \"engine pr_mvm\"")
	peachRecAnimations = append(peachRecAnimations, "}")

	return peachRecAnimations
}
