package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generateManifest(manifest string, peachRecManifest []string) {
	// Update existing manifest to include PeachREC
	file, err := os.Create(manifest)
	if err != nil {
		fmt.Println("Error generating hudanimations_manifest.txt:", err)
		pressToExit()
	}
	defer file.Close()

	// Write peachRecManifest over hudanimations_manifest.txt
	for _, line := range peachRecManifest {
		file.WriteString(line + "\n")
	}

	fmt.Println("Added PeachREC to hudanimation_manifest.txt")
}

func generateDefaultManifest(workingDir string) {
	// Check that _PeachREC directory exists
	filePath := filepath.Join(workingDir, "_PeachREC")
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		os.Mkdir(filePath, os.ModePerm)
	}

	// Check that _PeachREC/scripts directory exists
	filePath = filepath.Join(filePath, "scripts")
	_, er := os.Stat(filePath)
	if errors.Is(er, os.ErrNotExist) { // If scripts does not exist, create one
		os.Mkdir(filePath, os.ModePerm)
	}

	// Generate hudanimations_manifest using default code
	fileName := filepath.Join(filePath, "hudanimations_manifest.txt")
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error generating hudanimations_manifest.txt with default code:", err)
		pressToExit()
	}
	defer file.Close()

	// Write hudanimations_manifest.txt default code
	file.WriteString("hudanimations_manifest\n")
	file.WriteString("{\n")
	file.WriteString("\tfile\tscripts/hudanimations_tf.txt\n")
	file.WriteString("\tfile\tscripts/hudanimations.txt\n")
	file.WriteString("}")

	fmt.Println("Created _PeachREC/scripts/hudanimations_manifest.txt")
}

func generateAnimations(workingDir string, peachRecAnimations []string) {
	// Generate hudanimations_peachrec.txt
	filePath := filepath.Join(workingDir, "_PeachREC")
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		os.Mkdir(filePath, os.ModePerm)
	}

	file, err := os.Create(filepath.Join(filePath, "hudanimations_peachrec.txt"))
	if err != nil {
		fmt.Println("Error generating hudanimations_peachrec.txt:", err)
		pressToExit()
	}
	defer file.Close()

	// Write peachRecAnimations to hudanimations_peachrec.txt
	for _, line := range peachRecAnimations {
		file.WriteString(line + "\n")
	}

	fmt.Println("Created _PeachREC/hudanimations_peachrec.txt")
}

func generateConfig(workingDir string) {
	// Check that _PeachREC directory exists
	filePath := filepath.Join(workingDir, "_PeachREC")
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		os.Mkdir(filePath, os.ModePerm)
	}

	// Check that _PeachREC/cfg directory exists
	filePath = filepath.Join(filePath, "cfg")
	_, er := os.Stat(filePath)
	if errors.Is(er, os.ErrNotExist) { // If cfg does not exist, create one
		os.Mkdir(filePath, os.ModePerm)
	}

	// Generate cfg/peachrec.cfg
	fileName := filepath.Join(filePath, "peachrec.cfg")
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error generating peachrec.cfg:", err)
		pressToExit()
	}
	defer file.Close()

	// Write peachrec.cfg
	file.WriteString("// Edit these to customize the message in chat when recording is started or stopped\n")
	file.WriteString("alias pr_message_start \"say_team [PeachREC] Recording...\"\n")
	file.WriteString("alias pr_message_stop \"say_team [PeachREC] Stop record.\"\n\n")

	file.WriteString("// Edit these to customize the sound when recording is on standby, started, or stopped\n")
	file.WriteString("alias pr_sound_standby \"play ui/mm_xp_chime.wav\"\n")
	file.WriteString("alias pr_sound_start \"play buttons/button17.wav\"\n")
	file.WriteString("alias pr_sound_stop \"play buttons/button17.wav\"\n\n\n\n")

	file.WriteString("//============================//\n")
	file.WriteString("//  PeachREC script contents  //\n")
	file.WriteString("//   no reason to edit this   //\n")
	file.WriteString("//============================//\n")
	file.WriteString("alias peachrec \"player_ready_toggle\"\n\n")

	file.WriteString("alias pr1 \"alias peachrec pr2;alias pr_open pr_open_test;player_ready_toggle\"\n")
	file.WriteString("alias pr2 \"pr_trigger\"\n")
	file.WriteString("alias pr_reset \"alias peachrec pr1\"\n")
	file.WriteString("alias pr_trigger \"alias peachrec;alias pr_close pr_end;pr_sound_start;pr_message_start;ds_record;echo =====PeachREC.started.recording=====\"\n")
	file.WriteString("alias pr_end \"ds_stop;alias peachrec player_ready_toggle;alias pr_open pr_open_init;alias pr_close pr_close_init;pr_sound_stop;pr_message_stop;echo =====PeachREC.stopped.recording=====\"\n\n")

	file.WriteString("alias pr_open \"pr_open_init\"\n")
	file.WriteString("alias pr_close \"pr_close_init\"\n\n")

	file.WriteString("alias pr_open_init \"alias peachrec pr1;alias pr_open pr_open_nat;alias pr_close pr_close_nat;player_ready_toggle;pr_sound_standby;echo =====PeachREC.waiting.for.match.to.start=====\"\n")
	file.WriteString("alias pr_close_init \"player_ready_toggle\"\n\n")

	file.WriteString("alias pr_open_nat \"alias pr_close pr_close_nat\"\n")
	file.WriteString("alias pr_close_nat \"alias pr_close pr_close_newserver\"\n\n")

	file.WriteString("alias pr_open_test \"pr_reset;alias pr_open pr_open_nat;alias pr_close pr_close_nat;player_ready_toggle\"\n\n")

	file.WriteString("alias pr_close_newserver \"alias peachrec player_ready_toggle;alias pr_open pr_open_init;alias pr_close pr_close_init;echo =====PeachREC.detected.new.server=====\"\n")
	file.WriteString("alias pr2_newserver \"alias pr2 pr_trigger;echo =====PeachREC.detected.new.match.server=====\"\n")
	file.WriteString("alias pr_mvm \"alias peachrec;echo =====PeachREC.detected.MvM.server=====\"\n\n")

	file.WriteString("echo\necho ===============\necho PeachREC Active\necho ===============\necho")

	fmt.Println("Created _PeachREC/cfg/peachrec.cfg")
}

func generateAutoexec(file string) {
	input, err := os.Open(file)
	if errors.Is(err, os.ErrNotExist) {
		os.Mkdir(filepath.Dir(file), os.ModePerm)
	} else if err != nil {
		fmt.Printf("Error opening %v to add PeachREC to autoexec: %v\n", file, err)
		pressToExit()
	}

	// Copy contents
	var contents []string
	scnr := bufio.NewScanner(input)
	for scnr.Scan() {
		line := scnr.Text()
		if !strings.Contains(line, "exec peachrec") {
			contents = append(contents, line)
		}
	}

	// Rewrite contents to file
	output, err := os.Create(file)
	if err != nil {
		fmt.Println("Error generating autoexec:", err)
		pressToExit()
	}
	defer output.Close()

	for _, line := range contents {
		output.WriteString(line + "\n")
	}

	// Insert PeachREC exec
	output.WriteString("exec peachrec")

	fmt.Println("Added PeachREC to autoexec.cfg")
}
