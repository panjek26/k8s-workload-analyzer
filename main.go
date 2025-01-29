package main

import (
    "context"
    "log"
    "net"
    "sync"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    healthpb "google.golang.org/grpc/health/grpc_health_v1"
    "google.golang.org/grpc/reflection"
    pb "github.com/panjek26/k8s-workload-analyzer/proto"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"

    "github.com/panjek26/k8s-workload-analyzer/config"
    "github.com/panjek26/k8s-workload-analyzer/deepseek"
)

const (
    port = ":8080"
)

// HealthChecker handles the health check state
type HealthChecker struct {
    healthServer *health.Server
    mu          sync.Mutex
    ready       bool
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
    return &HealthChecker{
        healthServer: health.NewServer(),
        ready:       false,
    }
}

// SetReady sets the readiness state
func (h *HealthChecker) SetReady(ready bool) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.ready = ready
    if ready {
        h.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
    } else {
        h.healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
    }
}

type WorkloadAnalyzerServer struct {
    pb.UnimplementedWorkloadAnalyzerServer
    k8sClient *kubernetes.Clientset
}

func NewWorkloadAnalyzerServer() (*WorkloadAnalyzerServer, error) {
    config, err := rest.InClusterConfig()
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, err
    }

    return &WorkloadAnalyzerServer{
        k8sClient: clientset,
    }, nil
}

func (s *WorkloadAnalyzerServer) AnalyzeWorkload(ctx context.Context, req *pb.WorkloadRequest) (*pb.WorkloadAnalysis, error) {
    // Initialize Deepseek client (you'll need to implement this based on their API)
    deepseekClient := NewDeepseekClient()

    // Collect workload metrics from Kubernetes
    metrics, err := s.collectWorkloadMetrics(req.Namespace, req.DeploymentName)
    if err != nil {
        return nil, err
    }

    // Prepare prompt for Deepseek
    prompt := prepareAnalysisPrompt(metrics, req.Metrics)

    // Get analysis from Deepseek
    analysis, err := deepseekClient.Analyze(ctx, prompt)
    if err != nil {
        return nil, err
    }

    // Convert Deepseek analysis to response
    return &pb.WorkloadAnalysis{
        AnalysisResult:   analysis.Result,
        ConfidenceScore: analysis.Confidence,
        Recommendations: convertRecommendations(analysis.Recommendations),
    }, nil
}

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    lis, err := net.Listen("tcp", cfg.GRPCPort)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // Create Deepseek client
    deepseekClient := deepseek.NewDeepseekClient(cfg.DeepseekAPIKey)

    // Create a new gRPC server
    server := grpc.NewServer()

    // Create and register health checker
    healthChecker := NewHealthChecker()
    healthpb.RegisterHealthServer(server, healthChecker.healthServer)

    // Create and register workload analyzer
    workloadAnalyzer, err := NewWorkloadAnalyzerServer()
    if err != nil {
        log.Fatalf("failed to create workload analyzer: %v", err)
    }
    pb.RegisterWorkloadAnalyzerServer(server, workloadAnalyzer)

    // Enable reflection for grpcurl and other tools
    reflection.Register(server)

    // Initially set as not ready
    healthChecker.SetReady(false)

    // Simulate startup delay (replace with your actual startup logic)
    go func() {
        time.Sleep(5 * time.Second)
        healthChecker.SetReady(true)
        log.Printf("Service is now ready")
    }()

    log.Printf("Health check server listening at %v", lis.Addr())
    if err := server.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
