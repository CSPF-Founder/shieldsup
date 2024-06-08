package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scanner struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ScannerName     string             `json:"scanner_name" bson:"scanner_name"`
	ScannerUsername string             `json:"scanner_username" bson:"scanner_username"`
	ScannerIP       string             `json:"scanner_ip" bson:"scanner_ip"`
}
