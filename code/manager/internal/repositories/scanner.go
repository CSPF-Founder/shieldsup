package repositories

import (
	"context"
	"fmt"

	"github.com/CSPF-Founder/shieldsup/common/manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ScannerRepository struct {
	collection *mongo.Collection
}

func NewScannerRepository(db *mongo.Database) ScannerRepository {
	return ScannerRepository{collection: db.Collection("scanners")}
}

func (s *ScannerRepository) FindById(ctx context.Context, scannerID primitive.ObjectID) (*models.Scanner, error) {
	if scannerID.IsZero() {
		return nil, fmt.Errorf("Invalid Mongodb object id given")
	}

	var scanner models.Scanner
	err := s.collection.FindOne(ctx, bson.M{"_id": scannerID}).Decode(&scanner)
	if err != nil {
		return nil, err
	}

	return &scanner, nil
}

func (s *ScannerRepository) FindByScannerIP(ctx context.Context, scannerIP string) (*models.Scanner, error) {
	if scannerIP == "" {
		return nil, fmt.Errorf("Scanner IP is empty")
	}

	var scanner models.Scanner
	err := s.collection.FindOne(ctx, bson.M{"scanner_ip": scannerIP}).Decode(&scanner)
	if err != nil {
		return nil, err
	}

	return &scanner, nil
}
