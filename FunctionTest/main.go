package main

import (
	"fmt"
	"log"
)

// RunCommand takes a full command line as a single string,
// splits it into the command and its arguments, executes it,
// and returns the combined output.
//func RunCommand(cmdLine string) (string, error) {
//	parts := strings.Fields(cmdLine)
//	if len(parts) == 0 {
//		return "", fmt.Errorf("no command provided")
//	}
//	cmd := exec.Command(parts[0], parts[1:]...)
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		return "", fmt.Errorf("failed to execute command: %v\nOutput: %s", err, output)
//	}
//	return string(output), nil
//}

func main() {
	//ctx := context.Background()

	// Check Git username
	username, err := RunCommand("cd ..")
	if err != nil {
		log.Fatalf("Error getting Git username: %v", err)
	}

	// Check Git email
	//email, err := RunCommand("git config --get user.email")
	//if err != nil {
	//	log.Fatalf("Error getting Git email: %v", err)
	//}
	fmt.Println(username)
	//

}
