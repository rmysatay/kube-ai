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
	auditInputFile string
	auditResName   string
	auditNamespace string
)

var AuditCmd = &cobra.Command{
	Use:   "audit [resource-type] [resource-name] [flags] [optional: question]",
	Short: "Audit Kubernetes resources for security risks using AI",
	Long:  "Analyze Kubernetes resources (from file or live cluster) to detect security risks, misconfigurations, and policy violations using AI.",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		var auditData string
		var userQuestion string

		// Eğer inputFile verilmişse dosyadan oku
		if auditInputFile != "" {
			content, err := os.ReadFile(auditInputFile)
			if err != nil {
				fmt.Printf("❌ Failed to read file %s: %v\n", auditInputFile, err)
				return
			}
			auditData = string(content)
		} else if auditResName != "" && auditNamespace != "" {
			// Eğer resourceName ve namespace verilmişse cluster'dan oku
			kubectlCmd := exec.Command("kubectl", "get", auditResName, "-n", auditNamespace, "-o", "yaml")
			output, err := kubectlCmd.CombinedOutput()
			if err != nil {
				fmt.Printf("❌ Failed to fetch resource from cluster: %v\nOutput: %s\n", err, output)
				return
			}
			auditData = string(output)
		} else if len(args) > 0 {
			// Eğer sadece prompt verilmişse
			userQuestion = strings.Join(args, " ")
		} else {
			fmt.Println("❌ Please provide a file (-f), a resource name (--name and --ns), or a question as an argument.")
			return
		}

		// Eğer doğrudan soru yoksa otomatik audit promptu oluştur
		if userQuestion == "" {
			userQuestion = "Please audit the following Kubernetes manifest or output for security risks and best practice violations."
		}

		fullPrompt := fmt.Sprintf(`Kubernetes Resource to Audit:
---
%s
---
Task: %s
`, auditData, userQuestion)

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are a Kubernetes security auditor.
Your job is to detect any security vulnerabilities, misconfigurations, and best practice violations in Kubernetes manifests or outputs.
Focus on issues like missing resource limits, excessive permissions, absent network policies, and insecure container settings.`,
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

		fmt.Println("\n🔍 AI Audit Result:")
		fmt.Println(strings.TrimSpace(resp.Choices[0].Message.Content))
	},
}

func init() {
	AuditCmd.Flags().StringVarP(&auditInputFile, "file", "f", "", "Path to a file containing Kubernetes manifest")
	AuditCmd.Flags().StringVar(&auditResName, "name", "", "Kubernetes resource type/name (e.g., pod/mypod)")
	AuditCmd.Flags().StringVar(&auditNamespace, "ns", "", "Namespace of the resource")
}
