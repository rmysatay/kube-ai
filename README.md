
```markdown
# kube-ai

A simple AI-powered CLI tool to assist with Kubernetes commands and operations using OpenAI's GPT API.

## ✨ Features

- 🌐 Uses OpenAI GPT models to process and generate Kubernetes-related information
- ⚙️ CLI structure built with Cobra
- 🔐 API key managed via environment variable

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- OpenAI API key

### Installation

Clone the repo:

```bash
git clone https://github.com/your-username/kube-ai.git
cd kube-ai
```

Install dependencies:

```bash
go mod tidy
```

Build the app:

```bash
go build -o kube-ai
```

### Usage

Before using the CLI, make sure to set your OpenAI API key:

#### On Linux/macOS

```bash
export OPENAI_API_KEY=your_api_key_here
```

#### On Windows PowerShell

```powershell
$env:OPENAI_API_KEY="your_api_key_here"
```

Run the CLI:

```bash
./kube-ai "How to scale deployments in Kubernetes?"
```

## 🛠️ Project Structure

```bash
.
├── cmd
│   └── root.go         # Main CLI command logic
├── go.mod              # Module dependencies
├── go.sum              # Checksums
└── main.go             # Entry point
```

## 🧠 Powered by

- [OpenAI GPT](https://platform.openai.com/docs)
- [Cobra CLI](https://github.com/spf13/cobra)

## 📄 License

This project is licensed under the MIT License.
```

---
