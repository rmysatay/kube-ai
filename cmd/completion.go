package cmd

import (
	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completion scripts",
	Long: `To load completions:

Bash:
  source <(kube-ai completion bash)

Zsh:
  source <(kube-ai completion zsh)

Fish:
  kube-ai completion fish | source

PowerShell:
  kube-ai completion powershell | Out-String | Invoke-Expression
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = RootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			_ = RootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			_ = RootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			_ = RootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
		default:
			cmd.Print("Shell not supported. Use: bash|zsh|fish|powershell")
		}
	},
}

func init() {
	RootCmd.AddCommand(CompletionCmd)
}
