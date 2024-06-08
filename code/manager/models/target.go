package models

import (
	"time"

	"github.com/CSPF-Founder/shieldsup/common/manager/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomerUsername  string             `json:"customer_username" bson:"customer_username"`
	TargetAddress     string             `json:"target_address" bson:"target_address"`
	Flag              int                `json:"flag" bson:"flag"`
	ScanStatus        enums.TargetStatus `json:"scan_status" bson:"scan_status"`
	CreatedAt         *time.Time         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ScannerIP         *string            `json:"scanner_ip,omitempty" bson:"scanner_ip,omitempty"`
	ScannerUsername   *string            `json:"scanner_username,omitempty" bson:"scanner_username,omitempty"`
	ScanStartedTime   *time.Time         `json:"scan_started_time,omitempty" bson:"scan_started_time,omitempty"`
	ScanCompletedTime *time.Time         `json:"scan_completed_time,omitempty" bson:"scan_completed_time,omitempty"`
	ScanInitiatedTime *time.Time         `json:"scan_initiated_time,omitempty" bson:"scan_initiated_time,omitempty"`
}

func (t *Target) GetScanStartedTime() string {
	if t.ScanStartedTime == nil {
		return ""
	}
	return t.ScanStartedTime.Format("2006-01-02 15:04")
}

func (t *Target) GetScanCompletedTime() string {
	if t.ScanCompletedTime == nil {
		return ""
	}
	return t.ScanCompletedTime.Format("2006-01-02 15:04")
}
