package main

import (
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var fileWriteSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"fileName": {
			Type: genai.TypeString,
			Description: "The name or path of the file to write to. " +
				"If only a file name is provided, the file will be created in the current working directory. " +
				"You may specify any file extension.",
		},
		"content": {
			Type:        genai.TypeString,
			Description: "The text content to write to the file.",
		},
	},
	Required: []string{"fileName", "content"},
}

var FileTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name: "file_write",
			Description: "Writes a file to the local file system with the specified name/path and content. " +
				"If the file already exists, it will be overwritten.",
			Parameters: fileWriteSchema,
		},
	},
}

func WriteDesktop(fileName string, content string) error {
	var fullPath string

	// If fileName is an absolute path or contains a directory separator, use it as-is.
	// Otherwise, assume it's just a file name and use the current working directory.
	if filepath.IsAbs(fileName) || strings.Contains(fileName, string(os.PathSeparator)) {
		fullPath = fileName
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		fullPath = filepath.Join(cwd, fileName)
	}

	// Replace literal "\n" with actual newlines in the content.
	formattedContent := strings.ReplaceAll(content, "\\n", "\n")

	// Write the file, using 0644 permissions.
	// This will override the file if it already exists.
	err := os.WriteFile(fullPath, []byte(formattedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file at '%s': %v", fullPath, err)
	}
	return nil
}

func scanDirectory(dir string) (string, error) {
	// Ensure the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory '%s' does not exist", dir)
	}

	var content strings.Builder
	var indentLevel int

	ignorePatterns, err := loadIgnorePatterns(dir)
	if err != nil {
		log.Printf("Warning: Failed to load .fileignore: %v\n", err)
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Skip this file or folder
		}

		relPath, _ := filepath.Rel(dir, path)

		if isIgnored(relPath, info, ignorePatterns) {
			log.Printf("Skipping ignored file: %s\n", path)
			return nil
		}

		indentLevel = strings.Count(relPath, string(os.PathSeparator))
		indentation := strings.Repeat("  ", indentLevel)

		if info.IsDir() {
			content.WriteString(fmt.Sprintf("%s📁 Folder: %s\n", indentation, path))
			return nil
		}

		content.WriteString(fmt.Sprintf("%s📄 File: %s (Size: %d bytes)\n", indentation, path, info.Size()))
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
			Description: "Reads the contents of a specified file and send it to you to analyze it.",
			Parameters:  ReadFileSchema,
		},
	},
}

//	func GetSystemInfo() (map[string]string, error) {
//		info := make(map[string]string)
//
//		// OS and CPU Architecture
//		info["os"] = runtime.GOOS
//		info["architecture"] = runtime.GOARCH
//
//		// Detect terminal shell
//		shell := os.Getenv("SHELL")
//		if shell == "" {
//			shell = "Unknown"
//		}
//		info["shell"] = shell
//
//		// Get CPU details
//		cpuInfo, err := getCommandOutput("lscpu") // Linux
//		if err != nil {
//			cpuInfo, err = getCommandOutput("wmic cpu get name") // Windows
//		}
//		if err == nil {
//			info["cpu"] = strings.TrimSpace(cpuInfo)
//		} else {
//			info["cpu"] = "Unknown"
//		}
//
//		// Get GPU details (NVIDIA only)
//		gpuInfo, err := getCommandOutput("nvidia-smi --query-gpu=name --format=csv,noheader")
//		if err == nil {
//			info["gpu"] = strings.TrimSpace(gpuInfo)
//		} else {
//			info["gpu"] = "No NVIDIA GPU detected"
//		}
//
//		return info, nil
//	}
//
// // Helper function to run shell commands
//
//	func getCommandOutput(cmdStr string) (string, error) {
//		cmdParts := strings.Split(cmdStr, " ")
//		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
//		var out bytes.Buffer
//		cmd.Stdout = &out
//		err := cmd.Run()
//		if err != nil {
//			return "", err
//		}
//		return out.String(), nil
//	}
//
//	var systemInfoSchema = &genai.Schema{
//		Type:        genai.TypeObject,
//		Properties:  map[string]*genai.Schema{},
//		Description: "Retrieves system information including OS, CPU, GPU, and terminal shell.",
//	}
//
//	var SystemInfoTool = &genai.Tool{
//		FunctionDeclarations: []*genai.FunctionDeclaration{
//			{
//				Name:        "get_system_info",
//				Description: "Gets details about the operating system, CPU, GPU, and terminal shell.",
//				Parameters:  systemInfoSchema,
//			},
//		},
//	}
func RunCommand(cmdLine string) (string, error) {
	parts := strings.Fields(cmdLine)
	if len(parts) == 0 {
		return "", fmt.Errorf("no command provided")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v\nOutput: %s", err, output)
	}
	outStr := string(output)
	// If output is empty, return a default success message.
	if strings.TrimSpace(outStr) == "" {
		outStr = "Command executed successfully, but no output was returned."
	}
	return outStr, nil
}

var runCommandSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"cmdLine": {
			Type:        genai.TypeString,
			Description: "The full terminal command to execute, including the command and all its arguments (e.g., 'ls -la /home/user').",
		},
	},
	Required: []string{"cmdLine"},
}

var RunCommandTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "run_command",
			Description: "Executes a simple terminal command  and returns its output.",
			Parameters:  runCommandSchema,
		},
	},
}
