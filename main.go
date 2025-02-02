package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
	// "os/exec"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

//var cs *genai.ChatSession

func main() {
	// Load ..env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No ..env file found, falling back to system environment variables")
	}

	// Get the API key
	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		log.Fatal("GEMINI_API_KEY is not set. Please provide it in the ..env file or as an environment variable.")
	}

	// Create a new GenAI client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
	if err != nil {
		log.Fatalf("Failed to create GenAI client: %v", err)
	}
	defer client.Close()

	// Specify the directory to scan
	dir := "./" // Current directory
	content, err := scanDirectory(dir)
	if err != nil {
		log.Fatalf("Error scanning directory: %v", err)
	}

	// Initialize chat context
	cs := initializeChat(client, ctx, content)

	// Start interactive conversation
	fmt.Println("Context sent! You can now chat with Gemini.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ") // Prompt for user input
		if !scanner.Scan() {
			break
		}
		userInput := scanner.Text()

		if strings.ToLower(userInput) == "update" || strings.ToLower(userInput) == "relook" {
			fmt.Println("Re-scanning directory...")
			content, err := scanDirectory(dir)
			if err != nil {
				log.Printf("Error scanning directory: %v\n", err)
				continue
			}
			updateChatContext(cs, content)
			fmt.Println("Directory context updated!")
			continue
		}

		// Send user message to Gemini
		iter := cs.SendMessageStream(ctx, genai.Text(userInput))
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("Error sending message: %v\n", err)
				break
			}

			functionCallHandled := false
			for _, part := range resp.Candidates[0].Content.Parts {
				// Use a type assertion to check if this part is a function call.
				if fc, ok := part.(genai.FunctionCall); ok {
					fmt.Printf("Received function call: %s\n", fc.Name)
					var apiResult map[string]any

					// Handle the function call based on its name.
					switch fc.Name {
					case "getCurrentTime":
						apiResult = map[string]any{
							"currentTime": getCurrentTime(),
						}
					// Add more cases if you register more tools/functions.
					default:
						apiResult = map[string]any{
							"result": "Function not implemented",
						}
					}

					fmt.Printf("Sending API result:\n%v\n\n", apiResult)
					// Send the API result back to Gemini as a function response.
					//resp, err = cs.SendMessage(ctx, genai.FunctionResponse{
					//	Name:     fc.Name,
					//	Response: apiResult,
					//})
					//if err != nil {
					//	log.Printf("Error sending function response: %v\n", err)
					//} else {
					//	printResponse(resp)
					//}
					funcResp := &genai.Content{
						Parts: []genai.Part{
							genai.FunctionResponse{
								Name:     fc.Name,
								Response: apiResult,
							},
						},
						Role: "function",
					}

					// Append to chat history
					cs.History = append(cs.History, funcResp)

					// Get model's response to the function result
					iter := cs.SendMessageStream(ctx) // No message needed - uses history
					for {
						resp, err := iter.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							log.Printf("Error getting response: %v", err)
							break
						}
						printResponse(resp)

						// 4. Add model's response to history
						if len(resp.Candidates) > 0 {
							cs.History = append(cs.History, resp.Candidates[0].Content)
						}
					}
					functionCallHandled = true
					break // Exit the loop after handling the function call.
				}
			}

			// If no function call was detected, simply print out Gemini's response.
			if !functionCallHandled {
				printResponse(resp)
			}
		}
	}
	// Check if Gemini requests a function call
	//if functionName := detectFunctionCall(resp); functionName != "" {
	//	executeAndRespondToFunction(ctx, cs, functionName, dir)
	//	break
	//}
	// Print Gemini's response
	//printResponse(resp)

	// if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
	// 	// Iterate over each part to extract text
	// 	var commandText string
	// 	for _, part := range resp.Candidates[0].Content.Parts {
	// 		// Check if the part is of type genai.Text
	// 		switch p := part.(type) {
	// 		case genai.Text:
	// 			// Extract text from genai.Text and append it to commandText
	// 			commandText = string(p)
	// 		case *genai.Text:
	// 			// If it's a pointer to genai.Text, dereference it and append to commandText
	// 			commandText = string(*p)
	// 		}
	// 	}

	// // Check if the commandText contains "run command:"
	// if strings.Contains(commandText, "run command:") {
	// 	// Extract the actual command after "run command:"
	// 	command := strings.TrimSpace(strings.TrimPrefix(commandText, "run command:"))

	// 	// Execute the command Gemini suggested
	// 	output, err := executeCommand(command)
	// 	if err != nil {
	// 		fmt.Printf("Error executing command: %v\n", err)
	// 	}

	// 	// Send the output back to Gemini
	// 	_, sendErr := cs.SendMessage(ctx, genai.Text(output))
	// 	if sendErr != nil {
	// 		fmt.Printf("Error sending message to Gemini: %v\n", sendErr)
	// 	}
	// }
	//	}
	//
	//}
}

// detectFunctionCall checks if Gemini is requesting a function execution
//func detectFunctionCall(resp *genai.GenerateContentResponse) string {
//	for _, part := range resp.Candidates[0].Content.Parts {
//		text := fmt.Sprintf("%v", part) // Convert response part to text
//
//		// Check if the response asks for a function
//		if strings.Contains(text, "Call: scanDirectory") {
//			return "scanDirectory"
//		} else if strings.Contains(text, "Call: fetchSystemInfo") {
//			return "fetchSystemInfo"
//		}
//	}
//	return "" // No function call detected
//}

// func executeCommand(cmdStr string) (string, error) {
// 	// Prepare the command
// 	cmd := exec.Command("bash", "-c", cmdStr)

// 	// Create pipes for stdout and stderr
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return "", fmt.Errorf("error creating stdout pipe: %v", err)
// 	}
// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		return "", fmt.Errorf("error creating stderr pipe: %v", err)
// 	}

// 	// Start the command
// 	if err := cmd.Start(); err != nil {
// 		return "", fmt.Errorf("error starting command: %v", err)
// 	}

// 	// Create readers for stdout and stderr
// 	// stdoutReader := bufio.NewReader(stdout)
// 	// stderrReader := bufio.NewReader(stderr)

// 	// Channels to capture the complete output
// 	outputCh := make(chan string)
// 	errorCh := make(chan string)

// 	// Stream stdout and stderr in separate goroutines
// 	// go streamOutput(stdoutReader, "STDOUT", outputCh)
// 	// go streamOutput(stderrReader, "STDERR", errorCh)

// 	// Collect the output from both channels
// 	var stdoutOutput, stderrOutput string
// 	go func() {
// 		for line := range outputCh {
// 			stdoutOutput += line
// 		}
// 	}()
// 	go func() {
// 		for line := range errorCh {
// 			stderrOutput += line
// 		}
// 	}()

// 	// Wait for the command to finish
// 	if err := cmd.Wait(); err != nil {
// 		return fmt.Sprintf("Command failed:\nSTDOUT:\n%s\nSTDERR:\n%s\n", stdoutOutput, stderrOutput), err
// 	}

// 	// Combine stdout and stderr
// 	combinedOutput := fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s\n", stdoutOutput, stderrOutput)
// 	return combinedOutput, nil
// }

// streamOutput reads and prints output line-by-line in real-time
//
//	func streamOutput(reader *bufio.Reader, label string, ch chan string) {
//		defer close(ch)
//		for {
//			line, err := reader.ReadString('\n')
//			if len(line) > 0 {
//				fmt.Printf("[%s] %s", label, line) // Print to user terminal
//				ch <- line                        // Send line to output channel
//			}
//			if err != nil {
//				break
//			}
//		}
//	}
func getCurrentTime() string {
	return time.Now().Format(time.RFC1123)
}

// executeAndRespondToFunction runs the requested function and sends the result to Gemini
//func executeAndRespondToFunction(ctx context.Context, cs *genai.ChatSession, functionName string, dir string) {
//	var response string
//
//	switch functionName {
//	case "scanDirectory":
//		fmt.Println("🔍 Gemini requested: scanDirectory()")
//		content, err := scanDirectory(dir)
//		if err != nil {
//			response = fmt.Sprintf("Error scanning directory: %v", err)
//		} else {
//			response = fmt.Sprintf("Scan result:\n%s", content)
//		}
//
//	case "fetchSystemInfo":
//		fmt.Println("📡 Gemini requested: fetchSystemInfo()")
//		response = fetchSystemInfo()
//
//	default:
//		response = "Error: Unknown function requested."
//	}
//
//	// Send function result back to Gemini
//	_, err := cs.SendMessage(ctx, genai.Text(response))
//	if err != nil {
//		log.Printf("Error sending function result to Gemini: %v\n", err)
//	}
//}

// fetchSystemInfo gathers basic system info
func fetchSystemInfo() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("System Info:\nHostname: %s\nOS: %s", hostname, os.Getenv("OSTYPE"))
}

// initializeChat sets up the chat with initial file content
func initializeChat(client *genai.Client, ctx context.Context, content string) *genai.ChatSession {
	timeFunction := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{{
			Name:        "getCurrentTime",
			Description: "Returns the current date and time in RFC1123 format.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"currentTime": {
						Type:        genai.TypeString,
						Description: "The current date and time in RFC1123 format.",
					},
				},
				Required: []string{"currentTime"},
			},
		}},
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.Tools = []*genai.Tool{timeFunction}
	cs := model.StartChat()
	chatPrompt := fmt.Sprintf(`
	Here is the content of my files and folders:

	%s
    ### Available Functions:
	1. **scanDirectory** - Scan the current directory and list all files.
	2. **fetchSystemInfo** - Get system information (hostname, OS type).
	
	### How to Use:
	- If you need to scan the directory, respond with: "Call: scanDirectory".
	- If you need system info, respond with: "Call: fetchSystemInfo".

	You are the best copilot in world.
    Help the user to achieve most productivity
     Ask user for a nickname for you and follow that name tiil converstion ends
     Use the provided context and files to help the user`, content)

	// Initialize the chat with directory contents and command execution explanation
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(chatPrompt),
			},
			Role: "user",
		},
	}

	return cs
}

// updateChatContext appends new file content to the existing chat history
func updateChatContext(cs *genai.ChatSession, content string) {
	cs.History = append(cs.History, &genai.Content{
		Parts: []genai.Part{
			genai.Text(fmt.Sprintf("Here is the updated content of my files and folders:\n\n%s", content)),
		},
		Role: "user",
	})
}

// scanDirectory recursively scans a directory and reads file contents
// scanDirectory scans a directory, applies ignore rules, and converts binary files to Base64
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

// printResponse formats and prints Gemini's responses
func printResponse(resp *genai.GenerateContentResponse) {
	for _, part := range resp.Candidates[0].Content.Parts {
		// Assuming part is a string or implements a method to get the string content
		formatted := cleanMarkdown(fmt.Sprintf("%v", part)) // Convert part to string safely
		fmt.Println(formatted)                              // Print the cleaned response text
	}
}

// func printResponse(resp genai.GenerateContentResponse) {
// 	for _, part := range resp.Candidates[0].Content.Parts {
// 			textContent := ""
// 			switch p := part.(type) {
// 			case genai.Text: // Or whatever the actual type is
// 					textContent = p
// 			case *genai.Text: // Or a pointer to the type
// 					textContent = p
// 			default:
// 					log.Printf("Unexpected part type : %T\n", part)
// 					continue // Skip this part if it's not the expected type
// 			}
// 			formatted := cleanMarkdown(textContent)
// 			fmt.Println(formatted)
// 	}
// }

// cleanMarkdown removes Markdown artifacts for better readability
func cleanMarkdown(input string) string {
	// Replace Markdown-specific symbols
	replacer := strings.NewReplacer(
		"**", "", // Bold
		"*", "", // Italic
		"- ", "", // List items
		"`", "", // Inline code
	)
	output := replacer.Replace(input)

	// Remove extra newlines
	output = strings.TrimSpace(output)
	return output
}
