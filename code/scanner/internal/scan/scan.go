package scan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/config"
	"github.com/CSPF-Founder/shieldsup/scanner/enums"
	"github.com/CSPF-Founder/shieldsup/scanner/internal/repositories"
	"github.com/CSPF-Founder/shieldsup/scanner/logger"
	"github.com/CSPF-Founder/shieldsup/scanner/models"
	"github.com/CSPF-Founder/shieldsup/scanner/schemas"
	"github.com/CSPF-Founder/shieldsup/scanner/utils"
)

// type RecordData map[string]any

type ScannerModule struct {
	logger              *logger.Logger
	config              *config.Config
	db                  *repositories.Repository
	Target              models.Target
	totalWaitTime       time.Duration
	resultCheckInterval time.Duration
}

// Constructor for ScannerModule
func NewScannerModule(logger *logger.Logger, conf *config.Config, db *repositories.Repository, targetID string) (*ScannerModule, error) {
	// scanDir := app.App.Config.ScannerPath
	// if _, err := os.Stat(scanDir); os.IsNotExist(err) {
	// 	err := os.MkdirAll(scanDir, os.ModePerm) // os.ModePerm is 777, might need to check
	// 	if err != nil {
	// 		return nil
	// 	}
	// }

	target, err := db.Target.FindByID(targetID)
	if err != nil {
		return nil, fmt.Errorf("Unable to find the target with object id: %s %w", targetID, err)
	}

	return &ScannerModule{
		logger:              logger,
		config:              conf,
		db:                  db,
		Target:              target,
		totalWaitTime:       30 * time.Minute, // 30 minutes
		resultCheckInterval: 30 * time.Second, // 30 seconds
	}, nil
}

func (s *ScannerModule) HandleScanFailed(scanStatus enums.TargetStatus) bool {
	s.Target.ScanStatus = scanStatus
	isUpdated, err := s.db.Target.UpdateScanStatus(&s.Target)
	if err != nil {
		s.logger.Error("Failed to update Scan Status in Scan Failed", err)
	}
	return isUpdated
}

func (s *ScannerModule) UpdateScanStatus(status enums.TargetStatus) bool {
	s.Target.ScanStatus = status
	isUpdated, err := s.db.Target.UpdateScanStatus(&s.Target)
	if err != nil {
		s.logger.Error("Failed to update Scan Status", err)
	}
	return isUpdated
}

func (s *ScannerModule) CalculateScanTimeout() error {
	ipCount := utils.GetIPCountIfRange(s.Target.TargetAddress)

	if ipCount > 256 {
		return errors.New("ip count is more than 256")
	}

	if ipCount > 1 {
		// 15 minutes for each IP + 30 minutes for the whole range
		ipRangeTime := 15*ipCount + 30
		s.totalWaitTime = time.Duration(ipRangeTime) * time.Minute
		// Result poll interval is 1 minute
		s.resultCheckInterval = time.Minute
	}
	return nil
}

func (s *ScannerModule) Run(ctx context.Context) bool {

	isUpdated := s.UpdateScanStatus(enums.TargetStatusScanStarted)
	if !isUpdated {
		s.logger.Error("Failed to update Scan Status", nil)
		return s.HandleScanFailed(enums.TargetStatusScanFailed)
	}
	s.logger.Info(fmt.Sprintf("Doing Scan %s | Customer: %s", s.Target.TargetAddress, s.Target.CustomerUsername))

	err := s.CalculateScanTimeout()
	if err != nil {
		s.logger.Error("Error calculating scan timeout", err)
		return s.HandleScanFailed(enums.TargetStatusScanFailed)
	}

	results := s.runScan(ctx)
	if results == nil {
		// might need to change it
		s.logger.Error(fmt.Sprintf("Scan failed for %s", s.Target.TargetAddress), nil)
		return s.HandleScanFailed(enums.TargetStatusScanFailed)
	}

	handler := ResultsHandler(*s.logger, s.config.ReporterBin, s.Target, &s.db.ScanResult)
	err = handler.Run(ctx, results)
	if err != nil {
		s.logger.Error("Error in handling results", err)
		return false
	}
	return true
}

func (s *ScannerModule) runScan(ctx context.Context) []schemas.InputRecord {

	s.logger.Info(fmt.Sprintf("Starting scan for target: %s", s.Target.TargetAddress))

	apiClient := NewAPIClient(s.logger, s.config.API.Key, s.config.API.URL)
	// ! Assumption: only one scan can be started at a time
	// ! cleanup before starting (this cleans up any previous scans)
	isCleaned := apiClient.Cleanup(ctx)

	if !isCleaned {
		s.logger.Error("Error Cleaning Up", nil)
		return nil
	}

	isSuccess := apiClient.StartScan(ctx, s.Target.TargetAddress)
	if !isSuccess {
		s.logger.Error("Error adding. Doing cleanup", nil)
		apiClient.Cleanup(ctx)
		return nil
	}

	results, err := s.retrieveResult(ctx, apiClient)
	if err != nil {
		s.logger.Error("Error retrieving results", err)
		return nil
	}

	return results

}

func (s *ScannerModule) retrieveResult(ctx context.Context, apiClient *APIClient) ([]schemas.InputRecord, error) {
	endTime := time.Now().Add(s.totalWaitTime)
	for time.Now().Before(endTime) {
		response := apiClient.GetResults(ctx, s.Target.TargetAddress, false)
		if response == nil {
			err := utils.SleepContext(ctx, s.resultCheckInterval)
			if err != nil {
				return nil, err
			}
			continue
		}

		switch response.ScanStatus {
		case enums.APIScanStatusCompleted:
			return response.Data, nil
		case enums.APIScanStatusNotCompleted:
			err := utils.SleepContext(ctx, s.resultCheckInterval)
			if err != nil {
				return nil, fmt.Errorf("Error sleeping in getting result: %w", err)
			}
		case enums.APIScanStatusDoesNotExist:
			return nil, errors.New("Scan does not exist")
		case enums.APIScanStatusCouldNotCompleted:
			return nil, errors.New("Scan could not be completed")
		default:
			return nil, errors.New("Unknown scan status")
		}
	}

	s.logger.Info(fmt.Sprintf("Partial retrieval for %s", string(s.Target.TargetAddress)))
	response := apiClient.GetResults(ctx, s.Target.TargetAddress, true)
	if response == nil {
		return nil, errors.New("Error getting partial results")
	}

	return response.Data, nil
}
