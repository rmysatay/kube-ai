package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	inputFile string
	resName   string
	namespace string
)

var AnalyzeCmd = &cobra.Command{
	Use:   "analyze [resource-type] [resource-name] --ns <namespace> [question]",
	Short: "Analyze Kubernetes resources or errors using AI",
	Long:  "Analyze raw kubectl outputs (describe, logs, events, etc.) or YAML manifests and diagnose Kubernetes issues with the help of AI.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		// Eğer --name boşsa ve en az 2 argüman varsa (resource-type ve resource-name), otomatik doldur
		if resName == "" && len(args) >= 2 {
			resourceType := args[0]
			resourceName := args[1]
			resName = fmt.Sprintf("%s/%s", resourceType, resourceName)
			args = args[2:]
		}

		// Eğer hem file hem name yoksa hata ver
		if resName == "" && inputFile == "" {
			fmt.Println("❌ Please provide either a file (-f) or a resource name (--name) with namespace (--ns).")
			return
		}

		// Eğer kaynak adı kullanılıyorsa namespace zorunlu olsun
		if resName != "" && namespace == "" {
			fmt.Println("❌ Namespace (--ns) is required when using resource name (--name) or positional args.")
			return
		}

		question := strings.Join(args, " ")
		if question == "" {
			fmt.Println("❌ Please provide a question for AI analysis.")
			return
		}

		SaveToHistory("analyze", fmt.Sprintf("resName=%s ns=%s file=%s question=%s", resName, namespace, inputFile, question))

		kubeData, err := getKubernetesData()
		if err != nil {
			fmt.Printf("❌ Error collecting Kubernetes data: %v\n", err)
			return
		}

		fullPrompt := fmt.Sprintf(`Analyze the following Kubernetes output and answer the user's question.

--- Start of Kubernetes Output ---
%s
--- End of Kubernetes Output ---

User question: %s`, kubeData, question)

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are a certified Kubernetes expert.
Your task is to analyze raw kubectl command outputs (describe, logs, events, etc.) or YAML manifests and help diagnose issues.
If a pod is crashing, stuck, or unhealthy, identify root causes such as image pull errors, readiness probe failures, or insufficient resources.
Provide detailed reasoning, potential root causes, and suggested fixes.`,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: fullPrompt,
					},
				},
				MaxTokens: MaxTokens,
			},
		)

		if err != nil {
			fmt.Printf("❌ OpenAI error: %v\n", err)
			return
		}

		fmt.Println("\n🤖 AI Analysis:")
		fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
	},
}

// Kubernetes çıktısını dosyadan ya da cluster'dan al
func getKubernetesData() (string, error) {
	if inputFile != "" {
		content, err := os.ReadFile(inputFile)
		if err != nil {
			return "", fmt.Errorf("failed to read file %s: %w", inputFile, err)
		}
		return string(content), nil
	} else if resName != "" && namespace != "" {
		cmd := exec.Command("kubectl", "get", resName, "-n", namespace, "-o", "yaml")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("failed to get resource: %w\nOutput: %s", err, output)
		}
		return string(output), nil
	}
	return "", fmt.Errorf("please provide either --file or both --name and --ns parameters")
}

func init() {
	AnalyzeCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Path to a file containing kubectl output")
	AnalyzeCmd.Flags().StringVar(&resName, "name", "", "Name and type of Kubernetes resource (e.g. deployment/nginx)")
	AnalyzeCmd.Flags().StringVar(&namespace, "ns", "", "Namespace of the resource (required)")
}
