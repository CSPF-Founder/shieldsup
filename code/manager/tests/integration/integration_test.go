package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

// run is split from TestMain so deferred cleanup executes before os.Exit.
func run(m *testing.M) int {
	ctx := context.Background()

	pool, err := dockertest.NewPool(ctx, "")
	if err != nil {
		log.Printf("Could not construct pool: %s", err)
		return 1
	}
	defer func() {
		if err := pool.Close(ctx); err != nil {
			log.Printf("Could not close pool: %s", err)
		}
	}()

	resource, err := pool.Run(ctx, "mongo", dockertest.WithoutReuse())
	if err != nil {
		log.Printf("Could not start resource: %s", err)
		return 1
	}

	dbURI := fmt.Sprintf("mongodb://%s", resource.GetHostPort("27017/tcp"))

	if err := pool.Retry(ctx, time.Minute, func() error {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
		if err != nil {
			return err
		}

		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := client.Ping(pingCtx, nil); err != nil {
			return err
		}

		return client.Disconnect(pingCtx)
	}); err != nil {
		log.Printf("Could not connect to MongoDB: %s", err)
		return 1
	}

	return m.Run()
}
