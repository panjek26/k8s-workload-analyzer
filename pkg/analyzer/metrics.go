package analyzer

import (
    "context"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

func GetMetrics(client *kubernetes.Clientset, namespace, name string, config *rest.Config) (map[string]string, error) {
    // Create metrics client
    metricsClient, err := metricsv.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create metrics client: %v", err)
    }

    // Get deployment to find pod selector
    deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
        return nil, fmt.Errorf("failed to get deployment: %v", err)
    }

    // Get pods using deployment's selector
    selector := metav1.FormatLabelSelector(deployment.Spec.Selector)
    pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
        LabelSelector: selector,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list pods: %v", err)
    }

    if len(pods.Items) == 0 {
        return map[string]string{
            "cpu_utilization": "N/A",
            "memory_utilization": "N/A",
            "replica_count": "0",
        }, nil
    }

    var totalCPU, totalMemory int64
    podCount := 0

    // Get metrics for each pod
    for _, pod := range pods.Items {
        fmt.Printf("Getting metrics for pod: %s\n", pod.Name)
        podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), pod.Name, metav1.GetOptions{})
        if err != nil {
            fmt.Printf("Error getting metrics for pod %s: %v\n", pod.Name, err)
            continue
        }

        // Calculate total resource usage
        for _, container := range podMetrics.Containers {
            cpuQuantity := container.Usage.Cpu()
            memQuantity := container.Usage.Memory()
            
            totalCPU += cpuQuantity.MilliValue()
            totalMemory += memQuantity.Value()
        }
        podCount++
    }

    metrics := make(map[string]string)
    // Set default values
    metrics["efficiency_rate"] = "N/A"
    
    if podCount > 0 {
        avgCPU := totalCPU / int64(podCount)
        avgMemory := totalMemory / int64(podCount)
        
        metrics["cpu_utilization"] = fmt.Sprintf("%dm", avgCPU)
        
        // Format memory in Gi if over 1000Mi
        memoryMi := avgMemory / (1024 * 1024)
        if memoryMi >= 1000 {
            metrics["memory_utilization"] = fmt.Sprintf("%.2fGi", float64(memoryMi)/1024.0)
        } else {
            metrics["memory_utilization"] = fmt.Sprintf("%dMi", memoryMi)
        }

        // Calculate efficiency rate based on resource usage vs requests
        if len(deployment.Spec.Template.Spec.Containers) > 0 {
            container := deployment.Spec.Template.Spec.Containers[0]
            
            // Get CPU request
            cpuRequest := container.Resources.Requests["cpu"]
            memRequest := container.Resources.Requests["memory"]
            
            if !cpuRequest.IsZero() && !memRequest.IsZero() {
                cpuRequestValue := cpuRequest.MilliValue()
                memRequestValue := memRequest.Value() / (1024 * 1024)
                
                cpuEfficiency := float64(avgCPU) / float64(cpuRequestValue) * 100
                memEfficiency := float64(memoryMi) / float64(memRequestValue) * 100
                avgEfficiency := (cpuEfficiency + memEfficiency) / 2

                var efficiencyLevel string
                if avgEfficiency < 50 {
                    efficiencyLevel = "Low"
                } else if avgEfficiency < 80 {
                    efficiencyLevel = "Medium"
                } else {
                    efficiencyLevel = "High"
                }
                
                metrics["efficiency_rate"] = fmt.Sprintf("%s (%.1f%%)", efficiencyLevel, avgEfficiency)
            }
        } else {
            metrics["efficiency_rate"] = "N/A (No containers)"
        }
        
        metrics["replica_count"] = fmt.Sprintf("%d", podCount)
    } else {
        metrics["efficiency_rate"] = "N/A (No running pods)"
    }

    // Get HPA info if available
    hpa, err := client.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(context.Background(), name, metav1.GetOptions{})
    if err == nil {
        metrics["replica_count"] = fmt.Sprintf("%d/%d", hpa.Status.CurrentReplicas, hpa.Status.DesiredReplicas)
    }

    return metrics, nil
}