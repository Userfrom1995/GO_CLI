package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getAPIKey(reader *bufio.Reader) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	envFile := filepath.Join(homeDir, EnvFilePath)

	// Check if the key is already stored
	if key, err := os.ReadFile(envFile); err == nil {
		return strings.TrimSpace(string(key))
	}

	// If not found, ask the user
	fmt.Print("Enter your Gemini API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Save the key for future use
	err = os.WriteFile(envFile, []byte(apiKey), 0600) // Secure file permissions
	if err != nil {
		log.Fatalf("Error saving API key: %v", err)
	}

	fmt.Printf("API Key saved at %s\n", envFile)
	return apiKey
}
