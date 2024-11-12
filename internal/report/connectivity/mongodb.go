package connectivity

import (
	"context"
	"fmt"

	types "github.com/accuknox/rinc/types/connectivity"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Report reports the mongodb connectivity status.
func (r Reporter) mongodbReport(ctx context.Context) (*types.Mongodb, error) {
	opts := options.Client().ApplyURI(r.conf.Mongodb.URI)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongodb: %w", err)
	}
	defer client.Disconnect(ctx)

	// ping to verify that the deployment is up and the Client was configured
	// successfully.
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("pinging mongodb server: %w", err)
	}

	return &types.Mongodb{
		Reachable: true,
	}, nil
}
