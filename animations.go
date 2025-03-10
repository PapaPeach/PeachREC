package main

import (
	"fmt"
	"os"
	"path/filepath"
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
