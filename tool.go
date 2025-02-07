package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

var fileWriteSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"fileName": {
			Type:        genai.TypeString,
			Description: "The name of the file to write to. Do not include extension, it will be automatically added (.txt)",
		},
		"content": {
			Type:        genai.TypeString,
			Description: "The text content to write to the file",
		},
	},
	Required: []string{"fileName", "content"},
}

var FileTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "file_write",
			Description: "write a text file to user local file system with specified name and content.",
			Parameters:  fileWriteSchema,
		},
	},
}

func WriteDesktop(fileName string, content string) error {
	fileName = fileName + ".txt"
	home, _ := os.UserHomeDir()
	fullPath := filepath.Join(home, "Desktop", fileName)

	formattedContent := strings.ReplaceAll(content, "\\n", "\n")

	err := os.WriteFile(fullPath, []byte(formattedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func scanDirectory(dir string) (string, error) {
	var content strings.Builder

	// Load ignored patterns from .fileignore
	ignorePatterns, err := loadIgnorePatterns(dir)
	if err != nil {
		log.Printf("Warning: Failed to load .fileignore: %v\n", err)
	}

	// Walk through the directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Skip this file or folder
		}

		relPath, _ := filepath.Rel(dir, path) // Get relative path for ignoring

		// Skip ignored files and directories
		if isIgnored(relPath, info, ignorePatterns) {
			log.Printf("Skipping ignored file: %s\n", path)
			return nil
		}

		if info.IsDir() {
			content.WriteString(fmt.Sprintf("Folder: %s\n", path))
			return nil
		}

		// Read file content
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v\n", path, err)
			return nil
		}

		// Detect binary files (like images/videos)
		if !utf8.Valid(fileContent) {
			log.Printf("Processing binary file: %s\n", path)

			// Convert binary file to Base64
			base64Data := base64.StdEncoding.EncodeToString(fileContent)

			// Limit Base64 preview length (avoid large text)
			preview := base64Data
			if len(base64Data) > 1000 {
				preview = base64Data[:1000] + "..." // Show only first 1000 chars
			}

			// Format binary file info for Gemini
			content.WriteString(fmt.Sprintf(
				"Binary File: %s\nType: %s\nSize: %d bytes\nBase64 Preview:\n%s\n\n",
				path, detectFileType(path), info.Size(), preview,
			))

			return nil
		}

		// Append text file content
		content.WriteString(fmt.Sprintf("File: %s\n%s\n\n", path, string(fileContent)))
		return nil
	})

	return content.String(), err
}

// loadIgnorePatterns reads .fileignore and returns a slice of patterns
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

// isIgnored checks if a file matches any pattern in .fileignore
func isIgnored(path string, info os.FileInfo, patterns []string) bool {
	for _, pattern := range patterns {
		// Match directories
		if info.IsDir() && strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(path, strings.TrimSuffix(pattern, "/")) {
				return true
			}
		}

		// Match wildcards (e.g., *.jpg, *.gz)
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}

		// Match exact filename
		if pattern == path {
			return true
		}
	}
	return false
}

// detectFileType returns the file type based on the extension
func detectFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return "Image"
	case ".mp4", ".avi", ".mov", ".mkv":
		return "Video"
	case ".mp3", ".wav", ".ogg":
		return "Audio"
	case ".pdf":
		return "PDF Document"
	default:
		return "Unknown"
	}
}

var scanDirectorySchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"directory": {
			Type:        genai.TypeString,
			Description: "The directory path to scan. If empty, scans the current directory.",
		},
	},
	Required: []string{"directory"},
}

var ScanTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "scan_directory",
			Description: "Scans the specified directory and lists all files.",
			Parameters:  scanDirectorySchema,
		},
	},
}

//package main
//
//import (
//"fmt"
//"log"
//"os"
//"path/filepath"
//)

// ReadFile reads the content of a given file, assuming the current directory if only a filename is provided
func ReadFile(filePath string) (string, error) {
	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %v", err)
		}
		filePath = filepath.Join(cwd, filePath) // Ensure it resolves correctly
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(data), nil
}

var ReadFileSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"fileName": { // Change "directory" to "fileName"
			Type:        genai.TypeString,
			Description: "The name or path of the file to be read. If only a filename is provided, the file will be searched in the current working directory.",
		},
	},
	Required: []string{"fileName"}, // Change "directory" to "fileName"
}

var ReadFileTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "ReadFile",
			Description: "Reads the contents of a specified file and returns the text.",
			Parameters:  ReadFileSchema,
		},
	},
}

//func GetSystemInfo() (map[string]string, error) {
//	info := make(map[string]string)
//
//	// OS and CPU Architecture
//	info["os"] = runtime.GOOS
//	info["architecture"] = runtime.GOARCH
//
//	// Detect terminal shell
//	shell := os.Getenv("SHELL")
//	if shell == "" {
//		shell = "Unknown"
//	}
//	info["shell"] = shell
//
//	// Get CPU details
//	cpuInfo, err := getCommandOutput("lscpu") // Linux
//	if err != nil {
//		cpuInfo, err = getCommandOutput("wmic cpu get name") // Windows
//	}
//	if err == nil {
//		info["cpu"] = strings.TrimSpace(cpuInfo)
//	} else {
//		info["cpu"] = "Unknown"
//	}
//
//	// Get GPU details (NVIDIA only)
//	gpuInfo, err := getCommandOutput("nvidia-smi --query-gpu=name --format=csv,noheader")
//	if err == nil {
//		info["gpu"] = strings.TrimSpace(gpuInfo)
//	} else {
//		info["gpu"] = "No NVIDIA GPU detected"
//	}
//
//	return info, nil
//}
//
//// Helper function to run shell commands
//func getCommandOutput(cmdStr string) (string, error) {
//	cmdParts := strings.Split(cmdStr, " ")
//	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	err := cmd.Run()
//	if err != nil {
//		return "", err
//	}
//	return out.String(), nil
//}
//
//var systemInfoSchema = &genai.Schema{
//	Type:        genai.TypeObject,
//	Properties:  map[string]*genai.Schema{},
//	Description: "Retrieves system information including OS, CPU, GPU, and terminal shell.",
//}
//
//var SystemInfoTool = &genai.Tool{
//	FunctionDeclarations: []*genai.FunctionDeclaration{
//		{
//			Name:        "get_system_info",
//			Description: "Gets details about the operating system, CPU, GPU, and terminal shell.",
//			Parameters:  systemInfoSchema,
//		},
//	},
//}
