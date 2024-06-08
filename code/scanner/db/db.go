package db

import (
	"context"
	"time"

	"github.com/CSPF-Founder/shieldsup/scanner/config"
	"github.com/CSPF-Founder/shieldsup/scanner/internal/repositories"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*Setup opens a database connection to mongodb*/
func SetupDatabase(ctx context.Context, conf *config.Config) (*repositories.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(conf.DatabaseURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	mongo_db := client.Database(conf.DBName)
	repo := repositories.NewRepository(mongo_db)
	return repo, nil
}
