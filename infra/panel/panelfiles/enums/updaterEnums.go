package enums

// Type definition for UpdateStatus
type UpdateStatus int

const (
	UpdateStatusUpdating     UpdateStatus = 0
	UpdateStatusUpdated      UpdateStatus = 1
	UpdateStatusUpdateFailed UpdateStatus = 2
)

// Type alias with underlying type of IntEnumMap[UpdateStatus]
type UpdateStatusMapType = IntEnumMap[UpdateStatus]

// UpdateStatusMap is the map of UpdateStatus to string
var UpdateStatusMap = UpdateStatusMapType{
	UpdateStatusUpdating:     "Updating",
	UpdateStatusUpdated:      "Updated",
	UpdateStatusUpdateFailed: "Update Failed",
}

// Colors for the UpdateStatus
var UpdateStatusColors = map[UpdateStatus]string{
	UpdateStatusUpdating:     "badge bg-warning",
	UpdateStatusUpdated:      "badge bg-success",
	UpdateStatusUpdateFailed: "badge bg-danger",
}

// Get background color from status
func BGColorFromStatus(value UpdateStatus) string {
	if color, ok := UpdateStatusColors[value]; ok {
		return color
	}
	return "badge bg-light"
}
