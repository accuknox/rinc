package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	types "github.com/accuknox/rinc/types/rabbitmq"
)

// GetMetrics uses the RabbitMQ Management API to fetch relevant metrics for
// the report.
func (r Reporter) GetMetrics(ctx context.Context) (*types.Metrics, error) {
	overview := new(types.Overview)
	if err := r.callEndpoint(ctx, overviewEndpoint, overview); err != nil {
		return nil, fmt.Errorf("fetch overview metrics: %w", err)
	}
	nodes := new(types.Nodes)
	if err := r.callEndpoint(ctx, nodesEndpoint, nodes); err != nil {
		return nil, fmt.Errorf("fetch nodes metrics: %w", err)
	}
	queues := new(types.Queues)
	if err := r.callEndpoint(ctx, queuesEndpoint, queues); err != nil {
		return nil, fmt.Errorf("fetch queues metrics: %w", err)
	}
	consumers := new(types.Consumers)
	if err := r.callEndpoint(ctx, consumersEndpoint, consumers); err != nil {
		return nil, fmt.Errorf("fetch consumers metrics: %w", err)
	}
	exchanges := new(types.Exchanges)
	if err := r.callEndpoint(ctx, exchangesEndpoint, exchanges); err != nil {
		return nil, fmt.Errorf("fetch exchanges metrics: %w", err)
	}
	return &types.Metrics{
		IsClusterUp: true,
		Overview:    *overview,
		Nodes:       *nodes,
		Queues:      *queues,
		Consumers:   *consumers,
		Exchanges:   *exchanges,
	}, nil
}

func (r Reporter) callEndpoint(ctx context.Context, endp string, v any) error {
	endp, err := url.JoinPath(r.conf.Management.URL, endp)
	if err != nil {
		return fmt.Errorf("joining url path: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endp, nil)
	if err != nil {
		return fmt.Errorf("creating new http request: %w", err)
	}
	req.SetBasicAuth(r.conf.Management.Username, r.conf.Management.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("rabbitmq management api request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status: %s", resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if resp.StatusCode != 200 {
		return fmt.Errorf("decoding json body: %w", err)
	}
	return nil
}

func (r Reporter) callEndpointReturnStatus(ctx context.Context, endp string) (int, error) {
	endp, err := url.JoinPath(r.conf.Management.URL, endp)
	if err != nil {
		return 0, fmt.Errorf("joining url path: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endp, nil)
	if err != nil {
		return 0, fmt.Errorf("creating new http request: %w", err)
	}
	req.SetBasicAuth(r.conf.Management.Username, r.conf.Management.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("rabbitmq management api request: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

const (
	healthCheckEndpoint = "/api/health/checks/virtual-hosts"
	overviewEndpoint    = "/api/overview"
	nodesEndpoint       = "/api/nodes"
	queuesEndpoint      = "/api/queues"
	exchangesEndpoint   = "/api/exchanges"
	consumersEndpoint   = "/api/consumers"
)
