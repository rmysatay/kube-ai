package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION = "0.1.0" 

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the kube-ai version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kube-AI Version:", VERSION)
	},
}
