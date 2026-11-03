package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnectDB(t *testing.T) {
	pool := dockertest.NewPoolT(t, "")
	resource := pool.RunT(t, "mongo", dockertest.WithoutReuse())

	dbURI := fmt.Sprintf("mongodb://%s", resource.GetHostPort("27017/tcp"))

	if err := pool.Retry(t.Context(), time.Minute, func() error {
		client, err := mongo.Connect(t.Context(), options.Client().ApplyURI(dbURI))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		if err := client.Ping(ctx, nil); err != nil {
			return err
		}

		return client.Disconnect(ctx)
	}); err != nil {
		t.Fatalf("Could not connect to MongoDB: %s", err)
	}
}
