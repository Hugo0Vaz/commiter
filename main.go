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
	yFlag := flag.Bool("y", false, "Automatically commit without confirmation")
	flag.Parse()
	if !isGitRepository() {
		fmt.Println("Error: Current directory is not a git repository")
		os.Exit(1)
	}

	args := flag.Args()
	allSubcommand := false
	if len(args) > 0 && args[0] == "all" {
		allSubcommand = true
	}

	if allSubcommand {
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
		processAnalysis(analysis, *outputFlag, *yFlag, "")
	} else {
		// Process each staged file individually
		cmdFiles := exec.Command("git", "diff", "--cached", "--name-only")
		var out bytes.Buffer
		cmdFiles.Stdout = &out
		if err := cmdFiles.Run(); err != nil {
			log.Fatalf("Error getting staged files: %v", err)
		}
		files := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(files) == 0 || files[0] == "" {
			fmt.Println("No staged changes to analyze")
			os.Exit(0)
		}
		for _, file := range files {
			diff, err := getFileStagedDiff(file)
			if err != nil {
				log.Fatalf("Error getting staged diff for %s: %v", file, err)
			}
			if diff == "" {
				continue
			}
			analysis, err := getAIAnalysis(diff, *langFlag)
			if err != nil {
				log.Fatalf("Error getting AI analysis for %s: %v", file, err)
			}
			processAnalysis(analysis, *outputFlag, *yFlag, file)
		}
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

func processAnalysis(analysis, outputFormat string, yFlag bool, file string) {
	parts := strings.SplitN(analysis, "\n", 2)
	shortPart := parts[0]
	longPart := ""
	if len(parts) > 1 {
		longPart = parts[1]
	}
	if outputFormat == "json" {
		result := map[string]string{"short": shortPart, "long": longPart}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			log.Fatalf("Error creating JSON output: %v", err)
		}
		fmt.Println(string(jsonBytes))
	} else if outputFormat == "cmd" {
		if file != "" {
			fmt.Printf("git commit -m \"%s\" -m \"%s\" -- %s\n", shortPart, longPart, file)
		} else {
			fmt.Printf("git commit -m \"%s\" -m \"%s\"\n", shortPart, longPart)
		}
	} else {
		if !yFlag {
			fmt.Println("Commit message preview:")
			fmt.Printf("Short: %s\n", shortPart)
			fmt.Printf("Long: %s\n", longPart)
			fmt.Print("Do you want to commit? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(strings.TrimSpace(response)) != "y" {
				fmt.Println("Commit aborted.")
				return
			}
		}
		var cmd *exec.Cmd
		if file != "" {
			cmd = exec.Command("git", "commit", "-m", shortPart, "-m", longPart, "--", file)
		} else {
			cmd = exec.Command("git", "commit", "-m", shortPart, "-m", longPart)
		}
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error committing %s: %v", file, err)
		}
		if file != "" {
			fmt.Println("Committed", file)
		} else {
			fmt.Println("Committed!")
		}
	}
}
