package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var execFile string

var ExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: "Apply a Kubernetes manifest file to the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if execFile == "" {
			fmt.Println("❌ Please provide a file using --file or -f")
			return
		}

		if _, err := os.Stat(execFile); os.IsNotExist(err) {
			fmt.Printf("❌ File '%s' does not exist.\n", execFile)
			return
		}

		fmt.Println("🚀 Applying manifest to cluster...")
		applyCmd := exec.Command("kubectl", "apply", "-f", execFile)
		output, err := applyCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Failed to apply manifest:\n%s\n", output)
			return
		}

		fmt.Println("✅ Resource applied successfully:")
		fmt.Println(string(output))
	},
}

func init() {
	ExecuteCmd.Flags().StringVarP(&execFile, "file", "f", "", "Path to the YAML manifest file to apply")
}
