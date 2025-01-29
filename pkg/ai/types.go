package ai

type WorkloadAnalysis struct {
    Analysis         string `json:"analysis"`
    Recommendations []struct {
        Category        string `json:"category"`
        Description    string `json:"description"`
        Severity       string `json:"severity"`
        SuggestedAction string `json:"suggested_action"`
    } `json:"recommendations"`
}