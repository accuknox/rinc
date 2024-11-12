package connectivity

import (
	"context"
	"fmt"

	types "github.com/accuknox/rinc/types/connectivity"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Report reports the neo4j connectivity status.
func (r Reporter) neo4jReport(ctx context.Context) (*types.Neo4j, error) {
	tkn := neo4j.BasicAuth(r.conf.Neo4j.Username, r.conf.Neo4j.Password, "")
	driver, err := neo4j.NewDriverWithContext(r.conf.Neo4j.URI, tkn)
	if err != nil {
		return nil, fmt.Errorf("creating driver: %w", err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("verifying connectivity: %w", err)
	}

	return &types.Neo4j{
		Reachable: true,
	}, nil
}
