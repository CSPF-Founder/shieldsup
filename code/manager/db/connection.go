package db

import (
	"context"
	"fmt"
	"time"

	"github.com/CSPF-Founder/shieldsup/common/manager/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectDB(ctx context.Context, dbURI string, dbName string) (*mongo.Database, error) {
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))

	if err != nil {
		return nil, err
	}
	err = db.Ping(ctx, readpref.Primary())

	if err == nil {
		return nil, err
	}
	database := db.Database(dbName)
	return database, nil
}

func ConnectDBWithRetry(ctx context.Context, dbURI string, dbName string, maxRetries int, initialWaitTime time.Duration) (*mongo.Database, error) {
	var db *mongo.Database

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
		defer cancel()
		//Mongo
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))

		if err != nil {
			waitTime := initialWaitTime * time.Duration(attempt)
			err := utils.SleepContext(ctx, waitTime)
			if err != nil {
				return nil, err
			}
			continue
		}

		err = client.Ping(ctx, readpref.Primary())

		if err == nil {
			db = client.Database(dbName)
			return db, nil
		}

		waitTime := initialWaitTime * time.Duration(attempt)
		err = utils.SleepContext(ctx, waitTime)
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}
