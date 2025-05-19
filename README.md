# 🧠 kube-ai

> AI-powered Kubernetes CLI assistant for analysis, auditing, YAML generation, troubleshooting, and more.

`kube-ai` is a smart Kubernetes CLI assistant powered by OpenAI. It helps developers, platform engineers, and SREs to quickly generate manifests, analyze pod issues, audit security misconfigurations, and troubleshoot live clusters with natural language input.

---

## 🚀 Features

- 🔍 `analyze`: Ask AI to explain Kubernetes YAML, logs, or command outputs.
- 🧪 `audit`: Perform security audits of resources for best practices.
- 🛠️ `diagnose`: Detect issues like `CrashLoopBackOff`, `ImagePullBackOff`, `OOMKilled`, etc.
- ✍️ `generate`: Generate Kubernetes YAML manifests from plain English.
- 🪄 `modify`: Edit fields like namespace, name, and replicas in YAML files.
- 📦 `execute`: Apply Kubernetes manifests to your cluster.
- 💬 `chat`: Ask any K8s-related question and get CLI examples.
- 🧠 `completion`: Auto-generate shell completions (bash, zsh, fish, powershell).
- 🕘 `history`: View your past `kube-ai` commands.

---

## 📦 Installation

```bash
git clone https://github.com/rmysatay/kube-ai.git
cd kube-ai
go build -o kube-ai
