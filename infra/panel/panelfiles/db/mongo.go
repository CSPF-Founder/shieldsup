package db

import (
	"context"

	"github.com/CSPF-Founder/shieldsup/onpremise/panel/config"
	"github.com/CSPF-Founder/shieldsup/onpremise/panel/internal/repositories/datarepos"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*Setup opens a database connection to mongodb*/
func SetupMongo(ctx context.Context, c *config.Config) (*datarepos.Repository, error) {
	connectionURI := c.MongoDatabaseURI
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	mongo_db := client.Database(c.MongoDatabaseName)
	repo := datarepos.NewRepository(mongo_db)
	return repo, nil
}
