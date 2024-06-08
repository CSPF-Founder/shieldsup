package repositories

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	ScanResult ScanResultRepository
	Target     TargetRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		ScanResult: NewScanResultRepository(db),
		Target:     NewTargetRepository(db),
	}
}
