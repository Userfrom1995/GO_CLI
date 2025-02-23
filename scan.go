package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// scanDirectory scans the given directory using `tree` (or `find/dir`) while respecting .fileignore
func scanDirectory(dir string) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("error: directory path is required")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("error: failed to get absolute path: %v", err)
	}

	// Ensure the directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("error: directory '%s' does not exist", absPath)
	}

	// Load ignore patterns (optional)
	ignorePatterns, err := loadIgnorePatterns(absPath)
	if err != nil {
		log.Printf("Warning: Failed to load .fileignore: %v\n", err)
	}

	// Run tree command with ignore rules
	treeOutput, err := runTreeCommand(absPath, ignorePatterns)
	if err != nil {
		return "", fmt.Errorf("error: failed to run tree command: %v", err)
	}

	return treeOutput, nil
}

// runTreeCommand executes the system's directory listing command while skipping ignored files
func runTreeCommand(dir string, ignorePatterns []string) (string, error) {
	var cmd *exec.Cmd
	var ignoreArgs []string

	// Convert ignore patterns to command arguments
	for _, pattern := range ignorePatterns {
		if runtime.GOOS == "windows" {
			ignoreArgs = append(ignoreArgs, fmt.Sprintf("/S /B /A:-D | find /V \"%s\"", pattern))
		} else {
			ignoreArgs = append(ignoreArgs, fmt.Sprintf("! -path '%s/*'", filepath.Join(dir, pattern)))
		}
	}

	if runtime.GOOS == "windows" {
		// Windows: Use `dir /s /b` and filter ignored files
		cmd = exec.Command("cmd", "/c", "dir /s /b", dir, strings.Join(ignoreArgs, " "))
	} else {
		// Linux/macOS: Use `find` with ignore filters
		if _, err := exec.LookPath("tree"); err == nil {
			cmd = exec.Command("tree", "-a", "--noreport", "--prune", "--matchdirs", dir)
			for _, pattern := range ignorePatterns {
				cmd.Args = append(cmd.Args, "-I", pattern)
			}
		} else {
			findCmd := []string{"find", dir, "-type", "f"}
			findCmd = append(findCmd, ignoreArgs...)
			cmd = exec.Command(findCmd[0], findCmd[1:]...)
		}
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	return out.String(), err
}

func loadIgnorePatterns(dir string) ([]string, error) {
	ignoreFilePath := filepath.Join(dir, ".fileignore")
	data, err := ioutil.ReadFile(ignoreFilePath)
	if err != nil {
		return nil, err // No .fileignore found, return empty slice
	}

	lines := strings.Split(string(data), "\n")
	var patterns []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			patterns = append(patterns, trimmed)
		}
	}
	return patterns, nil
}
