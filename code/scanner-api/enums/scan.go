package enums

// Type definition for ScanStatus
type ScanStatus int

const (
	ScanStatusStarted           ScanStatus = 1
	ScanStatusNotCompleted      ScanStatus = 2
	ScanStatusCompleted         ScanStatus = 3
	ScanStatusCouldNotCompleted ScanStatus = 4
	ScanStatusDoesNotExist      ScanStatus = 501
)

// Type alias with underlying type of IntEnumMap[ScanStatus]
type ScanStatusMapType = IntEnumMap[ScanStatus]

// ScanStatusMap is the map of ScanStatus to string
var ScanStatusMap = ScanStatusMapType{
	ScanStatusStarted:           "Scan Started",
	ScanStatusNotCompleted:      "Not Completed",
	ScanStatusCompleted:         "Completed",
	ScanStatusCouldNotCompleted: "Could Not Completed",
	ScanStatusDoesNotExist:      "Does Not Exist",
}
