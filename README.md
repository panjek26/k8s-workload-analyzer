# Kubernetes Workload Analyzer

A gRPC service that analyzes Kubernetes workloads using AI to provide security, performance, and best practices recommendations.

## Features

- **Workload Analysis:** Supports Deployments, StatefulSets, and DaemonSets.
- **Recommendations:** Provides insights for:
  - Security enhancements
  - Performance optimizations
  - Kubernetes best practices
- **gRPC Interface:** Enables seamless integration.
- **Multi-Namespace Support:** Analyze workloads across different namespaces.
- **AI-Powered Insights:** Uses GPT to generate recommendations.

## Prerequisites

- Go 1.19 or later
- Access to a Kubernetes cluster
- GPT API key
- `grpcurl` (optional, for testing)

## Installation

```sh
git clone https://github.com/panjek26/k8s-workload-analyzer.git
cd k8s-workload-analyzer
go mod tidy
```

## Running the Server

```sh
go run cmd/main.go -api-key "your-gpt-api-key" -port 50052
```

## Usage Example

```sh
grpcurl -plaintext -d '{
  "namespace": "sit",
  "workload_type": "deployment",
  "workload_name": "your-deployment-name"
}' localhost:50052 analyzer.WorkloadAnalyzer/AnalyzeWorkload
```

## License

This project is licensed under the MIT License.
