package ui

import (
    "fmt"
    "strings"
    "github.com/charmbracelet/lipgloss"
    "k8s-workload-analyzer/pkg/analyzer"
)

var (
    titleStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.DoubleBorder()).
        BorderForeground(lipgloss.Color("87")).
        Padding(0, 1)

    labelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("246")).
        Width(20)

    valueStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("252"))

    sectionStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("240")).
        Padding(1, 2)

    // Update existing styles
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("82"))  // Brighter green

    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("214")) // Brighter yellow

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")) // Brighter red
)

func RenderAnalysis(details *analyzer.WorkloadDetails) string {
    // Format basic info
    basicInfo := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s",
        labelStyle.Render("Namespace"),
        valueStyle.Render(details.Namespace),
        labelStyle.Render("Deployment"),
        valueStyle.Render(details.Deployment),
        labelStyle.Render("Kind"),
        valueStyle.Render(details.Kind),
        labelStyle.Render("Main Container"),
        valueStyle.Render(details.MainContainer),
    )

    // Format metrics
    metrics := fmt.Sprintf("%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s",
        labelStyle.Render("Replica Count"),
        valueStyle.Render(details.ReplicaCount),
        labelStyle.Render("CPU Utilization"),
        valueStyle.Render(details.CPUUtilization),
        labelStyle.Render("Memory Utilization"),
        valueStyle.Render(details.MemoryUtilization),
        labelStyle.Render("Container Count"),
        valueStyle.Render(details.ContainerCount),
        labelStyle.Render("Efficiency Rate"),
        formatEfficiencyRate(details.EfficiencyRate),
    )

    // Format analysis sections
    analysis := fmt.Sprintf("%s: %s\n%s: %s",
        labelStyle.Render("Reliability Risk"),
        formatReliabilityRisk(details.ReliabilityRisk),
        labelStyle.Render("Analysis"),
        valueStyle.Render(details.Analysis),
    )

    return fmt.Sprintf(`
%s

%s

%s

%s

%s

%s

%s

%s`,
        titleStyle.Render("Workload Analysis"),
        sectionStyle.Render(basicInfo),
        sectionStyle.Render(metrics),
        sectionStyle.Render(analysis),
        formatSection("Opportunities", details.Opportunities, successStyle),
        formatSection("Cautions", details.Cautions, warningStyle),
        formatSection("Blockers", details.Blockers, errorStyle),
        formatSection("Recommendations", details.Recommendations, successStyle),
    )
}

func formatEfficiencyRate(rate string) string {
    if strings.Contains(rate, "High") {
        return successStyle.Render(rate)
    } else if strings.Contains(rate, "Medium") {
        return warningStyle.Render(rate)
    }
    return errorStyle.Render(rate)
}

func formatReliabilityRisk(risk string) string {
    if risk == "Low" {
        return successStyle.Render(risk)
    } else if risk == "Medium" {
        return warningStyle.Render(risk)
    }
    return errorStyle.Render(risk)
}

func formatSection(title string, items []string, style lipgloss.Style) string {
    if len(items) == 0 {
        return sectionStyle.Render(fmt.Sprintf("%s:\nNone", labelStyle.Render(title)))
    }

    content := fmt.Sprintf("%s:", labelStyle.Render(title))
    for _, item := range items {
        content += fmt.Sprintf("\n• %s", style.Render(item))
    }
    return sectionStyle.Render(content)
}

func renderBasicInfo(details *analyzer.WorkloadDetails) string {
    return fmt.Sprintf(`
Namespace           : %s
Deployment          : %s
Kind               : %s
Main Container     : %s
Pod QoS Class      : %s
Average Replica Count: %s
Container Count    : %d`,
        details.Namespace,
        details.Deployment,
        details.Kind,
        details.MainContainer,
        details.PodQoSClass,
        details.ReplicaCount,
        details.ContainerCount,
    )
}

func renderMetrics(details *analyzer.WorkloadDetails) string {
    return fmt.Sprintf(`
CPU Utilization    : %s
Memory Utilization : %s
Network Traffic    : %s
Opsani Flags       : %s

Efficiency Rate    : %s
Reliability Risk   : %s
Analysis          : %s`,
        details.CPUUtilization,
        details.MemoryUtilization,
        details.NetworkTraffic,
        details.OpsaniFlags,
        successStyle.Render(details.EfficiencyRate),
        errorStyle.Render(details.ReliabilityRisk),
        errorStyle.Render(details.Analysis),
    )
}

func renderSection(title string, items []string, style lipgloss.Style) string {
    if len(items) == 0 {
        return ""
    }

    var lines []string
    lines = append(lines, fmt.Sprintf("%-18s: %s", title, style.Render(items[0])))
    
    for _, item := range items[1:] {
        lines = append(lines, fmt.Sprintf("%-18s  %s", "", style.Render(item)))
    }

    return strings.Join(lines, "\n")
}

func formatList(items []string) string {
    if len(items) == 0 {
        return "None"
    }
    
    var formatted []string
    for _, item := range items {
        formatted = append(formatted, fmt.Sprintf("- %s", item))
    }
    return strings.Join(formatted, "\n")
}