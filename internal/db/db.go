package db

import (
	"context"
	"fmt"

	"github.com/accuknox/rinc/internal/conf"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	CollectionRabbitmq = "rabbitmq"
	CollectionCeph     = "ceph"
	CollectionImageTag = "imagetag"
	CollectionDass     = "dass"
	CollectionLongJobs = "longjobs"
)

var Collections = []string{
	CollectionRabbitmq,
	CollectionCeph,
	CollectionImageTag,
	CollectionDass,
	CollectionLongJobs,
}

// NewMongoDBClient creates a new client connection to the configured mongodb
// instance.
func NewMongoDBClient(conf conf.Mongodb) (*mongo.Client, error) {
	opts := options.
		Client().
		ApplyURI(conf.URI).
		SetAuth(options.Credential{
			Username: conf.Username,
			Password: conf.Password,
		})
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongodb: %w", err)
	}

	// ping to verify that the deployment is up and the Client was configured
	// successfully. As mentioned in the Ping documentation, this reduces
	// application resiliency as the server may be temporarily unavailable when
	// Ping is called.
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("pinging mongodb server: %w", err)
	}

	return client, nil
}

func Database(mongo *mongo.Client) *mongo.Database {
	return mongo.Database("rinc")
}
