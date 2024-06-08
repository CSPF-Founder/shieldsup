package services

// Service is the wrapper for all services
// To extend later if needed.
// Example:
// type Service struct {
// 	TargetService     TargetService
// 	ScannerService    ScannerService
// 	ScanResultService ScanResultService
// }

// type TargetService struct {
// 	TargetRepository TargetRepository
// }

// func ParseTimeString(timeString string) (time.Time, error) {
// 	// Trim any leading or trailing whitespaces
// 	timeString = strings.TrimSpace(timeString)

// 	layouts := []string{
// 		"2006-01-02T15:04:05Z",
// 		"2006-01-02 15:04:05", // Add other formats as needed
// 	}

// 	for _, layout := range layouts {
// 		parsedTime, err := time.Parse(layout, timeString)
// 		if err == nil {
// 			return parsedTime, nil
// 		}
// 	}

// 	return time.Time{}, fmt.Errorf("Failed to parse time string: %s", timeString)
// }
