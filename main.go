package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rmysatay/kube-ai/cmd"
)

func main() {
	// .env dosyası varsa yükle
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ Warning: .env file not loaded.")
	}

	// CLI komutlarını çalıştır
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println("❌ Error executing command:", err)
		os.Exit(1)
	}
}