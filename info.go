package main

import (
	"bytes"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %v\nOutput: %s", err, out.String())
	}
	return out.String(), nil
}

func GetSystemSpecs() (map[string]string, error) {
	specs := make(map[string]string)

	// Basic information
	specs["os"] = runtime.GOOS
	specs["architecture"] = runtime.GOARCH

	// Determine shell/command interpreter.
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = os.Getenv("COMSPEC")
	}
	if shell == "" {
		shell = "unknown"
	}
	specs["shell"] = shell

	// OS Version details:
	if runtime.GOOS == "linux" {
		// Read /etc/os-release for Linux distribution details.
		data, err := ioutil.ReadFile("/etc/os-release")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					value := strings.TrimPrefix(line, "PRETTY_NAME=")
					value = strings.Trim(value, `"`)
					specs["os_version"] = value
					break
				}
			}
		} else {
			specs["os_version"] = "unknown"
		}
	} else if runtime.GOOS == "windows" {
		// Use wmic to get the Windows OS caption.
		output, err := runCommand("wmic", "os", "get", "Caption")
		if err == nil {
			// The output typically contains a header and then the OS name.
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) >= 2 {
				specs["os_version"] = strings.TrimSpace(lines[1])
			} else {
				specs["os_version"] = "unknown"
			}
		} else {
			specs["os_version"] = "unknown"
		}
	} else {
		specs["os_version"] = "unknown"
	}

	// CPU information:
	if runtime.GOOS == "linux" {
		data, err := ioutil.ReadFile("/proc/cpuinfo")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "model name") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 {
						specs["cpu"] = strings.TrimSpace(parts[1])
						break
					}
				}
			}
		} else {
			specs["cpu"] = "unknown"
		}
	} else if runtime.GOOS == "windows" {
		output, err := runCommand("wmic", "cpu", "get", "Name")
		if err == nil {
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) >= 2 {
				specs["cpu"] = strings.TrimSpace(lines[1])
			} else {
				specs["cpu"] = "unknown"
			}
		} else {
			specs["cpu"] = "unknown"
		}
	} else {
		specs["cpu"] = "unknown"
	}

	// GPU information:
	if runtime.GOOS == "linux" {
		// First, try using nvidia-smi for NVIDIA GPUs.
		output, err := runCommand("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
		if err == nil && strings.TrimSpace(output) != "" {
			specs["gpu"] = strings.TrimSpace(output)
		} else {
			// Fallback: Use lspci to search for VGA or 3D controllers.
			output, err := runCommand("lspci")
			if err == nil {
				lines := strings.Split(output, "\n")
				var gpuInfo string
				for _, line := range lines {
					lower := strings.ToLower(line)
					if strings.Contains(lower, "vga controller") || strings.Contains(lower, "3d controller") {
						gpuInfo = line
						break
					}
				}
				if gpuInfo != "" {
					specs["gpu"] = gpuInfo
				} else {
					specs["gpu"] = "no dedicated GPU detected"
				}
			} else {
				specs["gpu"] = "no dedicated GPU detected"
			}
		}
	} else if runtime.GOOS == "windows" {
		output, err := runCommand("wmic", "path", "win32_VideoController", "get", "Name")
		if err == nil {
			lines := strings.Split(strings.TrimSpace(output), "\n")
			if len(lines) >= 2 {
				specs["gpu"] = strings.TrimSpace(lines[1])
			} else {
				specs["gpu"] = "no dedicated GPU detected"
			}
		} else {
			specs["gpu"] = "no dedicated GPU detected"
		}
	} else {
		specs["gpu"] = "unknown"
	}

	return specs, nil
}

var systemInfoSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"os": {
			Type:        genai.TypeString,
			Description: "The operating system identifier (e.g., 'linux', 'windows').",
		},
		"os_version": {
			Type:        genai.TypeString,
			Description: "The full operating system version or distribution details (e.g., 'Ubuntu 22.04 LTS' or 'Microsoft Windows 10 Pro').",
		},
		"architecture": {
			Type:        genai.TypeString,
			Description: "The CPU architecture (e.g., 'amd64', 'arm64').",
		},
		"shell": {
			Type:        genai.TypeString,
			Description: "The default terminal shell or command interpreter (e.g., '/bin/bash' or 'cmd.exe').",
		},
		"cpu": {
			Type:        genai.TypeString,
			Description: "Detailed CPU model information.",
		},
		"gpu": {
			Type:        genai.TypeString,
			Description: "Detailed GPU information if available, or a note if no dedicated GPU is detected.",
		},
	},
	Description: "Returns system specifications including OS, OS version, architecture, shell, CPU, and GPU information.",
}

var SystemInfoTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name: "get_system_info",
			Description: "Retrieves detailed system specifications, including operating system and its version, architecture, terminal shell, CPU, and GPU information. " +
				"This helps the assistant choose and tailor system commands appropriately.",
			Parameters: systemInfoSchema,
		},
	},
}
