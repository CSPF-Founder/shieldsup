package scan

import (
	"testing"
)

// TODO: Add tests for the scanner module
func TestScanner(t *testing.T) {
	// conf, err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatal("Error loading config", err)
	// }

	// targetID, err := primitive.ObjectIDFromHex("65eed219e7e71042c5399e3a")
	// if err != nil {
	// 	log.Fatal("Error loading config", err)
	// }
	// target := models.Target{
	// 	ID:               targetID,
	// 	CustomerUsername: "test",
	// 	TargetAddress:    "192.168.2.1",
	// 	ScanStatus:       99,
	// }
	// scanner := NewScannerModule(&logger.Logger{}, conf, target)
	// got := scanner.Run(context.Background())
	// want := false
	// if got != want {
	// 	t.Errorf("got %t, want %t", got, want)
	// }
}

// func TestResultHandler(t *testing.T) {
// 	targetID, err := primitive.ObjectIDFromHex("65eed219e7e71042c5399e3a")
// 	target := models.Target{
// 		ID:               targetID,
// 		CustomerUsername: "test",
// 		TargetAddress:    "192.168.2.1",
// 		ScanStatus:       99,
// 	}

// 	resultHandler := ResultsHandler(target, logger.Logger{})
// 	resultHandler.Run()
// }
