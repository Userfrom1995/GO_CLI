// package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/google/generative-ai-go/genai"
// 	"google.golang.org/api/iterator"
// 	"google.golang.org/api/option"
// 	"github.com/joho/godotenv"
// )
// package main

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/google/generative-ai-go/genai"
// 	"github.com/joho/godotenv"
// 	"google.golang.org/api/option"
// )

// func main() {
// 	// Load .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Warning: No .env file found, falling back to system environment variables")
// 	}

// 	// Get the API key
// 	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
// 	if GEMINI_API_KEY == "" {
// 		log.Fatal("GEMINI_API_KEY is not set. Please provide it in the .env file or as an environment variable.")
// 	}

// 	// Create a new GenAI client
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
// 	if err != nil {
// 		log.Fatalf("Failed to create GenAI client: %v", err)
// 	}
// 	defer client.Close()

// 	// Specify the directory to scan
// 	dir := "./your-directory-path"

// 	// Scan the directory for files and folders
// 	content, err := scanDirectory(dir)
// 	if err != nil {
// 		log.Fatalf("Error scanning directory: %v", err)
// 	}

// 	// Send context to Gemini
// 	prompt := fmt.Sprintf("Here is the content of my files and folders:\n\n%s", content)
// 	fmt.Println("Sending context to Gemini...")
// 	chat, err := client.CreateChat(ctx, &genai.CreateChatRequest{
// 		Prompt: prompt,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to create chat: %v", err)
// 	}

// 	// Start interactive conversation
// 	fmt.Println("Context sent! You can now chat with Gemini.")
// 	for {
// 		fmt.Print("> ")
// 		var userInput string
// 		fmt.Scanln(&userInput)

// 		// Send user input to Gemini and get response
// 		resp, err := client.SendMessage(ctx, &genai.SendMessageRequest{
// 			ChatID: chat.ChatID,
// 			Prompt: userInput,
// 		})
// 		if err != nil {
// 			log.Printf("Error sending message: %v\n", err)
// 			continue
// 		}

// 		// Print Gemini's response
// 		fmt.Println(resp.Content)
// 	}
// }

// // scanDirectory recursively scans a directory and reads file contents
// func scanDirectory(dir string) (string, error) {
// 	var content string
// 	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() { // If it's a file
// 			fileContent, err := ioutil.ReadFile(path)
// 			if err != nil {
// 				return err
// 			}
// 			content += fmt.Sprintf("File: %s\n%s\n\n", path, string(fileContent))
// 		} else { // If it's a folder
// 			content += fmt.Sprintf("Folder: %s\n", path)
// 		}
// 		return nil
// 	})
// 	return content, err
// }

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}
// 	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()
//     dir := "."
// 	content, err := scanDirectory(dir)
// 	if err != nil {
// 		log.Fatalf("Error scanning directory: %v", err)
// 	}
//     prompt := fmt.Sprintf("Here is the content of my files and folders:\n\n%s", content)
// 	fmt.Println("Sending context to Gemini...")
// 	model := client.GenerativeModel("gemini-1.5-flash")
// 	cs := model.StartChat()

// 	cs.History = []*genai.Content{
// 		{
// 			Parts: []genai.Part{
// 				genai.Text("look at this directory"),
// 			},
// 			Role: "user",
// 		},
// 		{
// 			Parts: []genai.Part{
// 				genai.Blob{prompt},
// 			},
// 			Role: "model",
// 		},
// 	}
// 	iter := cs.SendMessageStream(ctx, genai.Text("look at this directory"))
// for {
// 	resp, err := iter.Next()
// 	if err == iterator.Done {
// 		break
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	printResponse(resp)
// }

// 	// iter := client.ListFiles(ctx)
// 	// for {
// 	// 	ifile, err := iter.Next()
// 	// 	if err == iterator.Done {
// 	// 		break
// 	// 	}
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	fmt.Println(ifile.Name)
// 	// }

// }
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	// Prepare the context prompt
	prompt := fmt.Sprintf("Here is the content of my files and folders:\n\n%s", content)

	// Send context to Gemini
	fmt.Println("Sending context to Gemini...")
	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	// Initialize the chat with file contents
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(prompt),
			},
			Role: "user",
		},
	}

	// Start interactive conversation
	fmt.Println("Context sent! You can now chat with Gemini.")
	for {
		fmt.Print("> ") // Prompt for user input
		var userInput string
		fmt.Scanln(&userInput)

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
			printResponse(resp)
		}
	}
}

// scanDirectory recursively scans a directory and reads file contents
func scanDirectory(dir string) (string, error) {
	var content string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() { // If it's a file
			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content += fmt.Sprintf("File: %s\n%s\n\n", path, string(fileContent))
		} else { // If it's a folder
			content += fmt.Sprintf("Folder: %s\n", path)
		}
		return nil
	})
	return content, err
}

// printResponse formats and prints Gemini's responses
func printResponse(resp *genai.Text) {
	for _, part := range resp.Content.Parts {
		switch v := part.(type) {
		case genai.Text:
			fmt.Println(v.Text)
		default:
			fmt.Printf("Unhandled response part: %+v\n", part)
		}
	}
}
