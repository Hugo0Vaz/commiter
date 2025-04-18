package main

import (
	"bytes"
	"context"
	"encoding/json"
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

	if *outputFlag == "json" {
		parts := strings.SplitN(analysis, "\n\n", 2)
		shortPart := parts[0]
		longPart := ""
		if len(parts) > 1 {
			longPart = parts[1]
		}
		result := map[string]string{"short": shortPart, "long": longPart}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("Error creating JSON output: %v", err)
		}
		fmt.Println(string(jsonBytes))
	} else if *outputFlag == "cmd" {
		parts := strings.SplitN(analysis, "\n\n", 2)
		shortPart := parts[0]
		longPart := ""
		if len(parts) > 1 {
			longPart = parts[1]
		}
		fmt.Printf("git commit -m \"%s\" -m \"%s\"\n", shortPart, longPart)
	} else {
		fmt.Println(analysis)
	}
}

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
			Model: openai.GPT3Dot5Turbo,
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
