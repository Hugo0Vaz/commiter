package main

import (
	"bytes"
	"context"
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	outputFlag := flag.String("output", "text", "Output format: json or text")
	langFlag := flag.String("lang", "pt-br", "Language for analysis (pt-br or en)")
	flag.Parse()

	if !isGitRepository() {
		fmt.Println("Error: Current directory is not a git repository")
		os.Exit(1)
	}

	diff, err := getStagedDiff()
	if err != nil {
		log.Fatalf("Error getting git diff: %v", err)
	}
	
	if diff == "" {
		fmt.Println("No staged changes to analyze")
		os.Exit(0)
	}
	
	analysis, err := getAIAnalysis(diff, *langFlag)
	if err != nil {
		log.Fatalf("Error getting AI analysis: %v", err)
	}

	output, err := generateOutput(analysis, *outputFlag)
	if err != nil {
		log.Fatalf("Error generating output: %v", err)
	}

	fmt.Println(output)
	os.Exit(0)
}

// Returns if the CWD is a git repository
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func getAIAnalysis(diff string, lang string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)
	var prompt string
	if lang == "en" {
		prompt = englishPrompt(diff)
	} else {
		prompt = portuguesePrompt(diff)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func getFileStagedDiff(filename string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "diff", "--cached", "--", filename)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed for file %s: %w", filename, err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateOutput(output string, output_type string) (string, error) {
	if output_type == "text" {
		return output, nil
	} else {
		return "", nil
	}
}
