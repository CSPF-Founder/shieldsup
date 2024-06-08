package models

// Contains Models that will be stored in the MySQL database
// User/Role/Permission and other base models are stored in the database
// Bugtracker models are used to store information about the bugs that are found

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDB(d *gorm.DB) {
	db = d
}

// Response contains the attributes found in an API response
type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}
