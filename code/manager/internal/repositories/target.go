package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/CSPF-Founder/shieldsup/common/manager/enums"
	"github.com/CSPF-Founder/shieldsup/common/manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TargetRepository struct {
	collection     *mongo.Collection
	DeploymentType string
}

func NewTargetRepository(db *mongo.Database, deploymentType string) TargetRepository {
	return TargetRepository{
		collection:     db.Collection("targets"),
		DeploymentType: deploymentType,
	}
}

func (c *TargetRepository) FindById(ctx context.Context, targetID primitive.ObjectID) (*models.Target, error) {
	if targetID.IsZero() {
		return nil, errors.New("Invalid Mongodb object id given")
	}

	var target models.Target
	err := c.collection.FindOne(ctx, bson.M{"_id": targetID}).Decode(&target)
	if err != nil {
		return nil, err
	}

	return &target, nil
}

func (c *TargetRepository) UpdateScanStatus(ctx context.Context, target models.Target) error {

	now := time.Now()
	toUpdate := bson.M{}

	switch target.ScanStatus {
	case enums.TargetStatusInitiatingScan:
		target.ScanInitiatedTime = &now
		toUpdate["scan_initiated_time"] = now
		if c.DeploymentType == "cloud" {
			toUpdate["flag"] = enums.ScanFlagDontScan
		}
	case enums.TargetStatusScanStarted:
		target.ScanStartedTime = &now
		toUpdate["scan_started_time"] = now
	case enums.TargetStatusReportGenerated:
		target.ScanCompletedTime = &now
		toUpdate["scan_completed_time"] = now
	}

	toUpdate["scan_status"] = target.ScanStatus
	_, err := c.collection.UpdateOne(
		ctx,
		bson.M{"_id": target.ID},
		bson.M{"$set": toUpdate},
	)
	return err
}

func (c *TargetRepository) GetJobToScan(ctx context.Context) (*models.Target, error) {
	filters := bson.M{}

	if c.DeploymentType == "cloud" {
		filters["flag"] = enums.ScanFlagWaitingToStart
		filters["$or"] = []bson.M{
			{
				"scan_initiated_time": bson.M{
					"$lt": time.Now().Add(-6 * time.Hour),
				},
			},
			{
				"scan_initiated_time": bson.M{"$exists": false},
			},
		}
	} else {
		filters["$or"] = []bson.M{
			{"scan_status": enums.TargetStatusYetToStart},
			{"scan_status": bson.M{"$exists": false}},
		}
	}

	var target models.Target
	err := c.collection.FindOne(ctx, filters).Decode(&target)
	if err != nil {
		return nil, err
	}

	return &target, nil
}

// AnyRecentScanRunningForScanner checks if there is any recent scan running for the scanner
func (c *TargetRepository) AnyRecentScanRunningForScanner(ctx context.Context, scannerIP string) bool {
	today := time.Now()
	timeToCompare := today.Add(-time.Hour)

	filter := bson.M{
		"$and": []bson.M{
			{"scanner_ip": scannerIP},
			{
				"$and": []bson.M{
					{"scan_initiated_time": bson.M{"$exists": true}},
					{
						"scan_initiated_time": bson.M{
							"$not": bson.M{"$lte": timeToCompare},
						},
					},
				},
			},
			{
				"scan_status": bson.M{
					"$nin": []enums.TargetStatus{enums.TargetStatusReportGenerated},
				},
			},
		},
	}

	var result bson.M
	err := c.collection.FindOne(ctx, filter).Decode(&result)
	// If there is an error,
	// it means there is no recent scan running for the scanner
	return err == nil

}

func (c *TargetRepository) cursorToTargetModels(ctx context.Context, cursor *mongo.Cursor) ([]models.Target, error) {
	var targets []models.Target

	if err := cursor.All(ctx, &targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func (c *TargetRepository) GetJobListToScan(ctx context.Context) ([]models.Target, error) {
	filters := bson.M{}

	filters["flag"] = enums.ScanFlagWaitingToStart
	filters["$or"] = []bson.M{
		{
			"scan_initiated_time": bson.M{
				"$lt": time.Now().Add(-6 * time.Hour),
			},
		},
		{
			"scan_initiated_time": bson.M{"$exists": false},
		},
	}

	cursor, err := c.collection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	data, err := c.cursorToTargetModels(ctx, cursor)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *TargetRepository) MarkUnfinishedAsFailed(ctx context.Context) error {
	filters := bson.M{
		"scan_status": bson.M{
			"$nin": []enums.TargetStatus{
				enums.TargetStatusYetToStart,
				enums.TargetStatusReportGenerated,
				enums.TargetStatusScanFailed,
			},
		},
	}
	toUpdate := bson.M{"scan_status": enums.TargetStatusScanFailed}

	_, err := c.collection.UpdateMany(ctx, filters, bson.M{"$set": toUpdate})

	return err
}
