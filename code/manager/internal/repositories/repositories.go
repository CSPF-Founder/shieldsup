package repositories

import "go.mongodb.org/mongo-driver/mongo"

type Repository struct {
	Target     TargetRepository
	Scanner    ScannerRepository
	ScanResult ScanResultRepository
}

func NewRepository(db *mongo.Database, deploymentType string) *Repository {
	return &Repository{
		Target:     NewTargetRepository(db, deploymentType),
		Scanner:    NewScannerRepository(db),
		ScanResult: NewScanResultRepository(db),
	}
}
