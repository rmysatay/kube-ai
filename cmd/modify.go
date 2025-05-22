package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var (
	modifyFile   string
	newNamespace string
	newReplicas  int
	newName      string
)

var ModifyCmd = &cobra.Command{
	Use:   "modify --file <yaml> [options]",
	Short: "Modify a Kubernetes YAML manifest",
	Long:  "Modify fields like namespace, replicas, and metadata name inside an existing YAML manifest file. It does not apply the changes automatically; you can apply them later using execute.",
	Run: func(cmd *cobra.Command, args []string) {
		if modifyFile == "" {
			fmt.Println("❌ Please specify a YAML file with --file")
			return
		}

		SaveToHistory("modify", fmt.Sprintf("file=%s ns=%s name=%s replicas=%d", modifyFile, newNamespace, newName, newReplicas))

		data, err := os.ReadFile(modifyFile)
		if err != nil {
			fmt.Println("❌ Failed to read YAML file:", err)
			return
		}

		var manifest map[string]interface{}
		err = yaml.Unmarshal(data, &manifest)
		if err != nil {
			fmt.Println("❌ Failed to parse YAML:", err)
			return
		}

		if newNamespace != "" {
			if meta, ok := manifest["metadata"].(map[string]interface{}); ok {
				meta["namespace"] = newNamespace
			}
		}

		if newName != "" {
			if meta, ok := manifest["metadata"].(map[string]interface{}); ok {
				meta["name"] = newName
			}
		}

		if newReplicas > 0 {
			if spec, ok := manifest["spec"].(map[string]interface{}); ok {
				spec["replicas"] = newReplicas
			}
		}

		newYaml, err := yaml.Marshal(manifest)
		if err != nil {
			fmt.Println("❌ Failed to marshal YAML:", err)
			return
		}

		err = os.WriteFile(modifyFile, newYaml, 0644)
		if err != nil {
			fmt.Println("❌ Failed to save updated YAML:", err)
			return
		}

		fmt.Println("✅ YAML updated successfully:", modifyFile)
		fmt.Println("👉 If you want to apply it, run: kube-ai execute --file", modifyFile)
	},
}

func init() {
	ModifyCmd.Flags().StringVarP(&modifyFile, "file", "f", "", "YAML file to modify")
	ModifyCmd.Flags().StringVar(&newNamespace, "namespace", "", "New namespace to set")
	ModifyCmd.Flags().IntVar(&newReplicas, "replicas", 0, "New replica count to set")
	ModifyCmd.Flags().StringVar(&newName, "name", "", "New metadata name to set")
}
