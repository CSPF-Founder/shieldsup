package models

import (
	"regexp"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Target struct {
	ID               primitive.ObjectID `bson:"_id"`
	CustomerUsername string             `bson:"customer_username"`
	TargetAddress    string             `bson:"target_address"`
	Flag             int                `bson:"flag"`
	ScanStatus       enums.TargetStatus `bson:"scan_status"`
	CreatedAt        time.Time          `bson:"createdAt"`

	ScannerIP         string    `bson:"scanner_ip"`
	ScannerUsername   string    `bson:"scanner_username"`
	ScanStartedTime   time.Time `bson:"scan_started_time"`
	ScanCompletedTime time.Time `bson:"sscan_completed_time"`

	OverallCVSSScore *float64
	CVSSScoreByHost  map[string]float64
}

func (s *Target) IsIPRange() bool {
	ipRangeRegex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}$`)
	return ipRangeRegex.MatchString(s.TargetAddress)
}
