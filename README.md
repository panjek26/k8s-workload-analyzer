# Kubernetes Workload Analyzer

A gRPC service that analyzes Kubernetes workloads using AI to provide security, performance, and best practices recommendations.

## Features

- Analyzes Kubernetes Deployments, StatefulSets, and DaemonSets
- Provides detailed analysis and recommendations for:
  - Security improvements
  - Performance optimizations
  - Kubernetes best practices
- gRPC interface for easy integration
- Supports multiple namespaces
- AI-powered analysis using GPT

## Prerequisites

- Go 1.19 or later
- Access to a Kubernetes cluster
- GPT API key
- `grpcurl` (optional, for testing)

## Installation

```bash
git clone https://github.com/panjek26/k8s-workload-analyzer.git
cd k8s-workload-analyzer
go mod tidy

## Running the Server

```bash
go run cmd/main.go -api-key "your-gpt-api-key" -port 50052


## Usage Example
```bash
grpcurl -plaintext -d '{
  "namespace": "sit",
  "workload_type": "deployment",
  "workload_name": "your-deployment-name"
}' localhost:50052 analyzer.WorkloadAnalyzer/AnalyzeWorkload