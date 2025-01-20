// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"os/exec"
// )

// func main() {
// 	for {
// 		fmt.Print("Enter command to run (or type 'exit' to quit): ")
// 		var command string
// 		reader := bufio.NewReader(os.Stdin)
// 		command, _ = reader.ReadString('\n')
// 		command = command[:len(command)-1] // Remove the newline character

// 		if command == "exit" {
// 			fmt.Println("Exiting program.")
// 			break
// 		}

// 		err := runCommand(command)
// 		if err != nil {
// 			fmt.Printf("Error running command: %v\n", err)
// 		}
// 	}
// }

// func runCommand(command string) error {
// 	// Split the command and arguments
// 	cmdParts := parseCommand(command)
// 	if len(cmdParts) == 0 {
// 		return fmt.Errorf("invalid command")
// 	}

// 	// Create the command
// 	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

// 	// Get pipes for stdout and stderr
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return fmt.Errorf("error creating stdout pipe: %v", err)
// 	}

// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		return fmt.Errorf("error creating stderr pipe: %v", err)
// 	}

// 	// Start the command
// 	if err := cmd.Start(); err != nil {
// 		return fmt.Errorf("error starting command: %v", err)
// 	}

// 	// Create readers for the pipes
// 	stdoutReader := bufio.NewReader(stdout)
// 	stderrReader := bufio.NewReader(stderr)

// 	// Stream the output
// 	go streamOutput(stdoutReader, "STDOUT")
// 	go streamOutput(stderrReader, "STDERR")

// 	// Wait for the command to complete
// 	if err := cmd.Wait(); err != nil {
// 		return fmt.Errorf("command execution failed: %v", err)
// 	}

// 	fmt.Println("Command completed successfully.")
// 	return nil
// }

// // streamOutput reads and prints output line by line
// func streamOutput(reader *bufio.Reader, label string) {
// 	for {
// 		line, err := reader.ReadString('\n')
// 		if len(line) > 0 {
// 			fmt.Printf("[%s] %s", label, line)
// 		}
// 		if err != nil {
// 			break
// 		}
// 	}
// }

// // parseCommand splits the input into command and arguments
// func parseCommand(input string) []string {
// 	return exec.Command("sh", "-c", input).Args
// }
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for the API key
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	GEMINI_API_KEY := os.Getenv("GEMINI_API_KEY")

	// Initialize the Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(GEMINI_API_KEY))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer client.Close()

	// Define the function request
	request := genai.Text{
		Prompt: "Calculate the sum of two numbers.",
		FunctionCall: genai.FunctionCall{
			Name: "calculate_sum",
			Arguments: map[string]interface{}{
				"num1": 5,
				"num2": 7,
			},
		},
	}

	// Send the request
	iter := client.GenerateText(ctx, request)
	for {
		response, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error in generating response: %v", err)
		}

		// Process the response
		fmt.Printf("Gemini Response: %s\n", response.Text)

		// Check if the response includes function output
		if response.FunctionCall != nil {
			fmt.Printf("Function Name: %s\n", response.FunctionCall.Name)
			fmt.Printf("Function Output: %v\n", response.FunctionCall.Output)
		}
	}
}
