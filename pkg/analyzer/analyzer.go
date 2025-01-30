package analyzer

import (
    "context"
    "fmt"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

func GetWorkloadYAML(client *kubernetes.Clientset, namespace, workloadType, name string) (string, error) {
    ctx := context.Background()
    
    deployment, err := client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil {
        return "", fmt.Errorf("failed to get deployment: %v", err)
    }

    // Extract important deployment details
    yamlInfo := fmt.Sprintf(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
spec:
  replicas: %d
  template:
    spec:
      containers:`, deployment.Name, deployment.Namespace, *deployment.Spec.Replicas)

    // Add container details
    for _, container := range deployment.Spec.Template.Spec.Containers {
        yamlInfo += fmt.Sprintf(`
      - name: %s
        image: %s
        resources:
          limits:
            cpu: %s
            memory: %s
          requests:
            cpu: %s
            memory: %s`,
            container.Name,
            container.Image,
            container.Resources.Limits.Cpu().String(),
            container.Resources.Limits.Memory().String(),
            container.Resources.Requests.Cpu().String(),
            container.Resources.Requests.Memory().String(),
        )
    }

    return yamlInfo, nil
}

func AnalyzeWorkload(client *kubernetes.Clientset, namespace, workloadType, name string, config *rest.Config) (*WorkloadDetails, error) {
    var podSpec *corev1.PodSpec

    // Get workload based on type
    switch workloadType {
    case "deployment":
        deployment, err := client.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
        podSpec = &deployment.Spec.Template.Spec
    case "statefulset":
        sts, err := client.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
        podSpec = &sts.Spec.Template.Spec
    case "daemonset":
        ds, err := client.AppsV1().DaemonSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
        if err != nil {
            return nil, err
        }
        podSpec = &ds.Spec.Template.Spec
    default:
        return nil, fmt.Errorf("unsupported workload type: %s", workloadType)
    }

    // Get metrics
    metrics, err := GetMetrics(client, namespace, name, config)
    if err != nil {
        return nil, fmt.Errorf("failed to get metrics: %v", err)
    }

    fmt.Printf("Debug - Metrics map: %+v\n", metrics) // Add this debug line

    // Get main container and QoS class
    var mainContainer string
    if len(podSpec.Containers) > 0 {
        mainContainer = podSpec.Containers[0].Name
    }

    details := &WorkloadDetails{
        Namespace:         namespace,
        Deployment:       name,
        Kind:            workloadType,
        MainContainer:    mainContainer,
        PodQoSClass:     string(podSpec.PriorityClassName),
        ReplicaCount:    metrics["replica_count"],
        CPUUtilization:  metrics["cpu_utilization"],
        MemoryUtilization: metrics["memory_utilization"],
        EfficiencyRate:   metrics["efficiency_rate"],  // This line is important
        ContainerCount:   fmt.Sprintf("%d", len(podSpec.Containers)),
        NetworkTraffic:   metrics["network_traffic"],
        OpsaniFlags:      "N/A",
    }

    fmt.Printf("Debug - WorkloadDetails: %+v\n", details) // Add this debug line

    return details, nil
}

func GetWorkloadMetrics(client *kubernetes.Clientset, namespace, name string) (map[string]string, error) {
    // Get HPA metrics
    hpa, err := client.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(context.Background(), name, metav1.GetOptions{})
    if err != nil {
        return map[string]string{
            "cpu_utilization": "N/A",
            "memory_utilization": "N/A",
            "replica_count": "N/A",
        }, nil
    }

    metrics := make(map[string]string)
    
    // Get current metrics
    if hpa.Status.CurrentMetrics != nil {
        for _, metric := range hpa.Status.CurrentMetrics {
            if metric.Resource != nil {
                switch metric.Resource.Name {
                case "cpu":
                    metrics["cpu_utilization"] = fmt.Sprintf("%d%%", metric.Resource.Current.AverageUtilization)
                case "memory":
                    metrics["memory_utilization"] = fmt.Sprintf("%d%%", metric.Resource.Current.AverageUtilization)
                }
            }
        }
    }

    // Get replica count
    metrics["replica_count"] = fmt.Sprintf("%d/%d", hpa.Status.CurrentReplicas, hpa.Status.DesiredReplicas)

    return metrics, nil
}