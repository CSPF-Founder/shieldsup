package datamodels

import (
	"errors"
	"strings"
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/config"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
	ID                primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerName      string              `bson:"customer_username"`
	TargetAddress     string              `bson:"target_address"`
	Flag              enums.ScanFlag      `bson:"flag"`
	ScanStatus        enums.TargetStatus  `bson:"scan_status"`
	ScanStartedTime   *primitive.DateTime `bson:"scan_started_time,omitempty"`
	ScanCompletedTime *primitive.DateTime `bson:"scan_completed_time,omitempty"`
	TargetType        enums.TargetType    `bson:"target_type"`
	OverallCVSSScore  float64             `bson:"overall_cvss_score,omitempty"`
	CVSSScoreByHost   map[string]float64  `bson:"cvss_score_by_host,omitempty"`
	CreatedAt         primitive.DateTime  `bson:"created_at"`
}

func (t Target) ScanStartedTimeStr() string {
	if t.ScanStartedTime == nil {
		return "-"
	}
	return t.ScanStartedTime.Time().Format("02 Jan 06 03:04:05 PM")
}

func (t Target) ScanCompletedTimeStr() string {
	if t.ScanCompletedTime == nil {
		return "-"
	}
	return t.ScanCompletedTime.Time().Format("02 Jan 06 03:04:05 PM")
}

func (t Target) GetReportDir() (string, error) {
	// Load the config
	conf, _ := config.LoadConfig()
	if strings.Contains(t.CustomerName, "/") {
		return "", errors.New("Invalid username or target ID")
	}
	return conf.SheildsUpReportDir + "/" + t.CustomerName + "/" + t.ID.Hex() + "/", nil
}

func (t Target) GetReportPath() (string, error) {
	reportDir, err := t.GetReportDir()
	if err != nil {
		return "", err
	}
	return reportDir + "report.docx", nil
}

func (t Target) CanDelete() bool {
	// check if scan started time more than 24 hour
	if t.ScanStartedTime != nil && t.ScanStartedTime.Time().Add(24*time.Hour).Before(time.Now()) {
		return false
	}

	allowedStatus := map[enums.TargetStatus]bool{
		enums.TargetStatusYetToStart:      true,
		enums.TargetStatusReportGenerated: true,
		enums.TargetStatusScanFailed:      true,
	}
	return allowedStatus[t.ScanStatus]
}
