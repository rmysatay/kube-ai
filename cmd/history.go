package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// 🔧 Tüm komutlar tarafından kullanılan log fonksiyonu
func SaveToHistory(command string, input string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	historyPath := filepath.Join(home, ".kube-ai-history")
	entry := fmt.Sprintf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), command, input)

	f, err := os.OpenFile(historyPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(entry)
	}
}

// 🧾 "kube-ai history" komutu
var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show command history",
	Long:  "Displays previously used kube-ai commands stored in ~/.kube-ai-history.",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("❌ Could not detect home directory:", err)
			return
		}

		historyPath := filepath.Join(home, ".kube-ai-history")
		data, err := os.ReadFile(historyPath)
		if err != nil {
			fmt.Println("ℹ️ No history found yet.")
			return
		}

		fmt.Println("📜 Command History:")
		fmt.Println(string(data))
	},
}
