package analyzer

type WorkloadDetails struct {
    Namespace         string
    Deployment        string
    Kind             string
    MainContainer     string
    PodQoSClass      string
    ReplicaCount     string
    CPUUtilization   string
    MemoryUtilization string
    EfficiencyRate   string
    ReliabilityRisk  string
    ContainerCount   string    // Add this field
    NetworkTraffic   string    // Add this field
    OpsaniFlags      string    // Add this field
    Analysis         string
    Opportunities    []string
    Cautions        []string
    Blockers        []string
    Recommendations []string
}