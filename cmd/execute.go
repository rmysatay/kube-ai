package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var execFile string

var ExecuteCmd = &cobra.Command{
	Use:   "execute --file <yaml-file>",
	Short: "Apply a Kubernetes manifest file to the cluster",
	Long:  "Apply an existing Kubernetes YAML manifest to the cluster. The manifest content is displayed before applying.",
	Run: func(cmd *cobra.Command, args []string) {
		if execFile == "" {
			fmt.Println("❌ Please provide a file using --file or -f flag.")
			return
		}

		if _, err := os.Stat(execFile); os.IsNotExist(err) {
			fmt.Printf("❌ File '%s' does not exist.\n", execFile)
			return
		}

		// Dosyayı oku ve göster
		content, err := os.ReadFile(execFile)
		if err != nil {
			fmt.Printf("❌ Failed to read file '%s': %v\n", execFile, err)
			return
		}

		fmt.Println("\n📄 YAML Content to Apply:")
		fmt.Println("-----------------------------------")
		fmt.Println(string(content))
		fmt.Println("-----------------------------------")

		// YAML dosyasını apply et
		fmt.Println("🚀 Applying manifest to the cluster...")

		applyCmd := exec.Command("kubectl", "apply", "-f", execFile)
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr
		applyCmd.Stdin = os.Stdin

		if err := applyCmd.Run(); err != nil {
			fmt.Printf("❌ Failed to apply manifest: %v\n", err)
			return
		}

		fmt.Println("✅ Resource applied successfully!")
		fmt.Printf("📄 YAML remains at: %s\n", execFile)
	},
}

func init() {
	ExecuteCmd.Flags().StringVarP(&execFile, "file", "f", "", "Path to the YAML manifest file to apply")
}
