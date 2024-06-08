package repositories

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/enums"
	"github.com/CSPF-Founder/shieldsup/scanner/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TargetRepository struct {
	collection *mongo.Collection
}

func NewTargetRepository(db *mongo.Database) TargetRepository {
	return TargetRepository{
		collection: db.Collection("targets"),
	}
}

func (c *TargetRepository) FindByID(targetID string) (target models.Target, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, objectIDErr := primitive.ObjectIDFromHex(targetID)
	if objectIDErr != nil {
		return models.Target{}, objectIDErr
	}

	err = c.collection.FindOne(ctx, map[string]primitive.ObjectID{"_id": objectID}).Decode(&target)
	// Document found, return the target
	return target, err
}

func (c *TargetRepository) UpdateScanStatus(target *models.Target) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if target == nil || reflect.TypeOf(target) != reflect.TypeOf(&models.Target{}) {
		return false, errors.New("invalid target object given") // Might need to change later
	}

	now := time.Now()
	toUpdate := bson.M{}

	switch target.ScanStatus {
	case enums.TargetStatusScanStarted:
		target.ScanStartedTime = now
		toUpdate["scan_started_time"] = now
	case enums.TargetStatusReportGenerated:
		target.ScanCompletedTime = now
		toUpdate["scan_completed_time"] = now
	}

	toUpdate["scan_status"] = target.ScanStatus

	filter := bson.M{"_id": target.ID}
	update := bson.M{"$set": toUpdate}
	_, updateError := c.collection.UpdateOne(ctx, filter, update)
	if updateError != nil {
		return false, updateError
	}
	return true, nil
}

func (c *TargetRepository) UpdateScanStatusByID(ctx context.Context, targetID primitive.ObjectID, scanStatus enums.TargetStatus) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	now := time.Now()
	toUpdate := bson.M{}
	switch scanStatus {
	case enums.TargetStatusScanStarted:
		toUpdate["scan_started_time"] = now
	case enums.TargetStatusReportGenerated:
		toUpdate["scan_completed_time"] = now
	}

	filter := bson.M{"_id": targetID}
	update := bson.M{"$set": toUpdate}

	toUpdate["scan_status"] = scanStatus
	_, updateError := c.collection.UpdateOne(ctx, filter, update)
	if updateError != nil {
		return false, updateError
	}
	return true, nil
}

func (c *TargetRepository) MarkAsComplete(target *models.Target) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	toUpdate := bson.M{}
	toUpdate["scan_completed_time"] = now
	toUpdate["scan_status"] = enums.TargetStatusReportGenerated
	toUpdate["overall_cvss_score"] = target.OverallCVSSScore

	if target.IsIPRange() {
		toUpdate["cvss_score_by_host"] = target.CVSSScoreByHost
	}
	filter := bson.M{"_id": target.ID}
	update := bson.M{"$set": toUpdate}

	_, updateError := c.collection.UpdateOne(ctx, filter, update)
	if updateError != nil {
		return false, updateError
	}
	return true, nil
}

func (c *TargetRepository) GetAll() (targets []models.Target, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, findErr := c.collection.Find(ctx, map[string]string{})
	if findErr != nil {
		return nil, findErr
	}
	defer cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &targets); err != nil {
		return nil, err
	}

	return targets, nil
}
