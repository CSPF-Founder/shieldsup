package datarepos

import (
	"context"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models/datamodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScanResultRepository struct {
	collection *mongo.Collection
}

func NewScanResultRepository(db *mongo.Database) ScanResultRepository {
	return ScanResultRepository{collection: db.Collection("scan_results")}
}

func (s *ScanResultRepository) ListByTarget(ctx context.Context, target_id primitive.ObjectID) ([]datamodels.ScanResult, error) {
	var scanResult datamodels.ScanResult
	var scanResults []datamodels.ScanResult

	filter := map[string]any{
		"target_id": target_id,
	}
	options := options.Find()
	options.SetSort(map[string]int{"severity": 1})

	cursor, err := s.collection.Find(ctx, filter, options)
	if err != nil {
		defer cursor.Close(ctx)
		return scanResults, err
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&scanResult)
		if err != nil {
			return scanResults, err
		}
		scanResults = append(scanResults, scanResult)
	}

	return scanResults, nil
}

func (s *ScanResultRepository) ByID(ctx context.Context, id primitive.ObjectID) (datamodels.ScanResult, error) {
	var scanResult datamodels.ScanResult
	filter := map[string]any{
		"_id": id,
	}

	cursor := s.collection.FindOne(ctx, filter)
	if cursor.Err() != nil {
		return scanResult, cursor.Err()
	}
	err := cursor.Decode(&scanResult)
	if err != nil {
		return scanResult, err
	}
	return scanResult, nil
}

// IPListByAlert retrieves a list of distinct IPs from the collection based on the provided alert title and target ID.
func (s *ScanResultRepository) IPListByAlert(ctx context.Context, alertTitle string, targetID primitive.ObjectID) ([]string, error) {
	filter := bson.M{
		"vulnerability_title": alertTitle,
		"target_id":           targetID,
	}

	// projection := bson.M{"ip": 1}
	// opts := options.FindOne().SetProjection(projection)

	results, err := s.collection.Distinct(ctx, "ip", filter)
	if err != nil {
		return nil, err
	}

	var ips []string

	for _, result := range results {
		ips = append(ips, result.(string))
	}

	return ips, nil
}
