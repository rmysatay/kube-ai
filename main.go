package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"kube-ai/cmd"
)

const VERSION = "0.1.0"

var (
	model         string
	maxTokens     int
	countTokens   bool
	verbose       bool
	maxIterations int

	rootCmd = &cobra.Command{
		Use:     "kube-ai",
		Version: VERSION,
		Short:   "Kubernetes Assistant powered by AI",
	}
)

func init() {
	
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  .env file not loaded.")
	}

	rootCmd.PersistentFlags().StringVarP(&model, "model", "m", "gpt-4o", "AI model to use")
	rootCmd.PersistentFlags().IntVarP(&maxTokens, "max-tokens", "t", 2048, "Max tokens for the AI model")
	rootCmd.PersistentFlags().BoolVarP(&countTokens, "count-tokens", "c", false, "Print token usage")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().IntVarP(&maxIterations, "max-iterations", "x", 30, "Maximum iterations for agent")

	rootCmd.AddCommand(cmd.AnalyzeCmd)
	rootCmd.AddCommand(cmd.AuditCmd)
	rootCmd.AddCommand(cmd.DiagnoseCmd)
	rootCmd.AddCommand(cmd.GenerateCmd)
	rootCmd.AddCommand(cmd.ExecuteCmd)
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.SuggestCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
