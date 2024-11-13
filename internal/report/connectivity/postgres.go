package connectivity

import (
	"context"
	"database/sql"
	"fmt"

	types "github.com/accuknox/rinc/types/connectivity"

	_ "github.com/lib/pq"
)

// Report reports the postgres connectivity status.
func (r Reporter) postgresReport(ctx context.Context) (*types.Postgres, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable",
		r.conf.Postgres.Host,
		r.conf.Postgres.Port,
		r.conf.Postgres.Username,
		r.conf.Postgres.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging server: %w", err)
	}

	return &types.Postgres{
		Reachable: true,
	}, nil
}
