package enums

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

type ScanFlag int

const (
	ScanFlagDontScan       ScanFlag = 0
	ScanFlagWaitingToStart ScanFlag = 1
	ScanFlagIgnored        ScanFlag = 2
)
