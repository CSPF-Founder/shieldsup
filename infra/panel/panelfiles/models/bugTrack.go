package models

import (
	"strconv"
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"

	"gorm.io/gorm"
)

// Job represents the job model.
type BugTrack struct {
	ID                      uint64                 `gorm:"primaryKey;default:uuid_short()"`
	UserID                  int64                  `json:"user_id" sql:"not null"`
	Status                  enums.BugTrackStatus   `json:"status" sql:"null"`
	StatusText              string                 `gorm:"-"`
	Target                  string                 `gorm:"text" json:"target" sql:"null"`
	Severity                enums.BugTrackSeverity `json:"severity"`
	SeverityText            string                 `gorm:"-"`
	AlertTitle              string                 `gorm:"text" json:"alert_title"`
	Details                 string                 `gorm:"text" json:"details"`
	Poc                     string                 `gorm:"text" json:"poc"`
	Remediation             string                 `gorm:"text" json:"remediation"`
	Remarks                 string                 `gorm:"text" json:"remarks"`
	ToBeFixedBy             string                 `gorm:"text" json:"to_be_fixed_by"`
	FoundDate               time.Time              `json:"found_date" sql:"null"`
	FormatedFoundDate       string                 `gorm:"-"`
	RevalidatedDate         time.Time              `json:"revalidated_date" sql:"null"`
	FormatedRevalidatedDate string                 `gorm:"-"`
	Likelihood              enums.Likelihood       `json:"likelihood" gorm:"type:tinyint;default:0"`
	EffortsToExploit        enums.EffortsToExploit `json:"efforts_to_exploit" gorm:"type:tinyint;default:0"`
	DataLeakage             enums.DataLeakage      `json:"data_leakage" gorm:"type:tinyint;default:0"`
	CanWafStop              enums.CanWafStop       `json:"can_waf_stop" gorm:"type:tinyint;default:0"`
	ClarificationStatus     int                    `json:"clarification_status" gorm:"default:0"`
	Prioritization          enums.Prioritization   `json:"prioritization" gorm:"default:0"`
	PrioritizationText      string                 `gorm:"-"`
	IsReviewed              int                    `json:"is_reviewed" gorm:"type:tinyint;default:0"`
	TestingMethod           enums.TestingMethod    `json:"testing_method" gorm:"type:tinyint unsigned;default:1"`
}

type OverViewParameters struct {
	User   User
	Offset int
	Limit  int
}

// TableName overrides the table name used by User to `profiles`
func (BugTrack) TableName() string {
	return "bugtrack_entries"
}

func (e *BugTrack) AfterFind(tx *gorm.DB) (err error) {
	// e.CreatedAtString = e.CreatedAt.Format(DateTimeFormat)
	e.StatusText, err = enums.BTStatusMap.GetText(e.Status)
	if err != nil {
		return err
	}
	e.SeverityText, err = enums.BTSeverityMap.GetText(e.Severity)
	if err != nil {
		return err
	}

	e.PrioritizationText, err = enums.PrioritizationMap.GetText(e.Prioritization)
	if err != nil {
		return err
	}

	e.FormatedFoundDate = e.FoundDate.Format("02-01-2006")
	e.FormatedRevalidatedDate = e.RevalidatedDate.Format("02-01-2006")
	return
}

func GetOverviewListByUser(params OverViewParameters) ([]BugTrack, error) {

	if params.Limit <= 0 {
		params.Limit = 10000
	}

	var bug_track []BugTrack
	err := db.
		Where("user_id = ?", params.User.ID).
		Order("severity ASC").
		Offset(params.Offset).
		Limit(params.Limit).
		Find(&bug_track).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return bug_track, nil
}

// SaveJob saves the job to the database
func SaveBugTrack(u *BugTrack) error {
	err := db.Save(&u).Error
	db.Last(&u)
	return err
}

// Delete deletes the Bug from the database
func DeleteBugTrack(b *BugTrack) error {
	err := db.Delete(&b).Error
	return err
}

func FindBugTrackByIdAndUser(id string, user User) (BugTrack, error) {
	var bug_track BugTrack
	validID, err := strconv.Atoi(id)
	if err != nil {
		return bug_track, err
	}

	err = db.
		Where("id = ?", validID).
		Where("user_id = ?", user.ID).
		Find(&bug_track).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return bug_track, err
	}

	return bug_track, nil
}

func GetDetailBugTrackByUser(user User) ([]BugTrack, error) {
	var bug_track []BugTrack

	err := db.
		Where("user_id = ?", user.ID).
		Order("severity ASC").
		Find(&bug_track).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return bug_track, err
	}

	return bug_track, nil
}

func CheckBugTrackAlreadyExists(user int64, target string, severity enums.BugTrackSeverity, alert_title string, details string, poc string) (int64, error) {
	var bug_track_count int64

	err := db.
		Model(&BugTrack{}).
		Where("user_id = ?", user).
		Where("target = ?", target).
		Where("alert_title = ?", alert_title).
		Where("severity = ?", severity).
		Where("details = ?", details).
		Where("poc = ?", poc).
		Count(&bug_track_count).Error

	if err != nil {
		return bug_track_count, err
	}

	return bug_track_count, nil
}
