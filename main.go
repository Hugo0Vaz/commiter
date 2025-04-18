package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
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

	analysis, err := getAIAnalysis(diff)
	if err != nil {
		log.Fatalf("Error getting AI analysis: %v", err)
	}

	fmt.Println(analysis)
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

func getAIAnalysis(diff string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	client := openai.NewClient(apiKey)
	prompt := fmt.Sprintf(`Analise esse diff do git e forneça uma mensagem de confirmação curta e descritiva e outra mais longa e detalhada.

A saída deve seguir este modelo: 

(ação)[arquivo ou parte do sistema]: descrição curta.

	Descrição longa e mais detalhada aqui 

Diff:
%s`, diff)

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
