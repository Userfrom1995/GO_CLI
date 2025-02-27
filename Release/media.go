package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	"github.com/google/generative-ai-go/genai"
	// "google.golang.org/api/option"
)

// extractResponse concatenates only valid UTF-8 text parts from the response.
func extractResponse(resp *genai.GenerateContentResponse) string {
	var output string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					s := string(txt)
					if utf8.ValidString(s) {
						output += s
					} else {
						log.Printf("Skipping invalid UTF-8 part")
					}
				}
			}
		}
	}
	return output
}

// ReadFileContentWithAI uploads a media file (PDF, image, video, etc.) and asks Gemini
// to analyze its content using the provided prompt. If only a file name is provided, it is
// assumed to be in the current working directory.
func ReadFileContentWithAI(ctx context.Context, client *genai.Client, filePath, prompt string) (string, error) {
	// Resolve file path: if not absolute, use current working directory.
	if !filepath.IsAbs(filePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %v", err)
		}
		filePath = filepath.Join(cwd, filePath)
	}

	// Open the file.
	osf, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer osf.Close()

	// Upload the file using UploadFile.
	file, err := client.UploadFile(ctx, "", osf, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	// Clean up the uploaded file after processing.
	defer func() {
		if err := client.DeleteFile(ctx, file.Name); err != nil {
			log.Printf("warning: failed to delete file %s: %v", file.Name, err)
		}
	}()

	// Videos need to be processed before you can use them.
	for file.State == genai.FileStateProcessing {
		log.Printf("processing %s", file.Name)
		time.Sleep(5 * time.Second)
		var err error
		if file, err = client.GetFile(ctx, file.Name); err != nil {
			log.Fatal(err)
		}
	}
	if file.State != genai.FileStateActive {
		log.Fatalf("uploaded file has state %s, not active", file.State)
	}

	// Use the generative model (e.g., "gemini-1.5-pro") to analyze the file.
	resp, err := client.GenerativeModel("gemini-1.5-pro").GenerateContent(ctx,
		genai.FileData{URI: file.URI},
		genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %v", err)
	}

	return extractResponse(resp), nil
}

var fileContentSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"filePath": {
			Type:        genai.TypeString,
			Description: "The full path or file name (if in the current directory) of the file to analyze (e.g., PDF, image, video, etc.).",
		},
		"prompt": {
			Type:        genai.TypeString,
			Description: "A prompt describing the analysis to perform (e.g., 'Describe the contents of this video in detail' or 'Summarize this document').",
		},
	},
	Required: []string{"filePath", "prompt"},
}

var FileContentTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name: "read_file_content",
			Description: "Uploads a media file and returns its analysis. Supports PDFs, images, videos, and other documents. " +
				"If the file is a video, the tool waits until the file is fully processed before generating content.",
			Parameters: fileContentSchema,
		},
	},
}
