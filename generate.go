package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed embed/peachrec.cfg
var cfg string

//go:embed embed/hudanimations_manifest.txt
var hudanimationManifest string

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
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error generating hudanimations_manifest.txt:", err)
			pressToExit()
		}
	}

	fmt.Println("Added PeachREC to hudanimation_manifest.txt")
}

func generateDefaultManifest(workingDir string) {
	// Check that _PeachREC directory exists
	filePath := filepath.Join(workingDir, modName)
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating peachrec directory:", err)
			pressToExit()
		}
	}

	// Check that _PeachREC/scripts directory exists
	filePath = filepath.Join(filePath, "scripts")
	_, er := os.Stat(filePath)
	if errors.Is(er, os.ErrNotExist) { // If scripts does not exist, create one
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating peachrec/scripts directory:", err)
			pressToExit()
		}
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
	_, err = file.WriteString(hudanimationManifest)
	if err != nil {
		fmt.Println("Error generating hudanimations_manifest.txt with default code:", err)
		pressToExit()
	}

	fmt.Println("Created", modName+"/scripts/hudanimations_manifest.txt")
}

func generateAnimations(workingDir string, peachRecAnimations []string) {
	// Generate hudanimations_peachrec.txt
	filePath := filepath.Join(workingDir, modName)
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating peachrec directory:", err)
			pressToExit()
		}
	}

	file, err := os.Create(filepath.Join(filePath, "hudanimations_peachrec.txt"))
	if err != nil {
		fmt.Println("Error generating hudanimations_peachrec.txt:", err)
		pressToExit()
	}
	defer file.Close()

	// Write peachRecAnimations to hudanimations_peachrec.txt
	for _, line := range peachRecAnimations {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error generating hudanimations_peachrec.txt:", err)
			pressToExit()
		}
	}

	fmt.Println("Created", modName+"/hudanimations_peachrec.txt")
}

func generateConfig(workingDir string) {
	// Check that _PeachREC directory exists
	filePath := filepath.Join(workingDir, modName)
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) { // If _PeachREC does not exist, create one
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating peachrec directory:", err)
			pressToExit()
		}
	}

	// Check that _PeachREC/cfg directory exists
	filePath = filepath.Join(filePath, "cfg")
	_, er := os.Stat(filePath)
	if errors.Is(er, os.ErrNotExist) { // If cfg does not exist, create one
		err = os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating peachrec/cfg directory:", err)
			pressToExit()
		}
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
	_, err = file.WriteString(cfg)
	if err != nil {
		fmt.Println("Error generating peachrec.cfg:", err)
		pressToExit()
	}

	fmt.Println("Created", modName+"/cfg/peachrec.cfg")
}

func generateAutoexec(file string) {
	input, err := os.Open(file)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(filepath.Dir(file), os.ModePerm)
		if err != nil {
			fmt.Println("Error creating autoexec directory:", err)
			pressToExit()
		}
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
		_, err = output.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error generating autoexec:", err)
			pressToExit()
		}
	}

	// Insert PeachREC exec
	_, err = output.WriteString("exec peachrec")
	if err != nil {
		fmt.Println("Error generating autoexec:", err)
		pressToExit()
	}

	fmt.Println("Added PeachREC to autoexec.cfg")
}
