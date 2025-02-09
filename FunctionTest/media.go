package main

import (
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"os/exec"
	"strings"
	// "google.golang.org/api/option"
)

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
		"command": {
			Type:        genai.TypeString,
			Description: "The terminal command to execute. Examples include 'ls', 'rm', 'mv', 'cp', etc.",
		},
		"args": {
			Type:        genai.TypeString,
			Description: " space-separated arguments for the command.",
		},
	},
	Required: []string{"command"},
}

var RunCommandTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "run_command",
			Description: "Executes a simple terminal command with optional arguments and returns its output.",
			Parameters:  runCommandSchema,
		},
	},
}
