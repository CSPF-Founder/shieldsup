package services

import (
	"strings"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"
)

// DefaultRemediation constant
const DefaultRemediation = `For vulnerabilities related to external software or third-party components, ` +
	`it is advisable to reach out to the respective software vendor. ` +
	`In such cases, apply the patches or updates provided by the vendor ` +
	`to effectively remediate the vulnerability.
	
For vulnerabilities related to the application code, ` +
	`it is advisable to reach out to the application developer. ` +
	`In such cases, apply the patches or updates provided by ` +
	`the developer to effectively remediate the vulnerability.`

const NumberOfTARowsForDefaultRemediation = 6

// CalculateTextAreaRow calculates the number of rows for a text area
func CalculateTextAreaRow(inputData string) int {
	inputData = strings.TrimSpace(inputData)
	lineCount := strings.Count(inputData, "\n")
	lineCount += len(inputData) / 100
	lineCount = lineCount / 100
	return lineCount
}

// GetBgClassBySeverity returns the background class based on severity
func GetBgClassBySeverity(severity enums.Severity) string {
	switch severity {
	case enums.SeverityCritical:
		return "severity-bg-critical"
	case enums.SeverityHigh:
		return "severity-bg-high"
	case enums.SeverityMedium:
		return "severity-bg-medium"
	case enums.SeverityLow:
		return "severity-bg-low"
	case enums.SeverityInfo:
		return "severity-bg-info"
	default:
		return "severity-bg-info"
	}
}

// GetBgColorBySeverity returns the background color based on severity
func GetBgColorBySeverity(severity enums.Severity) string {
	switch severity {
	case enums.SeverityCritical:
		return "#e83123"
	case enums.SeverityHigh:
		return "#e77f34"
	case enums.SeverityMedium:
		return "#e6ac30"
	case enums.SeverityLow:
		return "#2fa84d"
	case enums.SeverityInfo:
		return "#0773b8"
	default:
		return "#0773b8"
	}
}
