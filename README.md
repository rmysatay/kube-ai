Elbette! Aşağıya `version` ve `completion` komutlarını da içeren güncellenmiş `README.md` sürümünü ekliyorum:

---

````markdown
# 🤖 kube-ai

> AI-powered Kubernetes CLI Assistant built in Go

**kube-ai** is a command-line tool that brings the power of AI (OpenAI GPT models) to Kubernetes. It helps DevOps, SRE, and platform engineers analyze, audit, generate, and troubleshoot Kubernetes resources using natural language.

---

## 🚀 Features

- 🔍 `analyze`: Analyze Kubernetes outputs (logs, describe, YAML) and ask questions
- 🔐 `audit`: Detect security risks and misconfigurations in YAML or live resources
- 🛠️ `diagnose`: Find root causes of pod failures (e.g., OOMKilled, ImagePullBackOff)
- 🧾 `generate`: Create YAML manifests with natural language prompts
- ✏️ `modify`: Edit existing YAML files (namespace, name, replicas)
- 💬 `chat`: Ask how-to questions and get CLI-based guidance
- ⚡ `execute`: Apply a manifest to the cluster
- 📜 `history`: View previously used commands and inputs
- ✅ `completion`: Generate shell autocompletions (bash, zsh, fish, powershell)
- 🔢 `version`: Show current CLI version
- 🔧 Customizable flags: AI model, max tokens, verbosity, etc.

---

## 📦 Installation

```bash
git clone https://github.com/rmysatay/kube-ai.git
cd kube-ai
go build -o kube-ai .
````

---

## 🔑 Prerequisites

* Go 1.20+
* `kubectl` installed and configured
* OpenAI API key (GPT-4o by default)

```bash
export OPENAI_API_KEY=your-api-key
```

Or create a `.env` file:

```
OPENAI_API_KEY=your-api-key
```

---

## 🧪 Usage Examples

### 🔍 Analyze a pod

```bash
kube-ai analyze pod nginx --ns default "Why is this crashing?"
```

### 🔐 Audit a resource

```bash
kube-ai audit --file deployment.yaml
```

### 🛠 Diagnose pod issues

```bash
kube-ai diagnose --name pod/nginx --ns default
```

### 🧾 Generate YAML from natural language

```bash
kube-ai generate "Create a Service and Deployment for redis" --replicas 1 --save
```

### ✏️ Modify YAML

```bash
kube-ai modify --file deploy.yaml --namespace prod --replicas 3 --name webapp
```

### ⚡ Apply a manifest

```bash
kube-ai execute --file output.yaml
```

### 💬 Ask CLI questions

```bash
kube-ai chat "How can I restart a deployment in Kubernetes?"
```

### 📜 View history

```bash
kube-ai history
```

---

## ✅ Completion

You can enable shell autocompletion to speed up usage:

```bash
# For Bash
source <(kube-ai completion bash)

# For Zsh
source <(kube-ai completion zsh)

# For Fish
kube-ai completion fish | source

# For PowerShell
kube-ai completion powershell | Out-String | Invoke-Expression
```

You can also add the appropriate command to your shell config file (`.bashrc`, `.zshrc`, etc.) to load on startup.

---

## 🔢 Version

To display the current version of `kube-ai`:

```bash
kube-ai version
```

This will print:

```
Kube-AI Version: 0.1.0
```

---

## 🧠 Global Flags

You can configure most commands using these flags:

```bash
--model, -m           # Choose AI model (default: gpt-4o)
--max-tokens, -t      # Max tokens per response (default: 2048)
--count-tokens, -c    # Show token usage after each command
--verbose, -v         # Show detailed debug output
--max-iterations, -x  # Reserved for advanced multi-step flows
```

---

## 🗂 Roadmap

* [ ] Helm chart analysis support
* [ ] Multi-resource batch audits
* [ ] Language localization (e.g., Turkish 🇹🇷)
* [ ] VS Code integration

---
