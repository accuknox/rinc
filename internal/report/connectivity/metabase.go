package connectivity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	types "github.com/accuknox/rinc/types/connectivity"
)

const metabaseHealthEndpoint = "/api/health"

// Report reports the metabase connectivity status.
func (r Reporter) metabaseReport(ctx context.Context) (*types.Metabase, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	url, err := url.JoinPath(r.conf.Metabase.BaseURL, metabaseHealthEndpoint)
	if err != nil {
		return nil, fmt.Errorf("joining url path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating new http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response. Status: %s", resp.Status)
	}

	return &types.Metabase{
		Healthy: true,
	}, nil
}
