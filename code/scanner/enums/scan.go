package enums

import "strings"

// Type definition for BugTrackSeverity
type Severity int

const (
	SeverityCritical Severity = 1
	SeverityHigh     Severity = 2
	SeverityMedium   Severity = 3
	SeverityLow      Severity = 4
	SeverityInfo     Severity = 5
)

// Type alias with underlying type of IntEnumMap[BugTrackSeverity]
type SeverityMapType = IntEnumMap[Severity]

// SeverityMap is the map of BugTrackSeverity to string
var SeverityMap = SeverityMapType{
	SeverityCritical: "Critical",
	SeverityHigh:     "High",
	SeverityMedium:   "Medium",
	SeverityLow:      "Low",
	SeverityInfo:     "Info",
}

func SeverityFromString(severity string) Severity {
	switch strings.ToLower(severity) {
	case "critical":
		return SeverityCritical
	case "high":
		return SeverityHigh
	case "medium":
		return SeverityMedium
	case "low":
		return SeverityLow
	default:
		return SeverityInfo
	}
}

// Type definition for TargetStatus
type TargetStatus int

const (
	TargetStatusYetToStart      TargetStatus = 0
	TargetStatusInitiatingScan  TargetStatus = 1
	TargetStatusScanStarted     TargetStatus = 2
	TargetStatusScanRetrieved   TargetStatus = 3
	TargetStatusReportGenerated TargetStatus = 4
	TargetStatusScanFailed      TargetStatus = 999
)

// Type alias with underlying type of IntEnumMap[TargetStatus]
type TargetStatusMapType = IntEnumMap[TargetStatus]

// TargetStatusMap is the map of TargetStatus to string
var TargetStatusMap = TargetStatusMapType{
	TargetStatusYetToStart:      "Yet To Start",
	TargetStatusInitiatingScan:  "Initiating Scan",
	TargetStatusScanStarted:     "Scan Started",
	TargetStatusScanRetrieved:   "Retrieved",
	TargetStatusReportGenerated: "Report Generated",
	TargetStatusScanFailed:      "Scan Failed",
}
