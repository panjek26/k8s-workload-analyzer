package ai

import (
    "encoding/json"
    "strings"
)

type gptAnalysis struct {
    MainContainer     string   `json:"main_container"`
    PodQoSClass      string   `json:"pod_qos_class"`
    ReplicaCount     string   `json:"replica_count"`
    CPUUtilization   string   `json:"cpu_utilization"`
    MemoryUtilization string  `json:"memory_utilization"`
    EfficiencyRate   string   `json:"efficiency_rate"`
    ReliabilityRisk  string   `json:"reliability_risk"`
    Analysis         string   `json:"analysis"`
    Opportunities    []string `json:"opportunities"`
    Cautions        []string `json:"cautions"`
    Blockers        []string `json:"blockers"`
    Recommendations []string `json:"recommendations"`
}

func parseGPTResponse(content string) (*WorkloadAnalysis, error) {
    var analysis gptAnalysis
    if err := json.NewDecoder(strings.NewReader(content)).Decode(&analysis); err != nil {
        return nil, err
    }

    return &WorkloadAnalysis{
        MainContainer:     analysis.MainContainer,
        PodQoSClass:      analysis.PodQoSClass,
        ReplicaCount:     analysis.ReplicaCount,
        CPUUtilization:   analysis.CPUUtilization,
        MemoryUtilization: analysis.MemoryUtilization,
        EfficiencyRate:   analysis.EfficiencyRate,
        ReliabilityRisk:  analysis.ReliabilityRisk,
        Analysis:         analysis.Analysis,
        Opportunities:    analysis.Opportunities,
        Cautions:        analysis.Cautions,
        Blockers:        analysis.Blockers,
        Recommendations: analysis.Recommendations,
    }, nil
}