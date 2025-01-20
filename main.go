package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
	// "os/exec"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)
var cs *genai.ChatSession
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, falling back to system environment variables")
	}

	// Get the API key
	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
	if GEMINI_API_KEY == "" {
		log.Fatal("GEMINI_API_KEY is not set. Please provide it in the .env file or as an environment variable.")
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
			// Print Gemini's response
			printResponse(resp)

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
			}
			
		}
	}


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
// func streamOutput(reader *bufio.Reader, label string, ch chan string) {
// 	defer close(ch)
// 	for {
// 		line, err := reader.ReadString('\n')
// 		if len(line) > 0 {
// 			fmt.Printf("[%s] %s", label, line) // Print to user terminal
// 			ch <- line                        // Send line to output channel
// 		}
// 		if err != nil {
// 			break
// 		}
// 	}
// }



// initializeChat sets up the chat with initial file content
func initializeChat(client *genai.Client, ctx context.Context, content string) *genai.ChatSession {
	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()
	chatPrompt := fmt.Sprintf(`
	Here is the content of my files and folders:
	
	%s
	
	So you're a coding assistant or partner you're provided with the current directory content already 
	and you can also ask user to provide you the current directory content again by running the update command 
	So you're a coding assistant or partner you're provided with the current directory content already and 
	you can also ask user to provide you the current directory content again by running the update command or 
	you can also run the command on a user's terminal whenever you wanna run command on the user terminal you 
	can send the command in the format run command or you can also run the command on a user's terminal whenever you wanna run 
	command on the user terminal you can send the command in the format "run command: desired command" and then the command will be executed on the user system and you will be provided with
	 output of the command and then the command will be executed on the user system and you will be provided 
	 without output of the command so your main job is to help the programmer to get the maximum productivity and be a 
	 companion so your main job is to help the programmer to get the maximum productivity and be a companion`, content)
	
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
func scanDirectory(dir string) (string, error) {
	var content strings.Builder
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil // Skip this file or folder
		}

		if !info.IsDir() { // If it's a file
			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v\n", path, err)
				return nil // Skip this file
			}

			// Skip binary files
			if !utf8.Valid(fileContent) {
				log.Printf("Skipping binary or invalid UTF-8 file: %s\n", path)
				return nil
			}

			content.WriteString(fmt.Sprintf("File: %s\n%s\n\n", path, string(fileContent)))
		} else { // If it's a folder
			content.WriteString(fmt.Sprintf("Folder: %s\n", path))
		}
		return nil
	})
	return content.String(), err
}

//printResponse formats and prints Gemini's responses
func printResponse(resp *genai.GenerateContentResponse) {
	for _, part := range resp.Candidates[0].Content.Parts {
		// Assuming part is a string or implements a method to get the string content
		formatted := cleanMarkdown(fmt.Sprintf("%v", part)) // Convert part to string safely
		fmt.Println(formatted) // Print the cleaned response text
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
		"*", "",  // Italic
		"- ", "", // List items
		"`", "",  // Inline code
	)
	output := replacer.Replace(input)

	// Remove extra newlines
	output = strings.TrimSpace(output)
	return output
}


