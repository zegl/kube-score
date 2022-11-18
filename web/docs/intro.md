---
sidebar_position: 1
---

# Introduction to kube-score

Let's discover **kube-score in less than 5 minutes**.

**kube-score** is a CLI tool that does static code analysis of your Kubernetes object definitions. The output is a list of recommendations of what you can improve to make your application more secure and resilient.

**kube-score** is open-source and available under the **MIT-license**. For more information about how to use kube-score, see zegl/kube-score on GitHub.

## Getting Started

### Install kube-score

```bash
# Install with Homebrew (recommended)
brew install kube-score

# Install with krew
kubectl krew install score

# Install with Docker
docker pull zegl/kube-score
```

### Run your first kube-score

```bash
kube-score score your/path/to/*.yaml
```

Thanks it, congrats!