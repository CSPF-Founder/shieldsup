package scan

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/CSPF-Founder/shieldsup/common/manager/enums"
	"github.com/CSPF-Founder/shieldsup/common/manager/internal/repositories"
	"github.com/CSPF-Founder/shieldsup/common/manager/logger"
	"github.com/CSPF-Founder/shieldsup/common/manager/models"
	"github.com/CSPF-Founder/shieldsup/common/manager/utils"
	iputils "github.com/CSPF-Founder/shieldsup/common/manager/utils/ip_utils"
)

type Scanner struct {
	DeploymentType string
	ScannerCmd     string
	SSHKeyPath     string
	db             *repositories.Repository
	logger         *logger.Logger
	scanLogsDir    string
}

func NewScanner(
	scannerCmd string,
	deploymentType string,
	sshKeyPath string,
	db *repositories.Repository,
	lgr *logger.Logger,
	scanLogsDir string,
) *Scanner {
	return &Scanner{
		ScannerCmd:     scannerCmd,
		DeploymentType: deploymentType,
		SSHKeyPath:     sshKeyPath,
		db:             db,
		logger:         lgr,
		scanLogsDir:    scanLogsDir,
	}
}

func (s *Scanner) Run(ctx context.Context) {
	if s.DeploymentType == "cloud" {
		s.runCloud(ctx)
	} else {
		s.runOnPremise(ctx)
	}
}

func (s *Scanner) runOnPremise(ctx context.Context) {
	// mark all unfinished scans as failed
	// ! NOTE: If it is parallel scanning, then this will be a problem
	// ! Or if it is background scanning, then this will be a problem
	// * That time should not use this approach
	// * This is perfect for sequential scanning
	err := s.db.Target.MarkUnfinishedAsFailed(ctx)
	if err != nil {
		s.logger.Error("Error marking unfinished targets as failed", err)
	}

	target, err := s.db.Target.GetJobToScan(ctx)
	if err != nil {
		return
	}

	target.ScanStatus = enums.TargetStatusInitiatingScan
	err = s.db.Target.UpdateScanStatus(ctx, *target)
	if err != nil {
		s.logger.Error("Error updating target status", err)
	}

	err = s.startScannerProcess(ctx, *target)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error in startScannerProcess, targetID: %s", target.ID.Hex()), err)
	}
}

func (s *Scanner) runCloud(ctx context.Context) {
	scanRunningIPs := make([]string, 0)

	targets, err := s.db.Target.GetJobListToScan(ctx)
	if err != nil || len(targets) == 0 {
		return
	}

	for _, target := range targets {
		s.db.ScanResult.DeleteByTarget(ctx, target.ID)
		if contains(scanRunningIPs, *target.ScannerIP) {
			s.logger.Info(fmt.Sprintf("Another scan is running currently for the scanner %s", *target.ScannerIP))
			continue
		}

		if s.db.Target.AnyRecentScanRunningForScanner(ctx, *target.ScannerIP) {
			scanRunningIPs = append(scanRunningIPs, *target.ScannerIP)
			s.logger.Info(fmt.Sprintf("Another scan is running currently for the scanner %s", *target.ScannerIP))
			continue
		}

		scanner, err := s.db.Scanner.FindByScannerIP(ctx, *target.ScannerIP)
		if err != nil {
			s.logger.Error("Error occured while running Scanner", err)
			target.ScanStatus = enums.TargetStatusScanFailed
			err = s.db.Target.UpdateScanStatus(ctx, target)
			if err != nil {
				s.logger.Error("Error updating target status", err)
			}
			continue
		}

		target.ScanStatus = enums.TargetStatusInitiatingScan
		err = s.db.Target.UpdateScanStatus(ctx, target)
		if err != nil {
			s.logger.Error("Error updating target status", err)
		}

		if err := s.startRemoteScannerProcess(ctx, target, scanner); err != nil {
			s.logger.Error(fmt.Sprintf("Error in startRemoteScannerProcess %s", target.ID.Hex()), err)
			continue
		}
	}

	err = utils.SleepContext(ctx, 10*time.Second)
	if err != nil {
		s.logger.Error("runCloud: Error sleeping", err)
	}

	if len(scanRunningIPs) > 0 {
		err = utils.SleepContext(ctx, 5*time.Minute)
		if err != nil {
			s.logger.Error("runCloud: Error sleeping", err)
		}
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func calcScanTimeout(target models.Target) time.Duration {
	ipCount, _ := iputils.GetIPCountIfRange(target.TargetAddress)

	if ipCount > 1 {
		ipRangeTime := 15*ipCount + 45
		return time.Duration(ipRangeTime) * time.Minute
	}

	return 45 * time.Minute
}

func (s *Scanner) startScannerProcess(ctx context.Context, target models.Target) error {
	s.logger.Info(fmt.Sprintf("Starting scan for target %s", target.ID.Hex()))
	scanTimeout := calcScanTimeout(target)
	scanCtx, cancel := context.WithTimeout(ctx, scanTimeout)
	defer cancel()

	s.logger.Info(fmt.Sprintf("Starting scan for target %s", target.ID.Hex()))
	cmd := exec.CommandContext(scanCtx, s.ScannerCmd, "-t", target.ID.Hex())

	scanLogPath := s.scanLogsDir + target.ID.Hex() + ".log"

	// set environment variables
	cmd.Env = append(cmd.Env, "LOG_FILE_PATH="+scanLogPath)

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		errMsg := ""
		errMsg = err.Error()

		if stdOut.String() != "" {
			errMsg = errMsg + "\n" + stdOut.String()
		}

		if stdErr.String() != "" {
			errMsg = errMsg + "\n" + stdErr.String()
		}
		return errors.New(errMsg)
	}

	if stdErr.Len() > 0 {
		return fmt.Errorf("scanner process error: %s", stdErr.String())
	}

	s.logger.Info(fmt.Sprintf("Scan process done for target %s", target.ID.Hex()))

	updatedTarget, err := s.db.Target.FindById(ctx, target.ID)
	if err != nil {
		return err
	}

	if updatedTarget.ScanStatus != enums.TargetStatusReportGenerated && updatedTarget.ScanStatus != enums.TargetStatusScanFailed {
		s.logger.Error(fmt.Sprintf("Scan failed for target %s", target.ID.Hex()), nil)
		updatedTarget.ScanStatus = enums.TargetStatusScanFailed
		err = s.db.Target.UpdateScanStatus(ctx, *updatedTarget)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) startRemoteScannerProcess(ctx context.Context, target models.Target, scanner *models.Scanner) error {

	s.logger.Info(fmt.Sprintf("Starting scan for target %s", target.ID.Hex()))
	if scanner.ScannerIP == "" {
		return errors.New("Scanner IP is empty")
	}

	targetIDStr := target.ID.Hex()

	sshClient, err := utils.NewSSHConnector().Connect(scanner.ScannerIP, scanner.ScannerUsername, s.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("Error connecting to scanner, %w", err)
	}
	defer sshClient.Close()

	cmdStr := s.ScannerCmd + " -t " + targetIDStr + " > /dev/null 2>&1 &"
	// Create a session. It is one session per command.
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("Error creating new session, %w", err)
	}
	defer session.Close()
	var b bytes.Buffer  // import "bytes"
	session.Stdout = &b // get output

	err = session.Run(cmdStr)
	if err != nil {
		return fmt.Errorf("Error while starting scan for target %d, %w", target.ID, err)
	}

	updateTarget, err := s.db.Target.FindById(ctx, target.ID)
	if err != nil {
		return fmt.Errorf("Target not found %s, %w", target.ID.Hex(), err)
	}

	if updateTarget.ScanStatus != enums.TargetStatusReportGenerated && updateTarget.ScanStatus != enums.TargetStatusScanFailed {
		updateTarget.ScanStatus = enums.TargetStatusScanFailed
		err = s.db.Target.UpdateScanStatus(ctx, *updateTarget)
		if err != nil {
			return fmt.Errorf("Error updating target status %s, %w", target.ID.Hex(), err)
		}
	}
	return nil
}
