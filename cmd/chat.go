package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var chatInputFile string

var ChatCmd = &cobra.Command{
	Use:   "chat [your question]",
	Short: "Chat with AI for Kubernetes help and CLI command suggestions",
	Long: `Ask AI how to perform Kubernetes tasks.
If you mention or provide a YAML file, it will read that file and give more context-aware command suggestions.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		question := strings.Join(args, " ")
		SaveToHistory("chat", fmt.Sprintf("file=%s question=%s", chatInputFile, question))

		yamlContent := ""

		if chatInputFile != "" {
			content, err := os.ReadFile(chatInputFile)
			if err != nil {
				fmt.Printf("❌ Failed to read file %s: %v\n", chatInputFile, err)
				return
			}
			yamlContent = string(content)
		} else {
			matches := regexp.MustCompile(`[\w\-_]+\.ya?ml`).FindStringSubmatch(question)
			if len(matches) > 0 {
				detectedFile := matches[0]
				if _, err := os.Stat(detectedFile); err == nil {
					content, err := os.ReadFile(detectedFile)
					if err == nil {
						yamlContent = string(content)
					}
				}
			}
		}

		systemPrompt := `You are a Kubernetes expert and CLI assistant.
Always answer user questions with short and clear Kubernetes CLI examples, YAML snippets, or precise step-by-step instructions.
Focus on practical guidance only. Example:
- "kubectl apply -f filename.yaml"
- "kubectl get pods --namespace=my-namespace"`

		if yamlContent != "" {
			systemPrompt += fmt.Sprintf(`

The user also has the following YAML file open:

--- YAML START ---
%s
--- YAML END ---`, yamlContent)
		}

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: question,
					},
				},
				MaxTokens: MaxTokens,
			},
		)

		if err != nil {
			fmt.Printf("❌ OpenAI error: %v\n", err)
			return
		}

		fmt.Println("\n🤖 AI Kubernetes Assistant:")
		fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
	},
}

func init() {
	ChatCmd.Flags().StringVarP(&chatInputFile, "file", "f", "", "Optional path to a YAML file for context-aware answers")
}
