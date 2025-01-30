package prompts

const WorkloadAnalysisTemplate = `As a Kubernetes expert, analyze this container configuration and provide a detailed assessment. Return a valid JSON object with comprehensive insights:
{
    "main_container": "container-name",
    "pod_qos_class": "qos-class",
    "replica_count": "current/desired",
    "cpu_utilization": "percentage",
    "memory_utilization": "percentage",
    "efficiency_rate": "efficiency-value",
    "reliability_risk": "risk-level",
    "analysis": "Detailed analysis of the workload configuration, focusing on:
        - Resource allocation and utilization patterns
        - Configuration quality and compliance with best practices
        - Potential performance bottlenecks
        - Security posture assessment
        - High availability and reliability considerations",
    "opportunities": [
        "Resource optimization opportunities with specific metrics",
        "Performance improvement suggestions with clear benefits",
        "Cost optimization strategies with estimated savings",
        "Scalability enhancements with concrete recommendations",
        "Security improvements with best practices references"
    ],
    "cautions": [
        "Resource constraints with impact analysis",
        "Configuration risks with potential consequences",
        "Security vulnerabilities with severity levels",
        "Scaling limitations with specific thresholds",
        "Reliability concerns with mitigation strategies"
    ],
    "recommendations": [
        "Specific, actionable steps for resource optimization",
        "Detailed configuration improvements with examples",
        "Security hardening measures with implementation guides",
        "Performance tuning suggestions with expected outcomes",
        "Scaling strategy improvements with configuration samples"
    ]
}

Container configuration to analyze:
%s`