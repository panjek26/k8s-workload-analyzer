package analyzer

import (
	"context"
	"fmt"

	"github.com/panjek26/k8s-workload-analyzer/pkg/api"
	"github.com/panjek26/k8s-workload-analyzer/pkg/models"
)

type AnalyzerService struct {
	api.UnimplementedWorkloadAnalyzerServer
	aiClient *models.DeepSeekClient
}

func NewAnalyzerService(apiKey string) *AnalyzerService {
	return &AnalyzerService{
		aiClient: models.NewDeepSeekClient(apiKey),
	}
}

func (s *AnalyzerService) AnalyzeWorkload(ctx context.Context, req *api.AnalyzeRequest) (*api.AnalyzeResponse, error) {
	// Prepare the prompt for AI analysis
	prompt := fmt.Sprintf(
		"Analyze Kubernetes workload:\nNamespace: %s\nType: %s\nName: %s\nMetrics: %v\n"+
			"Provide optimization recommendations based on the metrics.",
		req.Namespace, req.WorkloadType, req.WorkloadName, req.Metrics,
	)

	// Get AI analysis
	analysis, err := s.aiClient.GetAnalysis(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI analysis: %v", err)
	}

	// Process AI response and create recommendations
	recommendations := processAIResponse(analysis)

	return &api.AnalyzeResponse{
		Analysis:         analysis,
		Recommendations: recommendations,
	}, nil
}

func processAIResponse(analysis string) []*api.Recommendation {
	// This is a placeholder. You should implement proper parsing of AI response
	// and convert it into structured recommendations
	return []*api.Recommendation{
		{
			Category:        "Resource Optimization",
			Description:     analysis,
			ConfidenceScore: 0.85,
		},
	}
}