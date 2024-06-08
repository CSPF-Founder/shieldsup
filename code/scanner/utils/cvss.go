package utils

import "github.com/CSPF-Founder/shieldsup/scanner/enums"

// CalculateCVSSBySeverity calculates the CVSS score based on the severity.
func CalculateCVSSBySeverity(severity enums.Severity) float64 {
	switch severity {
	case enums.SeverityCritical:
		return 9.0
	case enums.SeverityHigh:
		return 7.0
	case enums.SeverityMedium:
		return 4.0
	case enums.SeverityLow:
		return 1.0
	default:
		return 0.0
	}
}
