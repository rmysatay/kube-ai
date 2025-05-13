package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var suggestFile string

var SuggestCmd = &cobra.Command{
	Use:   "suggest [question]",
	Short: "Ask AI to suggest appropriate kubectl commands for troubleshooting, applying, or inspecting Kubernetes resources",
	Long: `You can provide either:
  - A question like "How do I debug a pod in CrashLoopBackOff?"
  - Or a YAML file using --file/-f for kubectl command suggestions.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		client := openai.NewClient(apiKey)
		var userPrompt string

		if suggestFile != "" {
			content, err := os.ReadFile(suggestFile)
			if err != nil {
				fmt.Printf("❌ Failed to read file: %v\n", err)
				return
			}
			userPrompt = fmt.Sprintf(`This is a Kubernetes manifest:

---
%s
---

What would be the correct kubectl command to apply or interact with this manifest? Respond with only the shell command.`, string(content))
		} else if len(args) == 1 {
			userPrompt = fmt.Sprintf(`I want to troubleshoot or manage my Kubernetes resources.

Question:
%s

Please suggest valid kubectl command(s). Respond only with shell commands.`, args[0])
		} else {
			fmt.Println("❌ Please provide a question or a file.")
			return
		}

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "You are a Kubernetes expert assistant. Only return valid shell commands using kubectl. Avoid explanations unless requested.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: userPrompt,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("❌ OpenAI error: %v\n", err)
			return
		}

		fmt.Println("🤖 Suggested kubectl command(s):")
		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func init() {
	SuggestCmd.Flags().StringVarP(&suggestFile, "file", "f", "", "Path to a Kubernetes manifest file")
}
