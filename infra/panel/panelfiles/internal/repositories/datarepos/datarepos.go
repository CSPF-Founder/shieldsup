package datarepos

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Target     TargetRepository
	ScanResult ScanResultRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Target:     NewTargetRepository(db),
		ScanResult: NewScanResultRepository(db),
	}
}
