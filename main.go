package main

import (
	"bufio"
	"context"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const GenaiModel = "gemini-1.5-flash" // model to use

type App struct {
	client *genai.Client
	model  *genai.GenerativeModel
	cs     *genai.ChatSession
}

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	reader := bufio.NewReader(os.Stdin)

	genaiApp := &App{}

	apiKey := os.Getenv("GEMINI_API_KEY")
	genaiApp.client, err = NewClient(apiKey, context.Background())
	if err != nil {
		log.Fatalf("Error creating client")
	}

	genaiApp.model = NewModel(genaiApp.client, GenaiModel)
	genaiApp.model.Tools = []*genai.Tool{FileTool, ScanTool, ReadFileTool, RunCommandTool, SystemInfoTool}
	genaiApp.cs = genaiApp.model.StartChat()

	for {
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		response, err := genaiApp.cs.SendMessage(context.Background(), genai.Text(input))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}

		responseString := buildResponse(response, genaiApp.cs)

		log.Println("Response:", responseString)

	}
}

// buildResponse builds a string response based on content parts from candidates
func buildResponse(resp *genai.GenerateContentResponse, cs *genai.ChatSession) string {
	funcResponse := make(map[string]interface{})
	var err error

	for _, part := range resp.Candidates[0].Content.Parts {
		functionCall, ok := part.(genai.FunctionCall)
		if ok {
			log.Println("Function call:", functionCall.Name)
			switch functionCall.Name {
			case "file_write":
				fileName, fileNameOk := functionCall.Args["fileName"].(string)
				content, contentOk := functionCall.Args["content"].(string)

				if !fileNameOk || fileName == "" {
					funcResponse["error"] = "expected non-empty string at key 'fileName'"
					break
				}
				if !contentOk || content == "" {
					funcResponse["error"] = "expected non-empty string at key 'content'"
					break
				}
				err := WriteDesktop(fileName, content)
				if err != nil {
					funcResponse["error"] = "could not write file."
				} else {
					funcResponse["result"] = "file successfully written"
				}

			case "scan_directory":
				directory, ok := functionCall.Args["directory"].(string)
				if !ok {
					funcResponse = map[string]interface{}{"error": "Invalid directory path"}
					break
				}
				result, err := scanDirectory(directory)
				if err != nil {
					funcResponse["error"] = err.Error()
				} else {
					funcResponse["result"] = result
				}

			// default:
			// 	response = map[string]interface{}{"error": "Unknown function call"}
			//
			case "ReadFile":
				fileName, ok := functionCall.Args["fileName"].(string)
				if !ok || fileName == "" {
					fileName, ok = functionCall.Args["directory"].(string) // Fallback if Gemini is sending "directory"
					if !ok || fileName == "" {
						funcResponse["error"] = "expected non-empty string at key 'fileName' or 'directory'"
						break
					}
				}

				content, err := ReadFile(fileName)
				if err != nil {
					funcResponse["error"] = err.Error()
				} else {
					funcResponse["result"] = content
				}

			//case "get_system_info":
			//	sysInfo, err := GetSystemInfo()
			//	if err != nil {
			//		funcResponse["error"] = "failed to retrieve system information"
			//	} else {
			//		funcResponse["result"] = sysInfo
			//	}
			case "run_command":
				command, ok := functionCall.Args["command"].(string)
				if !ok || command == "" {
					funcResponse["error"] = "expected a non-empty string for 'command'"
					break
				}
				args, _ := functionCall.Args["args"].(string) // optional, so no need to check
				output, err := RunCommand(command, args)
				if err != nil {
					funcResponse["error"] = err.Error()
				} else {
					funcResponse["result"] = output
				}

			case "get_system_info":
				sysInfo, err := GetSystemSpecs()
				if err != nil {
					funcResponse["error"] = "failed to retrieve system information: " + err.Error()
				} else {
					// Convert map[string]string to map[string]interface{}
					convertedSysInfo := make(map[string]interface{})
					for k, v := range sysInfo {
						convertedSysInfo[k] = v
					}
					funcResponse["result"] = convertedSysInfo
				}

			default:
				funcResponse["error"] = "unknown function call"
			}
		}
	}

	if len(funcResponse) > 0 {
		resp, err = cs.SendMessage(context.Background(), genai.FunctionResponse{
			Name:     "Function_Call",
			Response: funcResponse,
		})
		if err != nil {
			return "Error sending message: " + err.Error()
		}
		funcResponse = nil
		return buildResponse(resp, cs)
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				res, ok := part.(genai.Text)
				if ok {
					return string(res)
				}

			}

		}
	}
	return ""
}
