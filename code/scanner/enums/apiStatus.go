package enums

// Type definition for APIScanStatus
type APIScanStatus int

const (
	APIScanStatusStarted           APIScanStatus = 1
	APIScanStatusNotCompleted      APIScanStatus = 2
	APIScanStatusCompleted         APIScanStatus = 3
	APIScanStatusCouldNotCompleted APIScanStatus = 4
	APIScanStatusDoesNotExist      APIScanStatus = 501
)

// Type alias with underlying type of IntEnumMap[APIScanStatus]
type APIScanStatusMapType = IntEnumMap[APIScanStatus]

// APIScanStatusMap is the map of APIScanStatus to string
var APIScanStatusMap = APIScanStatusMapType{
	APIScanStatusStarted:           "Scan Started",
	APIScanStatusNotCompleted:      "Not Completed",
	APIScanStatusCompleted:         "Completed",
	APIScanStatusCouldNotCompleted: "Could Not Completed",
	APIScanStatusDoesNotExist:      "Does Not Exist",
}
