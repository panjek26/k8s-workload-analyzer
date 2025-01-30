# Kubernetes Workload Analyzer

A command-line tool that analyzes Kubernetes workloads and provides detailed insights about resource utilization, efficiency, and best practices recommendations.

## Features

- Real-time resource metrics analysis
- Efficiency rate calculation
- AI-powered configuration analysis
- Best practices recommendations
- Security insights
- Performance optimization suggestions

## Prerequisites

- Go 1.19 or higher
- Access to a Kubernetes cluster
- OpenAI API key
- Kubernetes metrics server installed in your cluster

## Installation

```bash
git clone https://github.com/yourusername/k8s-workload-analyzer.git
cd k8s-workload-analyzer
go build -o kwa cmd/main.go