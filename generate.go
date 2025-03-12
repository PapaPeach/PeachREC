package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func generateManifest(manifest string, peachRecManifest []string) {
	// Update existing manifest to include PeachREC
	file, err := os.Create(manifest)
	if err != nil {
		fmt.Println("Error generating hudanimations_manifest.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write peachRecManifest over hudanimations_manifest.txt
	for _, line := range peachRecManifest {
		file.WriteString(line + "\n")
	}
}

func generateAnimations(hud string, peachRecAnimations []string) {
	// Generate hudanimations_peachrec.txt
	file, err := os.Create(filepath.Join(hud, "hudanimations_peachrec.txt"))
	if err != nil {
		fmt.Println("Error generating hudanimations_peachrec.txt:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write peachRecAnimations to hudanimations_peachrec.txt
	for _, line := range peachRecAnimations {
		file.WriteString(line + "\n")
	}
}

func generateConfig(hud string) {
	// Generate cfg/peachrec.cfg
	filePath := filepath.Join(hud, "cfg")
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(filePath, os.ModePerm)
	}
	fileName := filepath.Join(filePath, "peachrec.cfg")
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error generating peachrec.cfg:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write peachrec.cfg
	file.WriteString("alias peachrec \"player_ready_toggle\"\n\n")

	file.WriteString("alias pr1 \"alias peachrec pr2;alias pr_open pr_open_test;player_ready_toggle\"\n")
	file.WriteString("alias pr2 \"pr_trigger\"\n")
	file.WriteString("alias pr_reset \"alias peachrec pr1\"\n")
	file.WriteString("alias pr_trigger \"alias peachrec;alias pr_close pr_end;ds_record;echo =====PeachREC.started.recording=====\"\n")
	file.WriteString("alias pr_end \"ds_stop;alias peachrec player_ready_toggle;alias pr_open pr_open_init;alias pr_close pr_close_init;echo =====PeachREC.stopped.recording=====\"\n\n")

	file.WriteString("alias pr_open \"pr_open_init\"\n")
	file.WriteString("alias pr_close \"pr_close_init\"\n\n")

	file.WriteString("alias pr_open_init \"alias peachrec pr1;alias pr_open pr_open_nat;alias pr_close pr_close_nat;player_ready_toggle;echo =====PeachREC.waiting.for.match.to.start=====\"\n")
	file.WriteString("alias pr_close_init \"player_ready_toggle\"\n\n")

	file.WriteString("alias pr_open_nat \"alias pr_close pr_close_nat\"\n")
	file.WriteString("alias pr_close_nat \"alias pr_close pr_close_newserver\"\n\n")

	file.WriteString("alias pr_open_test \"pr_reset;alias pr_open pr_open_nat;alias pr_close pr_close_nat;player_ready_toggle\"\n\n")

	file.WriteString("alias pr_close_newserver \"alias peachrec player_ready_toggle;alias pr_open pr_open_init;alias pr_close pr_close_init;echo =====PeachREC.detected.new.server=====\"\n")
	file.WriteString("alias pr2_newserver \"alias pr2 pr_trigger;echo =====PeachREC.detected.new.match.server=====\"\n\n")

	file.WriteString("echo ===============\necho PeachREC Active\necho ===============")
}
