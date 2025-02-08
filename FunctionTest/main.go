package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// RunCommand is a helper to run a command and return its combined output.
func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command failed: %v\nOutput: %s", err, out.String())
	}
	return out.String(), nil
}

// GetSystemSpecs collects minimal system specifications, including OS version and GPU info.
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
		output, err := RunCommand("wmic", "os", "get", "Caption")
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
		output, err := RunCommand("wmic", "cpu", "get", "Name")
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
		output, err := RunCommand("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
		if err == nil && strings.TrimSpace(output) != "" {
			specs["gpu"] = strings.TrimSpace(output)
		} else {
			// Fallback: Use lspci to search for VGA or 3D controllers.
			output, err := RunCommand("lspci")
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
		output, err := RunCommand("wmic", "path", "win32_VideoController", "get", "Name")
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

func main() {
	specs, err := GetSystemSpecs()
	if err != nil {
		fmt.Printf("Error retrieving system specs: %v\n", err)
		return
	}
	fmt.Println("System Specifications:")
	for key, value := range specs {
		fmt.Printf("%s: %s\n", key, value)
	}
}
