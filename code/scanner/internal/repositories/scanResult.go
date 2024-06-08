package repositories

import (
	"context"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScanResultRepository struct {
	collection *mongo.Collection
}

func NewScanResultRepository(db *mongo.Database) ScanResultRepository {
	return ScanResultRepository{collection: db.Collection("scan_results")}
}
func (c *ScanResultRepository) AddMany(records []models.ScanResult) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if records == nil {
		return 0
	}

	var interfaceRecords []interface{}
	for _, r := range records {
		interfaceRecords = append(interfaceRecords, r)
	}

	insertedResult, insertError := c.collection.InsertMany(ctx, interfaceRecords)
	if insertError != nil {
		return 0
	}
	return len(insertedResult.InsertedIDs)
}

func (c *ScanResultRepository) Exists(query map[string]any) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if id, ok := query["_id"]; ok {
		// Checking  if _id is already an ObjectID
		if oid, isObjectID := id.(primitive.ObjectID); isObjectID {

			query["_id"] = oid // If _id is already an ObjectID, directly use it in the query

		} else if idStr, isString := id.(string); isString {
			// If _id is a string, try converting it to ObjectID
			correctObjectID, err := primitive.ObjectIDFromHex(idStr)
			if err != nil {
				// Handle the error if conversion fails
				return false
			}
			query["_id"] = correctObjectID
		} else {
			return false
		}
	}
	result := c.collection.FindOne(ctx, query)
	return result.Err() == nil
}
