package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	for {
		fmt.Print("Enter command to run (or type 'exit' to quit): ")
		var command string
		reader := bufio.NewReader(os.Stdin)
		command, _ = reader.ReadString('\n')
		command = command[:len(command)-1] // Remove the newline character

		if command == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		err := runCommand(command)
		if err != nil {
			fmt.Printf("Error running command: %v\n", err)
		}
	}
}

func runCommand(command string) error {
	// Split the command and arguments
	cmdParts := parseCommand(command)
	if len(cmdParts) == 0 {
		return fmt.Errorf("invalid command")
	}

	// Create the command
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// Create readers for the pipes
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	// Stream the output
	go streamOutput(stdoutReader, "STDOUT")
	go streamOutput(stderrReader, "STDERR")

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	fmt.Println("Command completed successfully.")
	return nil
}

// streamOutput reads and prints output line by line
func streamOutput(reader *bufio.Reader, label string) {
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			fmt.Printf("[%s] %s", label, line)
		}
		if err != nil {
			break
		}
	}
}

// parseCommand splits the input into command and arguments
func parseCommand(input string) []string {
	return exec.Command("sh", "-c", input).Args
}
