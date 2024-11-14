package pv

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	promV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	apiEndpoint = "/api/v1/query"

	metricCapacity = iota
	metricUsed
	metricAvailable
	metricUtilization

	queryCapacity    = `sum by (persistentvolumeclaim,namespace) (kubelet_volume_stats_capacity_bytes{namespace!=""})`
	queryUsed        = `sum by (persistentvolumeclaim,namespace) (kubelet_volume_stats_used_bytes{namespace!=""})`
	queryAvailable   = `sum by (persistentvolumeclaim,namespace) (kubelet_volume_stats_available_bytes{namespace!=""})`
	queryUtilization = `sum by (persistentvolumeclaim,namespace) (kubelet_volume_stats_used_bytes{namespace!=""}) / sum by (persistentvolumeclaim,namespace) (kubelet_volume_stats_capacity_bytes{namespace!=""}) * 100`
)

var queries = map[int]string{
	metricCapacity:    queryCapacity,
	metricUsed:        queryUsed,
	metricAvailable:   queryAvailable,
	metricUtilization: queryUtilization,
}

func query(ctx context.Context, api promV1.API, q string) (model.Vector, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	result, warnings, err := api.Query(ctx, q, time.Now())
	if err != nil {
		return nil, fmt.Errorf("querying prometheus: %w", err)
	}
	for _, w := range warnings {
		slog.LogAttrs(
			ctx,
			slog.LevelWarn,
			"prometheus warning",
			slog.String("query", q),
			slog.String("message", w),
		)
	}
	vector, ok := result.(model.Vector)
	if !ok {
		return nil, fmt.Errorf(
			"query %q result: want=%s got=%s",
			q, model.ValVector, result.Type(),
		)
	}
	return vector, nil
}
