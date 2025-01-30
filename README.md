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

./kwa -namespace=<namespace> -name=<deployment-name> -type=deployment -api-key=<your-openai-api-key>
```

### Parameters
- `-namespace` : Kubernetes namespace of the workload
- `-name` : Name of the workload (deployment, statefulset, etc.)
- `-type` : Type of workload (deployment, statefulset, daemonset)
- `-api-key` : OpenAI API key for AI analysis

### Example Output

```
╭─────────────────────╮
│  Workload Analysis  │
╰─────────────────────╯

┌──────────────────────────────────────┐
│ Basic Information                    │
├──────────────────────────────────────┤
│ Namespace      : default             │
│ Deployment     : my-app              │
│ Kind          : deployment           │
│ Main Container : app                 │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│ Resource Metrics                     │
├──────────────────────────────────────┤
│ CPU Utilization    : 45m             │
│ Memory Utilization : 256Mi           │
│ Efficiency Rate    : High (85.2%)    │
│ Container Count    : 2               │
└──────────────────────────────────────┘
```

## Features in Detail

### Resource Metrics Analysis
- Real-time CPU and memory utilization
- Resource efficiency calculation
- Container resource usage patterns

### Configuration Analysis
- Best practices validation
- Security configuration review
- Resource allocation assessment

### Recommendations
- Resource optimization suggestions
- Performance improvement tips
- Security enhancement recommendations

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
MIT License - feel free to use and modify this tool for your needs.

