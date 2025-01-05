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

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

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
			cs = initializeChat(client, ctx, content)
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
		}
	}
}

// initializeChat sets up the chat with initial file content
func initializeChat(client *genai.Client, ctx context.Context, content string) *genai.ChatSession {
	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	// Initialize the chat with directory contents
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(fmt.Sprintf("Here is the content of my files and folders:\n\n%s", content)),
			},
			Role: "user",
		},
	}
	return cs
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


