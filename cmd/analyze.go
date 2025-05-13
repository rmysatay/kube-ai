package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	inputFile string
	resName   string
	namespace string
)

var AnalyzeCmd = &cobra.Command{
	Use:   "analyze [question]",
	Short: "Analyze Kubernetes resources or errors using AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		question := args[0]
		fileContent := ""

		// Öncelik: Dosya okunursa onu al
		if inputFile != "" {
			content, err := os.ReadFile(inputFile)
			if err != nil {
				fmt.Printf("❌ Failed to read file %s: %v\n", inputFile, err)
				return
			}
			fileContent = string(content)
		} else if resName != "" && namespace != "" {
			// Kaynak adı ve namespace verilmişse, kubectl çıktısını al
			cmd := exec.Command("kubectl", "get", resName, "-n", namespace, "-o", "yaml")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Failed to get resource from cluster: %v\nOutput: %s\n", err, output)
				return
			}
			fileContent = string(output)
		} else {
			fmt.Println("❌ Please provide either --file or both --name and --ns parameters.")
			return
		}

		// AI'ye gönderilecek prompt
		fullPrompt := fmt.Sprintf(`Analyze the following Kubernetes data and answer the user's question.

--- Start of Kubernetes Output ---
%s
--- End of Kubernetes Output ---

User question: %s`, fileContent, question)

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are a certified Kubernetes expert.
Your task is to analyze raw kubectl command outputs (describe, logs, events, etc.) or YAML manifests and help diagnose issues.
If a pod is crashing, stuck, or unhealthy, identify root causes such as image pull errors, readiness probe failures, or insufficient resources.
Give detailed reasoning, possible root causes, and potential fixes.`,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fullPrompt,
					},
				},
			},
		)

		if err != nil {
			fmt.Printf("❌ OpenAI error: %v\n", err)
			return
		}

		fmt.Println("🤖 AI Analysis:")
		fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
	},
}

func init() {
	AnalyzeCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to a file containing kubectl output")
	AnalyzeCmd.Flags().StringVar(&resName, "name", "", "Name and type of Kubernetes resource (e.g. deployment/nginx)")
	AnalyzeCmd.Flags().StringVar(&namespace, "ns", "", "Namespace of the resource")
}
