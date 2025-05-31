# Dingo

An opinionated templating CLI for Infrastructure as Data workflows.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Overview

Dingo is a Go-based templating CLI inspired by Kapitan, designed to generate configuration files, infrastructure code, and scripts using a data-driven approach. It combines YAML data management with powerful Go templating enhanced by the Sprig function library.

## ✨ Key Features

### 🗂️ Data Management
-  Base/overlay environment configs
-  Automatic YAML merging
-  CUE-based validation

### 🎯 Templating
-  Go templates with 100+ Sprig functions

### 🔐 Secrets
-  Multiple backends (incl. GCP Secret Manager)
-  `$$secretname$$` syntax

### 🛡️ Validation
-  CUE schemas with type checks


## 🚀 Quick Start

### Installation
```bash
git clone https://github.com/alxndr13/dingo
cd dingo
make build
```

### Basic Usage
```bash
# Generate templates with default settings
./bin/dingo

# Specify custom paths
./bin/dingo --basepath ./data/base \
           --overlaypath ./data/overlays/prod \
           --templatepath ./templates

# Enable secret decryption
./bin/dingo --decryptor google
```

## 📁 Project Structure

```
project/
├── data/
│   ├── base/                 # Base configuration files
│   │   └── data.yaml
│   └── overlays/             # Environment-specific overrides
│       ├── dev/
│       │   └── data.yaml
│       └── prod/
│           └── data.yaml
├── templates/                # Go template files
│   ├── terraform/
│   │   └── main.tf
│   └── kubernetes/
│       └── deployment.yaml
├── output/                   # Generated files (auto-created)
└── schema.cue               # Data validation schema
```

## 🎨 Templating with Sprig

Dingo includes the full [Sprig function library](https://masterminds.github.io/sprig/), providing 100+ utility functions for advanced templating.

### String Functions
```yaml
# data.yaml
app:
  name: "my-awesome-app"
  environment: "production"
```

```go
# template.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .app.name | kebabcase }}
  labels:
    app: {{ .app.name | quote }}
    env: {{ .app.environment | upper }}
    version: {{ .app.name | sha256sum | trunc 8 }}
```

### Math & Logic Functions
```yaml
# data.yaml
replicas:
  min: 2
  max: 10
resources:
  cpu: 500
  memory: 1024
```

```go
# template.yaml
spec:
  replicas: {{ max .replicas.min 3 }}
  template:
    spec:
      containers:
      - name: app
        resources:
          requests:
            cpu: {{ .resources.cpu }}m
            memory: {{ .resources.memory }}Mi
          limits:
            cpu: {{ mul .resources.cpu 2 }}m
            memory: {{ mul .resources.memory 1.5 | int }}Mi
```

### Date & Time Functions
```go
# Generated timestamp
creationTimestamp: {{ now | date "2006-01-02T15:04:05Z" }}
# Expiry date (30 days from now)
expiryDate: {{ now | dateModify "+720h" | date "2006-01-02" }}
```

### List & Dictionary Functions
```yaml
# data.yaml
services:
  - name: "web"
    port: 80
  - name: "api"
    port: 8080
  - name: "db"
    port: 5432
```

```go
# template.yaml
services:
{{- range .services }}
  - name: {{ .name }}
    port: {{ .port }}
{{- end }}

# First and last services
primary: {{ (first .services).name }}
backup: {{ (last .services).name }}

# Join service names
all_services: {{ .services | pluck "name" | join "," }}
```

### Advanced Sprig Examples

```go
# Conditional logic with defaults
database_url: {{ .database.url | default "localhost:5432" }}

# Complex string manipulation
config_name: {{ printf "%s-%s" .app.name .environment | lower | replace "_" "-" }}

# Random generation for testing
test_password: {{ randAlphaNum 16 }}
test_uuid: {{ uuidv4 }}

# Base64 encoding
secret_data: {{ .secret | b64enc }}

# URL manipulation
api_endpoint: {{ .base_url | trimSuffix "/" }}/api/v1

# File operations in templates
{{- if .features.monitoring }}
monitoring_config: |
{{ .monitoring | toYaml | indent 2 }}
{{- end }}
```

## 🔐 Secret Management

### Google Secret Manager
```bash
# Setup authentication
gcloud auth application-default login

# Use in templates
./bin/dingo --decryptor google
```

```yaml
# data.yaml
database:
  host: "prod-db.example.com"
  password: "$$projects/alxndr13/secrets/password/versions/latest$$"  # Retrieved from Secret Manager
api:
  key: "$$projects/alxndr13/secrets/api-key/versions/latest$$"
```

### Custom Decryptors
Implement the `Decryptor` interface for other secret backends:
```go
type Decryptor interface {
    Init() error
    Decrypt(secretName string) (string, error)
}
```

## 📋 CLI Options

| Flag | Default | Description |
|------|---------|-------------|
| `--basepath` | `data/base` | Base directory for YAML files |
| `--overlaypath` | `data/overlays/dev` | Overlay directory for environment-specific data |
| `--templatepath` | `templates` | Directory containing template files |
| `--logmode` | `human` | Logging mode (`human` or `json`) |
| `--decryptor` | (none) | Secret decryptor (`example` or `google`) |

## 🧪 Development

### Running Tests
```bash
make test                    # Run unit tests
make integration_tests       # Run integration tests (requires GCP auth)
```

### Building
```bash
make build                   # Build binary to ./bin/dingo
make run                     # Run with default settings
```

## 🤝 Background

Dingo was born from real-world experience with Kapitan and the need for a Go-native solution that's easy to distribute and integrate. At my last employer, I was introduced to [Kapitan](kapitan.dev), which became the foundation for many projects by enabling multiple environments regardless of technology. We mainly used it to generate Terraform HCL, bash scripts, HAProxy configurations, PostgreSQL scripts, Kubernetes manifests, and more.

What I disliked about Kapitan was its Python ecosystem and the hassle of installing it on Docker images or new machines. You always had to ensure the exact Python version, or something would break.

When tasked with building a new foundation for a client, my colleagues and I chose to build on a [custom Go library](https://github.com/lukasjarosch/skipper) inspired by Kapitan (Shoutout to my friend [@lukasjarosch](https://github.com/lukasjarosch) here!). We even presented our "Infrastructure as Data" approach [at a local German conference](https://www.continuouslifecycle.de/veranstaltung-21248-0-wie-wir-mit-infrastructure-as-data-eine-plattform-gebaut-haben.html).
