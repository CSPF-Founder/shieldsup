package models

import (
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/enums"

	"gorm.io/gorm"
)

const DateTimeFormat = "2006-01-02 3:4 PM"
const DELAY_BEFORE_NEW_REQUEST = 1 // Delay in hours

// Job represents the job model.
type UpdateState struct {
	ID              uint64             `gorm:"primaryKey;default:uuid_short()"`
	Status          enums.UpdateStatus `json:"status" sql:"null"`
	StatusText      string             `gorm:"-"`
	CreatedAt       time.Time          `json:"created_at" sql:"null"`
	CreatedAtString string             `gorm:"-"`
	UpdatedAt       time.Time          `json:"updated_at" sql:"null"`
	UpdatedAtString string             `gorm:"-"`
	UserID          int64              `json:"user_id" sql:"not null"`
}

// TableName overrides the table name used by User to `profiles`
func (UpdateState) TableName() string {
	return "update_state"
}

func (e *UpdateState) AfterFind(tx *gorm.DB) (err error) {
	e.CreatedAtString = e.CreatedAt.Format(DateTimeFormat)
	e.StatusText, err = enums.UpdateStatusMap.GetText(e.Status)
	if err != nil {
		return err
	}
	return
}

func FindUpdateStateByUser(id int64) (*UpdateState, error) {
	var updateState UpdateState
	err := db.Where("user_id=?", id).First(&updateState).Error
	if err != nil {
		return nil, err
	}
	return &updateState, err
}

func (us *UpdateState) GetLastUpdatedDifference() bool {
	if us.UpdatedAt.IsZero() {
		return false
	}

	timeDifference := time.Since(us.UpdatedAt)
	return int(timeDifference.Hours()) <= DELAY_BEFORE_NEW_REQUEST
}

// SaveJob saves the job to the database
func SaveUpdateState(u *UpdateState) error {
	err := db.Save(&u).Error
	db.Last(&u)
	return err
}
