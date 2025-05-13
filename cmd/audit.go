package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var AuditCmd = &cobra.Command{
	Use:   "audit [prompt]",
	Short: "Audit Kubernetes resources using AI guidance",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		apiKey := os.Getenv("OPENAI_API_KEY")

		if apiKey == "" {
			fmt.Println("OPENAI_API_KEY environment variable not set.")
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
						Content: "You are a Kubernetes auditor. Help users detect security risks and misconfigurations in their cluster.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("OpenAI error: %v\n", err)
			return
		}

		fmt.Println("🔍 AI Audit Result:")
		fmt.Println(resp.Choices[0].Message.Content)
	},
}
