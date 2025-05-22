package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Model         string
	MaxTokens     int
	CountTokens   bool
	Verbose       bool
	MaxIterations int

	RootCmd = &cobra.Command{
		Use:     "kube-ai",
		Version: "0.1.0",
		Short:   "AI-powered Kubernetes CLI assistant",
		Long:    "Kube-AI provides AI-powered analysis, auditing, YAML generation, and troubleshooting for Kubernetes resources.",
	}
)

func init() {
	// Global flags for all subcommands
	RootCmd.PersistentFlags().StringVarP(&Model, "model", "m", "gpt-4o", "AI model to use (default: gpt-4o)")
	RootCmd.PersistentFlags().IntVarP(&MaxTokens, "max-tokens", "t", 2048, "Maximum tokens for AI responses (default: 2048)")
	RootCmd.PersistentFlags().BoolVarP(&CountTokens, "count-tokens", "c", false, "Print token usage after request")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().IntVarP(&MaxIterations, "max-iterations", "x", 30, "Maximum iterations for multi-step analysis (reserved)")

	// Register subcommands
	RootCmd.AddCommand(
		AnalyzeCmd,
		AuditCmd,
		DiagnoseCmd,
		GenerateCmd,
		ExecuteCmd,
		ChatCmd,
		VersionCmd,
		ModifyCmd,
		CompletionCmd,
		HistoryCmd, // 🟡 Bu komutun history.go içinde tanımlı olduğundan emin olun
	)
}