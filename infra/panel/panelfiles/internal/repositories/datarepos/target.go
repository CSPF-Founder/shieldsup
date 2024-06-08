package datarepos

import (
	"context"
	"errors"
	"time"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/validator"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/models/datamodels"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TargetRepository struct {
	collection *mongo.Collection
}

func NewTargetRepository(db *mongo.Database) TargetRepository {
	return TargetRepository{
		collection: db.Collection("targets"),
	}
}

func NewTargetService(db *mongo.Database) TargetRepository {
	return TargetRepository{
		collection: db.Collection("targets"),
	}
}

func (s *TargetRepository) SaveTarget(ctx context.Context, t *datamodels.Target) (*datamodels.Target, error) {
	t.CreatedAt = primitive.NewDateTimeFromTime(time.Now().UTC())
	result, err := s.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}
	t.ID = result.InsertedID.(primitive.ObjectID)
	return t, err
}

func (s *TargetRepository) CountByCustomerUsername(ctx context.Context, customerUsername string) (int, error) {
	if !validator.IsValidUsername(customerUsername) {
		return 0, errors.New("Invalid Username")
	}

	filter := bson.M{
		"customer_username": customerUsername,
	}

	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *TargetRepository) ListByCustomerUsername(
	ctx context.Context,
	customerUsername string,
	offset int64,
	limit int64,
) ([]datamodels.Target, error) {
	if !validator.IsValidUsername(customerUsername) {
		return nil, errors.New("Invalid Username")
	}

	if limit < 0 {
		limit = 100000
	}

	options := options.Find()
	options.SetSort(map[string]int{"created_at": -1})
	options.SetLimit(limit)
	options.SetSkip(offset)

	filter := map[string]any{
		"customer_username": customerUsername,
	}

	cursor, err := s.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []datamodels.Target
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *TargetRepository) findByFilters(ctx context.Context, filter any) (*datamodels.Target, error) {

	var target datamodels.Target
	err := s.collection.FindOne(ctx, filter).Decode(&target)
	if err != nil {
		return nil, err
	}

	return &target, err
}

func (s *TargetRepository) ByIdAndCustomerUsername(ctx context.Context, id string, customerUsername string) (*datamodels.Target, error) {
	if id == "" || !validator.IsValidUsername(customerUsername) {
		return nil, errors.New("Invalid Username")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid ObjectID")
	}

	filter := bson.D{
		{Key: "customer_username", Value: customerUsername},
		{Key: "_id", Value: objectID},
	}

	return s.findByFilters(ctx, filter)
}

func (s *TargetRepository) DeleteTargetByID(ctx context.Context, id primitive.ObjectID) (int, error) {
	deleteResult, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	if deleteResult.DeletedCount > 0 {
		return int(deleteResult.DeletedCount), nil
	}

	return 0, errors.New("No document deleted")
}

func (s *TargetRepository) ListByIdsAndCustomerUsername(ctx context.Context, targetIDs []string, customerUsername string) ([]datamodels.Target, error) {
	if len(targetIDs) == 0 || customerUsername == "" || !validator.IsValidUsername(customerUsername) {
		return nil, errors.New("Invalid inputs")
	}

	dbIDs := []primitive.ObjectID{}
	for _, targetID := range targetIDs {
		if targetID != "" {
			dbID, _ := primitive.ObjectIDFromHex(targetID)
			dbIDs = append(dbIDs, dbID)
		}
	}

	filters := bson.M{
		"$and": []bson.M{
			{"customer_username": customerUsername},
			{"_id": bson.M{"$in": dbIDs}},
		},
	}

	return s.getListByFilters(ctx, filters)
}

func (s *TargetRepository) getListByFilters(ctx context.Context, filters bson.M) ([]datamodels.Target, error) {
	if len(filters) == 0 {
		return nil, errors.New("Invalid filters")
	}

	documents, err := s.collection.Find(ctx, filters)
	if err != nil {
		return nil, err
	}

	objectList := []datamodels.Target{}
	for documents.Next(ctx) {
		var document datamodels.Target
		if err := documents.Decode(&document); err != nil {
			return nil, err
		}
		objectList = append(objectList, document)
	}

	return objectList, nil
}
