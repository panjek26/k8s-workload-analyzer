package ai

type WorkloadAnalysis struct {
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