package integration

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var testDB *sql.DB

func TestMain(m *testing.M) {
	// Create a new Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// Uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Start a MongoDB container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "latest",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	// MongoDB connection string
	dbURI := fmt.Sprintf("mongodb://localhost:%s", resource.GetPort("27017/tcp"))

	if err := pool.Retry(func() error {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURI))
		if err != nil {
			return err
		}

		// Ping the MongoDB server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			return err
		}

		return client.Disconnect(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
		return
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
