package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"
    "bytes"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"  // Add this import
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/runtime/serializer/json"
    "k8s.io/client-go/kubernetes/scheme"
    runtime "k8s.io/apimachinery/pkg/runtime"
    pb "k8s-workload-analyzer/api/proto"
    "k8s-workload-analyzer/pkg/ai"
)

type server struct {
    pb.UnimplementedWorkloadAnalyzerServer
    k8sClient *kubernetes.Clientset
    aiClient  *ai.GPTClient
}

func (s *server) AnalyzeWorkload(ctx context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
    // Get workload YAML
    var obj interface{}
    var err error
    
    switch req.WorkloadType {
    case "deployment":
        obj, err = s.k8sClient.AppsV1().Deployments(req.Namespace).Get(ctx, req.WorkloadName, metav1.GetOptions{})
    case "statefulset":
        obj, err = s.k8sClient.AppsV1().StatefulSets(req.Namespace).Get(ctx, req.WorkloadName, metav1.GetOptions{})
    case "daemonset":
        obj, err = s.k8sClient.AppsV1().DaemonSets(req.Namespace).Get(ctx, req.WorkloadName, metav1.GetOptions{})
    default:
        return nil, fmt.Errorf("unsupported workload type: %s", req.WorkloadType)
    }

    if err != nil {
        return nil, fmt.Errorf("failed to get workload: %v", err)
    }

    // Convert to YAML
    serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
    var yamlData bytes.Buffer
    err = serializer.Encode(obj.(runtime.Object), &yamlData)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal workload to YAML: %v", err)
    }

    // Analyze with Claude
    analysis, err := s.aiClient.AnalyzeWorkload(yamlData.String(), req.Namespace, req.WorkloadType, req.WorkloadName)
    if err != nil {
        return nil, fmt.Errorf("analysis failed: %v", err)
    }

    // Convert to gRPC response
    recommendations := make([]*pb.Recommendation, len(analysis.Recommendations))
    for i, rec := range analysis.Recommendations {
        recommendations[i] = &pb.Recommendation{
            Category:        rec.Category,
            Description:    rec.Description,
            Severity:       rec.Severity,
            SuggestedAction: rec.SuggestedAction,
        }
    }

    return &pb.AnalyzeResponse{
        Analysis:        analysis.Analysis,
        Recommendations: recommendations,
    }, nil
}

func main() {
    port := flag.String("port", "50051", "The server port")
    apiKey := flag.String("api-key", "", "GPT API Key")
    flag.Parse()

    if *apiKey == "" {
        log.Fatal("API key is required")
    }

    // Setup Kubernetes client
    kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Failed to get kubeconfig: %v", err)
    }

    k8sClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    // Create server
    s := &server{
        k8sClient: k8sClient,
        aiClient:  ai.NewGPTClient(*apiKey),
    }

    lis, err := net.Listen("tcp", ":"+*port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterWorkloadAnalyzerServer(grpcServer, s)
    
    // Enable reflection
    reflection.Register(grpcServer)
    
    log.Printf("Server listening on port %s", *port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}