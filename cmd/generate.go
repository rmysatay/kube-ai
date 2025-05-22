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

var (
	saveToFile      bool
	outputFile      string
	customNamespace string
	customReplicas  int
	customName      string
)

var GenerateCmd = &cobra.Command{
	Use:   "generate [resource description]",
	Short: "Generate Kubernetes YAML manifest using AI",
	Long:  "Use AI to generate Kubernetes YAML manifests (Deployments, StatefulSets, DaemonSets, Services, etc.) based on user description. You can specify additional parameters like namespace, replicas, and metadata name.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("❌ OPENAI_API_KEY environment variable not set.")
			return
		}

		basePrompt := strings.Join(args, " ")
		extraPrompt := ""
		if customNamespace != "" {
			extraPrompt += fmt.Sprintf(" Use namespace '%s'.", customNamespace)
		}
		if customReplicas > 0 {
			extraPrompt += fmt.Sprintf(" Set replicas to %d.", customReplicas)
		}
		if customName != "" {
			extraPrompt += fmt.Sprintf(" Set metadata name to '%s'.", customName)
		}

		finalPrompt := basePrompt + extraPrompt +
			" ONLY return raw Kubernetes YAML. Do not include any titles, explanations, or code blocks."

		SaveToHistory("generate", fmt.Sprintf("desc='%s' ns=%s replicas=%d name=%s save=%v output=%s",
			basePrompt, customNamespace, customReplicas, customName, saveToFile, outputFile))

		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: Model,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are a Kubernetes YAML generator.
Return only raw YAML manifests without any markdown, code blocks, or titles.
Do not include any text like 'Deployment manifest', 'Service manifest', or 'yaml'. Only valid YAML content.`,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: finalPrompt,
					},
				},
				MaxTokens: MaxTokens,
			},
		)

		if err != nil {
			fmt.Println("❌ OpenAI error:", err)
			return
		}

		output := strings.TrimSpace(resp.Choices[0].Message.Content)

// Clean unwanted markdown, headers, and formatting
re := regexp.MustCompile("(?i)^(---|Deployment\\.yaml:|Service\\.yaml:|Deployment YAML:|Service YAML:|yaml|```yaml|```)\\s*$")

lines := strings.Split(output, "\n")
cleaned := []string{}

for _, line := range lines {
	trimmed := strings.TrimSpace(line)

	// Remove unwanted markers
	if re.MatchString(trimmed) {
		continue
	}

	// Strip leading/trailing quotes and backticks
	trimmed = strings.Trim(trimmed, "`\"")
	if trimmed != "" {
		cleaned = append(cleaned, trimmed)
	}
}

// Join all cleaned lines with a single newline
output = strings.Join(cleaned, "\n")
output = strings.TrimSpace(output)



		// Print output
		fmt.Println("\n📄 Generated Kubernetes YAML:")
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
				fmt.Println("❌ Failed to save YAML to file:", err)
				return
			}
			fmt.Println("✅ YAML saved to file:", file)
		}
	},
}

func init() {
	GenerateCmd.Flags().BoolVarP(&saveToFile, "save", "s", false, "Save the generated YAML output to a file")
	GenerateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify output filename (default: output.yaml)")
	GenerateCmd.Flags().StringVar(&customNamespace, "namespace", "", "Specify a custom namespace")
	GenerateCmd.Flags().IntVar(&customReplicas, "replicas", 0, "Specify number of replicas")
	GenerateCmd.Flags().StringVar(&customName, "name", "", "Specify a custom metadata name")
}