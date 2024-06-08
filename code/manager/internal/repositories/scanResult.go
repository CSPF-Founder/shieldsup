package repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScanResultRepository struct {
	collection *mongo.Collection
}

func NewScanResultRepository(db *mongo.Database) ScanResultRepository {
	return ScanResultRepository{collection: db.Collection("scan_results")}
}

func (s *ScanResultRepository) DeleteByTarget(ctx context.Context, targetID primitive.ObjectID) bool {
	if targetID.IsZero() {
		return false
	}

	_, err := s.collection.DeleteMany(ctx, bson.M{"target_id": targetID})
	return err == nil
}
