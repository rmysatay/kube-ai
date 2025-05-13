package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"regexp"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	saveToFile bool
	outputFile string
)

var GenerateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Generate Kubernetes YAML manifest using AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		apiKey := os.Getenv("OPENAI_API_KEY")

		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		client := openai.NewClient(apiKey)

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "You are a helpful Kubernetes assistant. Your task is to output only valid and minimal YAML manifests based on user descriptions.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("❌ OpenAI error: %v\n", err)
			return
		}

		output := strings.TrimSpace(resp.Choices[0].Message.Content)

		// Clean up markdown code blocks like ```yaml ... ```
		re := regexp.MustCompile("(?s)```(?:yaml)?\\s*([\\s\\S]+?)\\s*```")
		if matches := re.FindStringSubmatch(output); len(matches) > 1 {
			output = matches[1]
		}

		output = strings.TrimSpace(output)

		fmt.Println("📄 Generated Kubernetes YAML:")
		fmt.Println("-----------------------------------")
		fmt.Println(output)
		fmt.Println("-----------------------------------")

		if saveToFile {
			file := "output.yaml"
			if outputFile != "" {
				file = outputFile
			}
			err := os.WriteFile(file, []byte(output), 0644)
			if err != nil {
				fmt.Printf("❌ Failed to save to file: %v\n", err)
				return
			}
			fmt.Printf("✅ YAML saved to file: %s\n", file)
		}
	},
}

func init() {
	GenerateCmd.Flags().BoolVarP(&saveToFile, "save", "s", false, "Save output to a file")
	GenerateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify output filename (default: output.yaml)")
}
