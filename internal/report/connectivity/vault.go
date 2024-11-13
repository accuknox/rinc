package connectivity

import (
	"context"
	"fmt"
	"time"

	types "github.com/accuknox/rinc/types/connectivity"

	"github.com/hashicorp/vault/api"
)

// Report reports the vault connectivity status.
func (r Reporter) vaultReport(ctx context.Context) (*types.Vault, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	client, err := api.NewClient(&api.Config{Address: r.conf.Vault.Addr})
	if err != nil {
		return nil, fmt.Errorf("creating api client: %w", err)
	}
	health, err := client.Sys().HealthWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching health: %w", err)
	}
	return &types.Vault{
		Reachable:   true,
		Initialized: health.Initialized,
		Sealed:      health.Sealed,
		Version:     health.Version,
		ClusterName: health.ClusterName,
	}, nil
}
