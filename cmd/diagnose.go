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
	diagnoseInputFile string
	diagnoseResName   string
	diagnoseNamespace string
)

var DiagnoseCmd = &cobra.Command{
	Use:   "diagnose [resource-type] [resource-name] [flags] [optional: question]",
	Short: "Diagnose problems in Kubernetes pods using AI",
	Long:  "Troubleshoot issues in Kubernetes pods by analyzing describe outputs, logs, or manifest configurations with the help of AI.",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		var diagnoseData string
		var userQuestion string

		// Eğer inputFile verilmişse dosyadan oku
		if diagnoseInputFile != "" {
			content, err := os.ReadFile(diagnoseInputFile)
			if err != nil {
				fmt.Printf("❌ Failed to read file %s: %v\n", diagnoseInputFile, err)
				return
			}
			diagnoseData = string(content)
		} else if diagnoseResName != "" && diagnoseNamespace != "" {
			// Eğer resourceName ve namespace verilmişse cluster'dan describe çek
			kubectlCmd := exec.Command("kubectl", "describe", diagnoseResName, "-n", diagnoseNamespace)
			output, err := kubectlCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Failed to describe resource from cluster: %v\nOutput: %s\n", err, output)
				return
			}
			diagnoseData = string(output)
		} else if len(args) > 0 {
			// Sadece prompt verilmişse
			userQuestion = strings.Join(args, " ")
		} else {
			fmt.Println("❌ Please provide a file (-f), a resource name (--name and --ns), or a direct question as an argument.")
			return
		}

		// Eğer kullanıcı özel bir soru yazmadıysa, otomatik default soru yaz
		if userQuestion == "" {
			userQuestion = "Please diagnose the issue in the following Kubernetes pod output."
		}

		fullPrompt := fmt.Sprintf(`Pod Output:
---
%s
---
Task: %s
`, diagnoseData, userQuestion)

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are a Kubernetes troubleshooter.
Analyze pod outputs such as describe results, logs, and events.
Identify problems like CrashLoopBackOff, OOMKilled, ImagePullBackOff, readiness probe failures, node pressure, etc.
Provide a clear diagnosis and suggest potential fixes.`,
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

		fmt.Println("\n🛠️ Diagnosis from AI:")
		fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
	},
}

func init() {
	DiagnoseCmd.Flags().StringVarP(&diagnoseInputFile, "file", "f", "", "Path to a file containing pod describe output or logs")
	DiagnoseCmd.Flags().StringVar(&diagnoseResName, "name", "", "Kubernetes resource type/name (e.g., pod/mypod)")
	DiagnoseCmd.Flags().StringVar(&diagnoseNamespace, "ns", "", "Namespace of the resource")
}
