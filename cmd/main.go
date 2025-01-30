package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s-workload-analyzer/pkg/analyzer"
    "k8s-workload-analyzer/pkg/ai"
    "k8s-workload-analyzer/pkg/ui"
)

func main() {
    namespace := flag.String("namespace", "", "Kubernetes namespace")
    workloadType := flag.String("type", "deployment", "Workload type (deployment, statefulset, daemonset)")
    workloadName := flag.String("name", "", "Workload name")
    apiKey := flag.String("api-key", "", "GPT API key")
    flag.Parse()

    if *namespace == "" || *workloadName == "" || *apiKey == "" {
        flag.Usage()
        os.Exit(1)
    }

    // Initialize Kubernetes client
    config, err := clientcmd.BuildConfigFromFlags("", clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename())
    if err != nil {
        log.Fatalf("Failed to get kubeconfig: %v", err)
    }

    k8sClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create kubernetes client: %v", err)
    }

    // Get workload YAML and analyze
    yaml, err := analyzer.GetWorkloadYAML(k8sClient, *namespace, *workloadType, *workloadName)
    if err != nil {
        log.Fatalf("Failed to get workload: %v", err)
    }

    // Initialize AI client and analyze
    aiClient := ai.NewGPTClient(*apiKey)
    analysis, err := aiClient.AnalyzeWorkload(yaml)
    if err != nil {
        log.Fatalf("Failed to analyze workload: %v", err)
    }

    // Get workload details with metrics
    details, err := analyzer.AnalyzeWorkload(k8sClient, *namespace, *workloadType, *workloadName, config)
    if err != nil {
        log.Printf("Warning: Failed to get workload details: %v", err)
    }

    // Debug output
    fmt.Printf("Debug - Before AI analysis - Efficiency Rate: %s\n", details.EfficiencyRate)

    // Update details with AI analysis
    if details != nil {
        // Don't overwrite efficiency rate from metrics
        // details.EfficiencyRate = analysis.EfficiencyRate
        
        details.ReliabilityRisk = analysis.ReliabilityRisk
        details.Analysis = analysis.Analysis
        details.Opportunities = analysis.Opportunities
        details.Cautions = analysis.Cautions
        details.Blockers = analysis.Blockers
        details.Recommendations = analysis.Recommendations
    }

    fmt.Printf("Debug - After AI analysis - Efficiency Rate: %s\n", details.EfficiencyRate)

    // Render and display results
    fmt.Println(ui.RenderAnalysis(details))
}