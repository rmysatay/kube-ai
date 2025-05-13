package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var DiagnoseCmd = &cobra.Command{
	Use:   "diagnose [prompt]",
	Short: "Diagnose problems in the Kubernetes cluster using AI",
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
						Content: "You are a Kubernetes troubleshooter. Help users identify and fix problems in their cluster.",
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

		fmt.Println("🛠️ Diagnosis from AI:")
		fmt.Println(resp.Choices[0].Message.Content)
	},
}
