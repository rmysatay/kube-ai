package cmd

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kube-ai",
	Short: "AI-powered Kubernetes helper",
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
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}

		fmt.Println("Response from ChatGPT:")
		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
