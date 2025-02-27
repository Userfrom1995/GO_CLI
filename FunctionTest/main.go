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
	"strconv"
	"strings"
)

// scanDirectory scans the given directory while respecting .fileignore and the depth limit.
func scanDirectory(dir string, depth int) (string, error) {
	if dir == "" {
		return "", fmt.Errorf("error: directory path is required")
	}

	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("error: failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("error: directory '%s' does not exist", absPath)
	}

	ignorePatterns, err := loadIgnorePatterns(absPath)
	if err != nil {
		log.Printf("Warning: Failed to load .fileignore: %v\n", err)
		ignorePatterns = []string{}
	}

	output, err := runDirectoryCommand(absPath, ignorePatterns, depth)
	if err != nil {
		return "", fmt.Errorf("error: failed to list directory: %v", err)
	}

	return output, nil
}

// runDirectoryCommand executes platform-specific directory listing with depth control.
func runDirectoryCommand(dir string, ignorePatterns []string, depth int) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// Use PowerShell for proper depth control
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf(`Get-ChildItem -Path "%s" -Recurse -Depth %d -Name`, dir, depth))
	} else {
		// Use `find` for Unix-like systems
		cmd = exec.Command("find", dir, "-maxdepth", strconv.Itoa(depth), "-print")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	lines := strings.Split(out.String(), "\n")
	filteredLines := filterIgnoredFiles(lines, ignorePatterns)
	return strings.Join(filteredLines, "\n"), nil
}

// loadIgnorePatterns reads .fileignore and returns a slice of patterns
func loadIgnorePatterns(dir string) ([]string, error) {
	ignoreFilePath := filepath.Join(dir, ".fileignore")
	data, err := ioutil.ReadFile(ignoreFilePath)
	if os.IsNotExist(err) {
		return []string{}, nil // No .fileignore found
	} else if err != nil {
		return nil, err
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

// filterIgnoredFiles removes paths that match .fileignore patterns
func filterIgnoredFiles(files []string, ignorePatterns []string) []string {
	var filtered []string
	for _, file := range files {
		ignored := false
		for _, pattern := range ignorePatterns {
			if strings.Contains(file, pattern) {
				ignored = true
				break
			}
		}
		if !ignored {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func main() {
	depth := 2 // Default depth
	if len(os.Args) > 1 {
		userDepth, err := strconv.Atoi(os.Args[1])
		if err == nil {
			depth = userDepth
		}
	}

	output, err := scanDirectory("/mnt", depth)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println(output)
}
